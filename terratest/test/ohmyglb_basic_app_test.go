package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Basic ohmyglb deployment test that is verifying that associated ingress is getting created
func TestOhmyglbBasicAppExample(t *testing.T) {
	t.Parallel()

	// Path to the Kubernetes resource config we will test
	kubeResourcePath, err := filepath.Abs("../examples/roundrobin.yaml")
	require.NoError(t, err)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := fmt.Sprintf("ohmyglb-test-%s", strings.ToLower(random.UniqueId()))

	// Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	// - Random namespace
	options := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.CreateNamespace(t, options, namespaceName)

	defer k8s.DeleteNamespace(t, options, namespaceName)

	defer k8s.KubectlDelete(t, options, kubeResourcePath)

	k8s.KubectlApply(t, options, kubeResourcePath)

	k8s.WaitUntilIngressAvailable(t, options, "test-gslb", 60, 1*time.Second)
	ingress := k8s.GetIngress(t, options, "test-gslb")
	require.Equal(t, ingress.Name, "test-gslb")

	helmRepoAdd := shell.Command{
		Command: "helm",
		Args:    []string{"repo", "add", "podinfo", "https://stefanprodan.github.io/podinfo"},
	}

	helmRepoUpdate := shell.Command{
		Command: "helm",
		Args:    []string{"repo", "update"},
	}
	shell.RunCommand(t, helmRepoAdd)
	shell.RunCommand(t, helmRepoUpdate)
	helmOptions := helm.Options{
		KubectlOptions: options,
	}
	helm.Install(t, &helmOptions, "podinfo/podinfo", "frontend")

	testAppFilter := metav1.ListOptions{
		LabelSelector: "app=frontend-podinfo",
	}

	k8s.WaitUntilNumPodsCreated(t, options, testAppFilter, 1, 60, 1*time.Second)

	var testAppPods []corev1.Pod

	testAppPods = k8s.ListPods(t, options, testAppFilter)

	for _, pod := range testAppPods {
		k8s.WaitUntilPodAvailable(t, options, pod.Name, 60, 1*time.Second)
	}

	k8s.WaitUntilServiceAvailable(t, options, "frontend-podinfo", 60, 1*time.Second)

	// Totally not ideal, but we need to wait until Gslb figures out Healthy status
	// We can optimize it by waiting loop with threshold later
	time.Sleep(10 * time.Second)

	ohmyglbServiceHealth, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "gslb", "test-gslb", "-o", "jsonpath='{.status.serviceHealth}'")
	if err != nil {
		t.Errorf("Failed to get ohmyglb status with kubectl (%s)", err)
	}

	want := "'map[app1.cloud.example.com:NotFound app2.cloud.example.com:NotFound app3.cloud.example.com:Healthy]'"
	assert.Equal(t, ohmyglbServiceHealth, want)

}
