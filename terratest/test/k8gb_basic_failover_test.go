package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

// Basic k8gb deployment test that is verifying that associated ingress is getting created
// Relies on two local clusters deployed by `$make deploy-two-local-clusters`
func TestK8gbBasicFailoverExample(t *testing.T) {
	t.Parallel()

	// Path to the Kubernetes resource config we will test
	kubeResourcePath, err := filepath.Abs("../examples/failover.yaml")
	require.NoError(t, err)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := fmt.Sprintf("k8gb-test-%s", strings.ToLower(random.UniqueId()))

	// Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	// - Random namespace
	optionsContext1 := k8s.NewKubectlOptions("k3d-test-gslb1", "", namespaceName)
	optionsContext2 := k8s.NewKubectlOptions("k3d-test-gslb2", "", namespaceName)

	k8s.CreateNamespace(t, optionsContext1, namespaceName)
	k8s.CreateNamespace(t, optionsContext2, namespaceName)
	defer k8s.DeleteNamespace(t, optionsContext1, namespaceName)
	defer k8s.DeleteNamespace(t, optionsContext2, namespaceName)

	gslbName := "test-gslb"

	createGslbWithHealthyApp(t, optionsContext1, kubeResourcePath, gslbName, "terratest-failover.cloud.example.com")

	createGslbWithHealthyApp(t, optionsContext2, kubeResourcePath, gslbName, "terratest-failover.cloud.example.com")

	expectedIPs := GetIngressIPs(t, optionsContext1, gslbName)

	beforeFailoverResponse, err := DoWithRetryWaitingForValueE(
		t,
		"Wait coredns to pickup dns values...",
		300,
		1*time.Second,
		func() ([]string, error) { return Dig(t, "localhost", 5053, "terratest-failover.cloud.example.com") },
		expectedIPs)
	require.NoError(t, err)

	assert.Equal(t, beforeFailoverResponse, expectedIPs)

	k8s.RunKubectl(t, optionsContext1, "scale", "deploy", "frontend-podinfo", "--replicas=0")

	assertGslbStatus(t, optionsContext1, gslbName, "terratest-failover.cloud.example.com:Unhealthy")

	t.Run("failover happens as expected", func(t *testing.T) {
		expectedIPsAfterFailover := GetIngressIPs(t, optionsContext2, gslbName)

		afterFailoverResponse, err := DoWithRetryWaitingForValueE(
			t,
			"Wait for failover to happen and coredns to pickup new values...",
			300,
			1*time.Second,
			func() ([]string, error) { return Dig(t, "localhost", 5053, "terratest-failover.cloud.example.com") },
			expectedIPsAfterFailover)
		require.NoError(t, err)

		assert.Equal(t, afterFailoverResponse, expectedIPsAfterFailover)
	})

}
