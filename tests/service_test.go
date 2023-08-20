package query_exporter

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"path/filepath"
	"testing"
)

func TestIngressEnabledCreatesIngress(t *testing.T) {
	t.Parallel()

	helmChartPath, err := filepath.Abs(filepath.Join("..", "..", "charts", "query-exporter"))
	require.NoError(t, err)

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
