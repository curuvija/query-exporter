//go:build integration
// +build integration

package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	httphelper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	t := testing.T{}

	kubectlOptions := k8s.KubectlOptions{
		ContextName: kubectlContext,
		Namespace:   namespaceName,
	}

	k8s.CreateNamespace(&t, &kubectlOptions, namespaceName)
	//defer k8s.DeleteNamespace(&t, &kubectlOptions, namespaceName)

	// kubectl -n terratest create secret generic --from-file=./config.yaml query-exporter-config-secret --output=yaml --dry-run=true > query-exporter-config-secret.yaml
	command := shell.Command{
		Command: "kubectl",
		Args: []string{
			"-n",
			"terratest",
			"create",
			"secret",
			"generic",
			"--from-file=./config.yaml",
			"query-exporter-config-secret",
		},
		WorkingDir: "",
		Env:        nil,
		Logger:     nil,
	}

	err := shell.RunCommandE(&t, command)
	require.NoError(&t, err)

	prometheusOptions := &helm.Options{
		ValuesFiles:    []string{"values/kube-prometheus-stack-values.yaml"},
		Version:        "48.3.2",
		KubectlOptions: &kubectlOptions,
		ExtraArgs: map[string][]string{
			"install": {"--wait"},
		},
	}

	// install prometheus
	helm.AddRepo(&t, prometheusOptions, "prometheus-community", "https://prometheus-community.github.io/helm-charts")
	helm.Install(&t, prometheusOptions, "prometheus-community/kube-prometheus-stack", "kube-prometheus-stack")

	// install postgres
	postgresOptions := &helm.Options{
		Version:        "12.8.5",
		KubectlOptions: &kubectlOptions,
		ValuesFiles:    []string{"values/postgres-values.yaml"},
		ExtraArgs: map[string][]string{
			"install": {"--wait"},
		},
	}
	helm.Install(&t, postgresOptions, "oci://registry-1.docker.io/bitnamicharts/postgresql", "postgres")

	// install query-exporter helm chart
	queryExporterOptions := &helm.Options{
		KubectlOptions: &kubectlOptions,
		ValuesFiles:    []string{"values/query-exporter-values.yaml"},
		ExtraArgs: map[string][]string{
			"install": {"--wait"},
		},
	}
	helm.Install(&t, queryExporterOptions, helmChartPath, releaseName)

	// wait for prometheus to scrape at least once
	time.Sleep(60 * time.Second)
	m.Run()
}

func TestQueryExporterServiceExposesMetrics(t *testing.T) {
	kubectlOptions := k8s.KubectlOptions{
		ContextName: kubectlContext,
		Namespace:   namespaceName,
	}
	tunnel := k8s.NewTunnel(&kubectlOptions, k8s.KubeResourceType(2), "query-exporter", 9560, 9560)
	tunnel.ForwardPort(t)
	status, body := httphelper.HttpGet(t, "http://localhost:9560/metrics", nil)
	require.Equal(t, 200, status)
	require.Contains(t, body, "database_errors_total")
}

func TestPrometheusScrapeMetricsUsingServiceMonitor(t *testing.T) {
	// curl 'http://localhost:9090/api/v1/targets' '
	kubectlOptions := k8s.KubectlOptions{
		ContextName: kubectlContext,
		Namespace:   namespaceName,
	}

	// port-forward do service
	tunnel := k8s.NewTunnel(
		&kubectlOptions,
		k8s.KubeResourceType(2),
		"kube-prometheus-stack-prometheus",
		9090,
		9090,
	)
	tunnel.ForwardPort(t)

	// check if servicemonitor configuration is present in prometheus configuration
	// curl 'http://localhost:9090/api/v1/targets?scrapePool=serviceMonitor/terratest/query-exporter/0' | jq '.data.activeTargets
	targetStatus, targetBody := httphelper.HttpGet(t, "http://localhost:9090/api/v1/targets?scrapePool=serviceMonitor/terratest/query-exporter/0", nil)
	require.Equal(t, 200, targetStatus)
	require.NotEmpty(t, targetBody)

	// check if prometheus scrapes query-exporter
	status, body := httphelper.HttpGet(t, "http://localhost:9090/api/v1/query?query=queries_total{job='query-exporter'}", nil)
	require.Equal(t, 200, status)
	require.NotEmpty(t, body)
}

func TestPrometheusScrapeMetricsUsingStandardConfiguration(t *testing.T) {
	kubectlOptions := k8s.KubectlOptions{
		ContextName: kubectlContext,
		Namespace:   namespaceName,
	}

	// port-forward do service
	tunnel := k8s.NewTunnel(
		&kubectlOptions,
		k8s.KubeResourceType(2),
		"kube-prometheus-stack-prometheus",
		9090,
		9090,
	)
	tunnel.ForwardPort(t)

	// check if scrapePool exists
	// curl 'http://localhost:9090/api/v1/targets?scrapePool=query-exporter-scrape' | jq '.data.activeTargets'
	targetStatus, targetBody := httphelper.HttpGet(t, "http://localhost:9090/api/v1/targets?scrapePool=query-exporter-scrape", nil)
	require.Equal(t, 200, targetStatus)
	require.NotEmpty(t, targetBody)

	// check if prometheus scrapes query-exporter
	status, body := httphelper.HttpGet(t, "http://localhost:9090/api/v1/query?query=queries_total{job='query-exporter-scrape'", nil)
	require.Equal(t, 200, status)
	require.NotEmpty(t, body)
}
