package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	//"github.com/gruntwork-io/terratest/modules/random"
)

func TestPodDeploysContainerImageHelmTemplateEngine(t *testing.T) {
	helmChartPath := "../query-exporter/"
	releaseName := "query-exporter"

	// we are working in default namespace using current kubectl context
	kubectlOptions := k8s.NewKubectlOptions("", "", "default")

	// We use a fullnameOverride so we can find the Pod later during verification
	serviceName := releaseName
	options := &helm.Options{
		SetValues: map[string]string{
			"image.repository": "adonato/query-exporter",
			"image.tag":        "latest",
		},
	}

	// This is where you render your helm chart
	output := helm.RenderTemplate(t, options, helmChartPath, releaseName, []string{})

	// Make sure to delete the resources at the end of the test
	defer k8s.KubectlDeleteFromString(t, kubectlOptions, output)

	// Now use kubectl to apply the rendered template
	k8s.KubectlApplyFromString(t, kubectlOptions, output)

	verifyService(t, kubectlOptions, serviceName)
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
	tunnel.ForwardPortE(t)

	// ... and now that we have the tunnel, we will verify that we get back a 200 OK with the nginx welcome page.
	// It takes some time for the Pod to start, so retry a few times.
	endpoint := fmt.Sprintf("http://%s", tunnel.Endpoint())
	http_helper.HttpGetWithRetryWithCustomValidation(
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
