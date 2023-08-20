//go:build all || template
// +build all template

package test

import (
	//"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	//corev1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
)

// https://github.com/gruntwork-io/terratest/blob/master/test/helm_basic_example_template_test.go

func TestDeployment(t *testing.T) {
	//t.Parallel()
	// Path to the helm chart we will test
	helmChartPath := "../../charts/query-exporter/"
	releaseName := "query-exporter"
	// require.NoError(t, err)

	// Setup the args. For this test, we will set the following input values:
	options := &helm.Options{
		SetValues: map[string]string{
			"image.repository": "adonato/query-exporter",
			"image.tag":        "latest",
		},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, releaseName, []string{"templates/deployment.yaml"})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	expectedContainerImage := "adonato/query-exporter:latest"
	deploymentContainers := deployment.Spec.Template.Spec.Containers

	// actual tests against deployment
	require.Equal(t, len(deploymentContainers), 1)
	require.Equal(t, deploymentContainers[0].Image, expectedContainerImage)
	require.Equal(t, deploymentContainers[0].Ports[0].ContainerPort, int32(9560))
	require.Equal(t, deploymentContainers[0].LivenessProbe.HTTPGet.Path, "/")
	require.Equal(t, deploymentContainers[0].LivenessProbe.HTTPGet.Port.IntVal, int32(9560))
}
