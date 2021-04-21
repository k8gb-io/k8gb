/*

		afterFailoverResponse, err := DoWithRetryWaitingForValueE(
			t,
			"Wait for failover to happen and coredns to pickup new values(cluster1)...",
			300,
			1*time.Second,
			func() ([]string, error) {
				return Dig(t, "localhost", dnsServer1Port, "terratest-failover-split."+dnsZone)
			},
			expectedIPsCluster1)
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
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var dnsZone = getEnv("GSLB_DOMAIN", "cloud.example.com")
var dnsServer1 = getEnv("DNS_SERVER1", "localhost")
var dnsServer1Port = getEnv("DNS_SERVER1_PORT", "5053")
var dnsServer2 = getEnv("DNS_SERVER2", "localhost")
var dnsServer2Port = getEnv("DNS_SERVER2_PORT", "5054")
var primaryGeoTag = getEnv("PRIMARY_GEO_TAG", "eu")
var secondaryGeoTag = getEnv("SECONDARY_GEO_TAG", "us")

// GetIngressIPs returns slice of IP's related to ingress
func GetIngressIPs(t *testing.T, options *k8s.KubectlOptions, ingressName string) []string {
	var ingressIPs []string
	ingress := k8s.GetIngress(t, options, ingressName)
	for _, ip := range ingress.Status.LoadBalancer.Ingress {
		ingressIPs = append(ingressIPs, ip.IP)
	}
	return ingressIPs
}

// Dig gets sorted slice of records related to dnsName
func Dig(t *testing.T, dnsServer string, dnsPort string, dnsName string, additionalArgs ...string) ([]string, error) {
	port := fmt.Sprintf("-p%s", dnsPort)
	dnsServer = fmt.Sprintf("@%s", dnsServer)

	digApp := shell.Command{
		Command: "dig",
		Args:    append([]string{port, dnsServer, dnsName, "+short"}, additionalArgs...),
	}

	digAppOut := shell.RunCommandAndGetOutput(t, digApp)
	digAppSlice := strings.Split(digAppOut, "\n")

	sort.Strings(digAppSlice)

	return digAppSlice, nil
}

// DoWithRetryWaitingForValueE Concept is borrowed from terratest/modules/retry and extended to our use case
func DoWithRetryWaitingForValueE(t *testing.T, actionDescription string, maxRetries int, sleepBetweenRetries time.Duration, action func() ([]string, error), expectedResult []string) ([]string, error) {
	var output []string
	var err error

	for i := 0; i <= maxRetries; i++ {

		output, err = action()
		if err != nil {
			t.Logf("%s returned an error: %s. Sleeping for %s and will try again.", actionDescription, err.Error(), sleepBetweenRetries)
			return output, nil
		}

		if reflect.DeepEqual(output, expectedResult) {
			return output, err
		}

		t.Logf("%s does not match expected result. Expected:(%s). Actual:(%s). Sleeping for %s and will try again.", actionDescription, expectedResult, output, sleepBetweenRetries)
		time.Sleep(sleepBetweenRetries)
	}

	return output, retry.MaxRetriesExceeded{Description: actionDescription, MaxRetries: maxRetries}
}

func createGslb(t *testing.T, options *k8s.KubectlOptions, kubeResourcePath string) {
	k8sManifestBytes, err := ioutil.ReadFile(kubeResourcePath)
	if err != nil {
		log.Fatal(err)
	}

	zoneReplacer := strings.NewReplacer("cloud.example.com", dnsZone, "eu", primaryGeoTag, "us", secondaryGeoTag)

	k8sManifestString := zoneReplacer.Replace(string(k8sManifestBytes))

	k8s.KubectlApplyFromString(t, options, k8sManifestString)
}

func installPodinfo(t *testing.T, options *k8s.KubectlOptions) {
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
		Version:        "5.2.0",
		SetValues: map[string]string{
			"image.repository": getEnv("PODINFO_IMAGE_REPO", "ghcr.io/stefanprodan/podinfo"),
		},
	}
	helm.Install(t, &helmOptions, "podinfo/podinfo", "frontend")

	testAppFilter := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=frontend-podinfo",
	}

	k8s.WaitUntilNumPodsCreated(t, options, testAppFilter, 1, 60, 1*time.Second)

	var testAppPods []corev1.Pod

	testAppPods = k8s.ListPods(t, options, testAppFilter)

	for _, pod := range testAppPods {
		k8s.WaitUntilPodAvailable(t, options, pod.Name, 60, 1*time.Second)
	}

	k8s.WaitUntilServiceAvailable(t, options, "frontend-podinfo", 60, 1*time.Second)

}

func createGslbWithHealthyApp(t *testing.T, options *k8s.KubectlOptions, kubeResourcePath string, gslbName string, hostName string) {

	createGslb(t, options, kubeResourcePath)

	k8s.WaitUntilIngressAvailable(t, options, gslbName, 60, 1*time.Second)
	ingress := k8s.GetIngress(t, options, gslbName)
	require.Equal(t, ingress.Name, gslbName)

	installPodinfo(t, options)

	serviceHealthStatus := fmt.Sprintf("%s:Healthy", hostName)
	assertGslbStatus(t, options, gslbName, serviceHealthStatus)
}

func assertGslbStatus(t *testing.T, options *k8s.KubectlOptions, gslbName string, serviceStatus string) {

	t.Helper()

	actualHealthStatus := func() ([]string, error) {
		//-o custom-columns=SERVICESTATUS:.status.serviceHealth --no-headers
		k8gbServiceHealth, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "gslb", gslbName, "-o",
			"custom-columns=SERVICESTATUS:.status.serviceHealth", "--no-headers")
		if err != nil {
			t.Logf("Failed to get k8gb status with kubectl (%s)", err)
		}
		return []string{k8gbServiceHealth}, nil
	}
	expectedHealthStatus := []string{fmt.Sprintf("map[%s]", serviceStatus)}
	_, err := DoWithRetryWaitingForValueE(
		t,
		"Wait for expected ServiceHealth status...",
		60,
		1*time.Second,
		actualHealthStatus,
		expectedHealthStatus)
	require.NoError(t, err)
}

func assertGslbSpec(t *testing.T, options *k8s.KubectlOptions, gslbName string, specPath string, expectedValue string) {
	t.Helper()
	actualValue, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "gslb", gslbName, "-o", fmt.Sprintf("custom-columns=SERVICESTATUS:%s", specPath), "--no-headers")
	require.NoError(t, err)
	assert.Equal(t, expectedValue, actualValue)
}

func assertDNSEndpointLabel(t *testing.T, options *k8s.KubectlOptions, label string) {
	t.Helper()
	k8s.RunKubectl(t, options, "get", "dnsendpoint", "-l", label)
}

func assertGslbDeleted(t *testing.T, options *k8s.KubectlOptions, gslbName string) {
	t.Helper()
	deletionExpected := []string{fmt.Sprintf("Error from server (NotFound): gslbs.k8gb.absa.oss \"%s\" not found", gslbName)}
	deletionActual, err := DoWithRetryWaitingForValueE(
		t,
		"Waiting for Gslb CR to be deleted...",
		300,
		1*time.Second,
		func() ([]string, error) {
			out, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "gslb", gslbName)
			return []string{out}, err
		},
		deletionExpected)
	require.NoError(t, err)

	assert.Equal(t, deletionExpected, deletionActual)
}

func waitForLocalGSLB(t *testing.T, options *k8s.KubectlOptions, dnsPort int, host string, expectedResult []string) (output []string, err error) {
	keepNamespace := options.Namespace
	options.Namespace = getEnv("K8GB_NAMESPACE", "k8gb")
	restoreOptionsNamespace := func(options *k8s.KubectlOptions, namespace string) {
		options.Namespace = namespace
	}
	defer restoreOptionsNamespace(options, keepNamespace)
	tunnel := k8s.NewTunnel(options, k8s.ResourceTypeService, "k8gb-coredns", dnsPort, 53)
	defer tunnel.Close()
	tunnel.ForwardPort(t)

	return DoWithRetryWaitingForValueE(
		t,
		"Wait for failover to happen and coredns to pickup new values...",
		100,
		time.Second*1,
		func() ([]string, error) { return Dig(t, "localhost", fmt.Sprint(dnsPort), host, "+tcp") },
		expectedResult)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
