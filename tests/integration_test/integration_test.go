package integration_test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	httphelper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

var (
	helmChartPath  = "../query-exporter"
	releaseName    = "query-exporter"
	kubectlContext = "docker-desktop"
	namespaceName  = "terratest"
)

func TestMain(m *testing.M) {
	// TODO: install kube-prometheus-stack https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
	t := testing.T{}

	kubectlOptions := k8s.KubectlOptions{
		ContextName: kubectlContext,
		Namespace:   namespaceName,
	}

	k8s.CreateNamespace(&t, &kubectlOptions, namespaceName)
	//defer k8s.DeleteNamespace(&t, &kubectlOptions, namespaceName)

	prometheusOptions := &helm.Options{
		ValuesFiles:    []string{"values/kube-prometheus-stack-values.yaml"},
		Version:        "48.3.2",
		KubectlOptions: &kubectlOptions,
	}

	helm.AddRepo(&t, prometheusOptions, "prometheus-community", "https://prometheus-community.github.io/helm-charts")
	helm.Install(&t, prometheusOptions, "prometheus-community/kube-prometheus-stack", "kube-prometheus-stack")

	testingChartOptions := &helm.Options{KubectlOptions: &kubectlOptions}
	helm.Install(&t, testingChartOptions, helmChartPath, releaseName)

	postgresOptions := &helm.Options{Version: "12.8.5", KubectlOptions: &kubectlOptions}
	helm.AddRepo(&t, postgresOptions, "", "oci://registry-1.docker.io/bitnamicharts/postgresql")
	helm.Install(&t, postgresOptions, "postgresql", "postgres")
	m.Run()
}

// verifyService will open a tunnel to the Pod and hit the endpoint to verify the nginx welcome page is shown.
func verifyService(t *testing.T, kubectlOptions *k8s.KubectlOptions, serviceName string) {
	// Wait for the pod to come up. It takes some time for the Pod to start, so retry a few times.
	retries := 15
	sleep := 5 * time.Second
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, serviceName, retries, sleep)

	// We will first open a tunnel to the pod, making sure to close it at the end of the test.
	tunnel := k8s.NewTunnel(kubectlOptions, k8s.ResourceTypeService, serviceName, 9560, 9560)
	defer tunnel.Close()
	err := tunnel.ForwardPortE(t)
	require.NoError(t, err)

	// ... and now that we have the tunnel, we will verify that we get back a 200 OK with the nginx welcome page.
	// It takes some time for the Pod to start, so retry a few times.
	endpoint := fmt.Sprintf("http://%s", tunnel.Endpoint())
	httphelper.HttpGetWithRetryWithCustomValidation(
		t,
		endpoint,
		nil,
		retries,
		sleep,
		func(statusCode int, body string) bool {
			return statusCode == 200 && strings.Contains(body, "Periodically run database queries and export results to Prometheus")
		},
	)
}

func TestVerifyService(t *testing.T) {
	kubectlOptions := k8s.NewKubectlOptions(kubectlContext, "", namespaceName)

	serviceName := releaseName
	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// This is where you render your helm chart
	output := helm.RenderTemplate(t, options, helmChartPath, releaseName, []string{})

	// Make sure to delete the resources at the end of the test
	defer k8s.KubectlDeleteFromString(t, kubectlOptions, output)

	// Now use kubectl to apply the rendered template
	k8s.KubectlApplyFromString(t, kubectlOptions, output)

	verifyService(t, kubectlOptions, serviceName)
}

// TODO: verify metrics
