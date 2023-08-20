//go:build all || template
// +build all template

package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestRenderTemplatesWithDefaultValuesShouldReturnNoError(t *testing.T) {
	var err error
	var templates []string

	releaseName := "query-exporter"
	helmChartPath, err := filepath.Abs(filepath.Join("..", "..", "charts", "query-exporter"))
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
