package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

// Basic k8gb deployment test that is verifying that associated ingress is getting created
func TestK8gbIngressAnnotationFailover(t *testing.T) {
	t.Parallel()

	// Path to the Kubernetes resource config we will test
	kubeResourcePath, err := filepath.Abs("../examples/ingress-annotation-failover.yaml")
	require.NoError(t, err)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := fmt.Sprintf("k8gb-basic-example-%s", strings.ToLower(random.UniqueId()))

	// Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	// - Random namespace
	options := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.CreateNamespace(t, options, namespaceName)

	defer k8s.DeleteNamespace(t, options, namespaceName)

	defer k8s.KubectlDelete(t, options, kubeResourcePath)

	k8s.KubectlApply(t, options, kubeResourcePath)

	ingress := k8s.GetIngress(t, options, "test-gslb-annotation-failover")
	require.Equal(t, ingress.Name, "test-gslb-annotation-failover")
	assertGslbStatus(t, options, "test-gslb-annotation-failover", "app1.cloud.example.com:NotFound app2.cloud.example.com:NotFound app3.cloud.example.com:NotFound")
	assertGslbSpec(t, options, "test-gslb-annotation-failover", ".spec.strategy.type", "failover")
	assertGslbSpec(t, options, "test-gslb-annotation-failover", ".spec.strategy.primaryGeoTag", "eu")
}
