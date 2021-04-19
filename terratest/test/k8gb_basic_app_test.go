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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Basic k8gb deployment test that is verifying that associated ingress is getting created
func TestK8gbBasicAppExample(t *testing.T) {
	t.Parallel()

	// Path to the Kubernetes resource config we will test
	kubeResourcePath, err := filepath.Abs("../examples/roundrobin.yaml")
	require.NoError(t, err)

	brokenResourcePath, err := filepath.Abs("../examples/broken-gslb.yaml")
	require.NoError(t, err)
	brokenNoHTTPResourcePath, err := filepath.Abs("../examples/broken-gslb-no-http.yaml")
	require.NoError(t, err)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := fmt.Sprintf("k8gb-test-%s", strings.ToLower(random.UniqueId()))

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

	// Path to the Kubernetes resource config we will test
	unhealthyAppPath, err := filepath.Abs("../examples/unhealthy-app.yaml")
	require.NoError(t, err)
	k8s.KubectlApply(t, options, unhealthyAppPath)

	helmRepoAdd := shell.Command{
		Command: "helm",
		Args:    []string{"repo", "add", "--force-update", "podinfo", "https://stefanprodan.github.io/podinfo"},
	}

	helmRepoUpdate := shell.Command{
		Command: "helm",
		Args:    []string{"repo", "update"},
	}
	shell.RunCommand(t, helmRepoAdd)
	shell.RunCommand(t, helmRepoUpdate)
	helmOptions := helm.Options{
		KubectlOptions: options,
		Version:        "4.0.6",
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

	assertGslbStatus(t, options, "test-gslb", "notfound."+dnsZone+":NotFound roundrobin."+dnsZone+":Healthy unhealthy."+dnsZone+":Unhealthy")
	// Ensure controller labels DNSEndpoint objects
	assertDNSEndpointLabel(t, options, "k8gb.absa.oss/dnstype")

	t.Run("Broken object rejected by API", func(t *testing.T) {
		err := k8s.KubectlApplyE(t, options, brokenResourcePath)
		require.Error(t, err)
		err = k8s.KubectlApplyE(t, options, brokenNoHTTPResourcePath)
		require.Error(t, err)
	})

}
