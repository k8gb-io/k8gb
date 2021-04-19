/*
Copyright 2021 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/
package test

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Basic k8gb deployment test that is verifying that associated ingress is getting created
// Relies on two local clusters deployed by `$make deploy-two-local-clusters`
func TestK8gbBasicRoundRobinExample(t *testing.T) {
	t.Parallel()
	var host = "roundrobin-test." + dnsZone
	const gslbName = "roundrobin-test-gslb"

	// Path to the Kubernetes resource config we will test
	kubeResourcePath, err := filepath.Abs("../examples/roundrobin2.yaml")
	require.NoError(t, err)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := fmt.Sprintf("k8gb-test-roundrobin-%s", strings.ToLower(random.UniqueId()))

	// Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	// - Random namespace
	optionsContext1 := k8s.NewKubectlOptions(getEnv("K8GB_CLUSTER1", "k3d-test-gslb1"), "", namespaceName)
	optionsContext2 := k8s.NewKubectlOptions(getEnv("K8GB_CLUSTER2", "k3d-test-gslb2"), "", namespaceName)

	k8s.CreateNamespace(t, optionsContext1, namespaceName)
	k8s.CreateNamespace(t, optionsContext2, namespaceName)
	defer k8s.DeleteNamespace(t, optionsContext1, namespaceName)
	defer k8s.DeleteNamespace(t, optionsContext2, namespaceName)

	createGslbWithHealthyApp(t, optionsContext1, kubeResourcePath, gslbName, host)

	createGslbWithHealthyApp(t, optionsContext2, kubeResourcePath, gslbName, host)

	ingressIPs1 := GetIngressIPs(t, optionsContext1, gslbName)
	ingressIPs2 := GetIngressIPs(t, optionsContext2, gslbName)
	var expectedIPs []string
	expectedIPs = append(expectedIPs, ingressIPs2...)
	expectedIPs = append(expectedIPs, ingressIPs1...)

	sort.Strings(expectedIPs)

	t.Run("round-robin on two concurrent clusters with podinfo running", func(t *testing.T) {
		resolvedIPsdnsServer1Port, err := waitForLocalGSLB(t, dnsServer1, dnsServer1Port, host, expectedIPs)
		require.NoError(t, err)
		resolvedIPsdnsServer2Port, err := waitForLocalGSLB(t, dnsServer2, dnsServer2Port, host, expectedIPs)
		require.NoError(t, err)

		assert.NotEmpty(t, resolvedIPsdnsServer1Port)
		assert.NotEmpty(t, resolvedIPsdnsServer2Port)
		assert.Equal(t, len(resolvedIPsdnsServer1Port), len(expectedIPs))
		assert.Equal(t, len(resolvedIPsdnsServer2Port), len(expectedIPs))
		assert.ElementsMatch(t, resolvedIPsdnsServer1Port, expectedIPs, "%s:%s", host, dnsServer1Port)
		assert.ElementsMatch(t, resolvedIPsdnsServer2Port, expectedIPs, "%s:%s", host, dnsServer2Port)
	})

	t.Run("kill podinfo on the first cluster", func(t *testing.T) {
		// kill app in the first cluster
		k8s.RunKubectl(t, optionsContext1, "scale", "deploy", "frontend-podinfo", "--replicas=0")

		assertGslbStatus(t, optionsContext1, gslbName, host+":Unhealthy")

		resolvedIPsdnsServer1Port, err := waitForLocalGSLB(t, dnsServer1, dnsServer1Port, host, ingressIPs2)
		require.NoError(t, err)
		resolvedIPsdnsServer2Port, err := waitForLocalGSLB(t, dnsServer2, dnsServer2Port, host, ingressIPs2)
		require.NoError(t, err)
		assert.ElementsMatch(t, resolvedIPsdnsServer1Port, resolvedIPsdnsServer2Port)
	})

	t.Run("kill podinfo on the second cluster", func(t *testing.T) {
		// kill app in the second cluster
		k8s.RunKubectl(t, optionsContext2, "scale", "deploy", "frontend-podinfo", "--replicas=0")

		assertGslbStatus(t, optionsContext2, gslbName, host+":Unhealthy")

		_, err = waitForLocalGSLB(t, dnsServer1, dnsServer1Port, host, []string{""})
		require.NoError(t, err)
		_, err = waitForLocalGSLB(t, dnsServer2, dnsServer2Port, host, []string{""})
		require.NoError(t, err)
	})

	t.Run("start podinfo on the both clusters", func(t *testing.T) {
		// start app in the both clusters
		k8s.RunKubectl(t, optionsContext1, "scale", "deploy", "frontend-podinfo", "--replicas=1")
		k8s.RunKubectl(t, optionsContext2, "scale", "deploy", "frontend-podinfo", "--replicas=1")

		assertGslbStatus(t, optionsContext1, gslbName, host+":Healthy")
		assertGslbStatus(t, optionsContext2, gslbName, host+":Healthy")

		resolvedIPsdnsServer1Port, err := waitForLocalGSLB(t, dnsServer1, dnsServer1Port, host, expectedIPs)
		require.NoError(t, err)
		resolvedIPsdnsServer2Port, err := waitForLocalGSLB(t, dnsServer2, dnsServer2Port, host, expectedIPs)
		require.NoError(t, err)

		assert.NotEmpty(t, resolvedIPsdnsServer1Port)
		assert.NotEmpty(t, resolvedIPsdnsServer2Port)
		assert.Equal(t, len(resolvedIPsdnsServer1Port), len(expectedIPs))
		assert.Equal(t, len(resolvedIPsdnsServer2Port), len(expectedIPs))
		assert.ElementsMatch(t, resolvedIPsdnsServer1Port, expectedIPs, "%s:%s", host, dnsServer1Port)
		assert.ElementsMatch(t, resolvedIPsdnsServer2Port, expectedIPs, "%s:%s", host, dnsServer2Port)
	})
}
