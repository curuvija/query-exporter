package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	prometheusoperatorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"path/filepath"
	"testing"
)

var (
	helmChartPath  = "../query-exporter"
	releaseName    = "query-exporter"
	kubectlContext = "docker-desktop"
	namespaceName  = "terratest"
)

func TestRenderTemplatesWithDefaultValuesShouldReturnNoError(t *testing.T) {
	var err error
	var templates []string

	require.NoError(t, err)

	templates = []string{
		"templates/deployment.yaml",
		"templates/configmap.yaml",
		"templates/service.yaml",
		"templates/serviceaccount.yaml",
		"templates/servicemonitor.yaml",
	}
	var valuesFile = []string{filepath.Join(helmChartPath, "values.yaml")}
	options := &helm.Options{
		ValuesFiles: valuesFile,
	}

	helm.RenderTemplate(t, options, helmChartPath, releaseName, templates)
}

func TestDeployment(t *testing.T) {

	// Set up the args. For this test, we will set the following input values:
	options := &helm.Options{
		SetValues: map[string]string{
			"image.repository": "adonato/query-exporter",
			"image.tag":        "latest",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, releaseName, []string{"templates/deployment.yaml"})
	fmt.Println(output)

	// Now we use kubernetes/client-go library to render the template output into the Pod struct. This will
	// ensure the Pod resource is rendered correctly.
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	expectedContainerImage := "adonato/query-exporter:latest"
	deploymentContainers := deployment.Spec.Template.Spec.Containers
	require.Equal(t, len(deploymentContainers), 1)
	require.Equal(t, deploymentContainers[0].Image, expectedContainerImage)
}

func TestServiceMonitorEnabledFalseDoesNotCreateServiceMonitor(t *testing.T) {
	t.Parallel()

	options := &helm.Options{
		ValuesFiles: []string{filepath.Join("..", "..", "charts", "query-exporter", "values.yaml")},
		SetValues:   map[string]string{"prometheus.monitor.enabled": "false"},
	}
	_, err := helm.RenderTemplateE(t, options, helmChartPath, "servicemonitor", []string{"templates/servicemonitor.yaml"})
	require.Error(t, err)
}

func TestServiceMonitorEnabledCreatesServiceMonitor(t *testing.T) {
	t.Parallel()

	options := &helm.Options{
		ValuesFiles: []string{
			filepath.Join("..", "..", "charts", "query-exporter", "values.yaml"),
		},
	}
	out := helm.RenderTemplate(t, options, helmChartPath, "servicemonitor", []string{"templates/servicemonitor.yaml"})

	rendered := prometheusoperatorv1.ServiceMonitor{}
	require.NoError(t, yaml.Unmarshal([]byte(out), &rendered))
	require.Equal(t, 1, len(rendered.Spec.Endpoints))

	//TODO: Configure other tests properly
	defaultEndpoint := rendered.Spec.Endpoints[0]
	//assert.Equal(t, "15s", defaultEndpoint.Interval)
	//assert.Equal(t, "10s", defaultEndpoint.ScrapeTimeout)
	assert.Equal(t, "/metrics", defaultEndpoint.Path)
	//assert.Equal(t, "http", defaultEndpoint.Port)
	//assert.Equal(t, "http", defaultEndpoint.Scheme)
}

func TestIngressEnabledCreatesIngress(t *testing.T) {
	t.Parallel()

	hostname := "query-exporter.mydomain.com"
	path := "/query-exporter"

	options := &helm.Options{
		ValuesFiles: []string{
			filepath.Join(helmChartPath, "values.yaml"),
		},
		SetValues: map[string]string{
			"ingress.enabled":                "true",
			"ingress.hosts[0].host":          hostname,
			"ingress.hosts[0].paths[0].path": path,
			"ingress.tls[0].secretName":      "tls-secretname",
			"ingress.tls[0].hosts[0]":        hostname,
		},
	}

	rendered, err := helm.RenderTemplateE(t, options, helmChartPath, "ingress", []string{"templates/ingress.yaml"})
	require.NoError(t, err)
	var ingress networkingv1.Ingress
	helm.UnmarshalK8SYaml(t, rendered, &ingress)

	ingressSpecRules := ingress.Spec.Rules
	require.Equal(t, ingressSpecRules[0].Host, hostname)
	require.Equal(t, ingressSpecRules[0].HTTP.Paths[0].Path, path)
	require.Equal(t, ingressSpecRules[0].HTTP.Paths[0].Backend.Service.Port.Number, int32(9560))
}

func TestDefaultValuesIngressNotEnabledDoesNotCreateIngress(t *testing.T) {
	t.Parallel()

	helmChartPath, err := filepath.Abs(filepath.Join("..", "..", "charts", "query-exporter"))
	require.NoError(t, err)

	options := &helm.Options{
		ValuesFiles: []string{
			filepath.Join(helmChartPath, "values.yaml"),
		},
	}

	_, err = helm.RenderTemplateE(t, options, helmChartPath, "ingress", []string{"templates/ingress.yaml"})
	expected := "error while running command: exit status 1; Error: could not find template templates/ingress.yaml in chart"
	require.EqualError(t, err, expected)

}

func TestServiceCreatedByDefault(t *testing.T) {
	t.Parallel()

	helmChartPath, err := filepath.Abs(filepath.Join("..", "..", "charts", "query-exporter"))
	require.NoError(t, err)

	releaseName := "query-exporter"

	options := &helm.Options{
		ValuesFiles: []string{
			filepath.Join(helmChartPath, "values.yaml"),
		},
	}

	rendered, err := helm.RenderTemplateE(t, options, helmChartPath, releaseName, []string{"templates/service.yaml"})
	var service corev1.Service
	helm.UnmarshalK8SYaml(t, rendered, &service)

	ports := service.Spec.Ports[0]
	require.Equal(t, ports.Port, int32(9560))
	require.Equal(t, ports.TargetPort, intstr.IntOrString{IntVal: 9560})
	require.Equal(t, ports.Protocol, corev1.Protocol("TCP"))
	require.Equal(t, ports.Name, "http")
	require.Equal(t, service.Spec.Type, corev1.ServiceType("ClusterIP"))
}
