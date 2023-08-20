//go:build all || template
// +build all template

package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"testing"

	prometheus_operator_v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/require"
)

func TestServiceMonitorEnabledFalseDoesNotCreateServiceMonitor(t *testing.T) {
	t.Parallel()
	helmChartPath, err := filepath.Abs(filepath.Join("..", "..", "charts", "query-exporter"))
	require.NoError(t, err)

	options := &helm.Options{
		ValuesFiles: []string{filepath.Join("..", "..", "charts", "query-exporter", "values.yaml")},
		SetValues:   map[string]string{"prometheus.monitor.enabled": "false"},
	}
	_, err = helm.RenderTemplateE(t, options, helmChartPath, "servicemonitor", []string{"templates/servicemonitor.yaml"})
	require.Error(t, err)
}

func TestServiceMonitorEnabledCreatesServiceMonitor(t *testing.T) {
	t.Parallel()

	helmChartPath, err := filepath.Abs(filepath.Join("..", "..", "charts", "query-exporter"))
	require.NoError(t, err)

	options := &helm.Options{
		ValuesFiles: []string{
			filepath.Join("..", "..", "charts", "query-exporter", "values.yaml"),
		},
	}
	out := helm.RenderTemplate(t, options, helmChartPath, "servicemonitor", []string{"templates/servicemonitor.yaml"})

	rendered := prometheus_operator_v1.ServiceMonitor{}
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
