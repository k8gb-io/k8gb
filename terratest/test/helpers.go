package test

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//GetIngressIPs returns slice of IP's related to ingress
func GetIngressIPs(t *testing.T, options *k8s.KubectlOptions, ingressName string) []string {
	var ingressIPs []string
	ingress := k8s.GetIngress(t, options, ingressName)
	for _, ip := range ingress.Status.LoadBalancer.Ingress {
		ingressIPs = append(ingressIPs, ip.IP)
	}
	return ingressIPs
}

//Dig gets sorted slice of records related to dnsName
func Dig(t *testing.T, dnsServer string, dnsPort int, dnsName string) ([]string, error) {
	port := fmt.Sprintf("-p%v", dnsPort)
	dnsServer = fmt.Sprintf("@%s", dnsServer)

	digApp := shell.Command{
		Command: "dig",
		Args:    []string{port, dnsServer, dnsName, "+short"},
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
			return output, nil
			t.Logf("%s returned an error: %s. Sleeping for %s and will try again.", actionDescription, err.Error(), sleepBetweenRetries)
		}

		if reflect.DeepEqual(output, expectedResult) {
			return output, err
		}

		t.Logf("%s does not match expected result. Expected:(%s). Actual:(%s). Sleeping for %s and will try again.", actionDescription, expectedResult, output, sleepBetweenRetries)
		time.Sleep(sleepBetweenRetries)
	}

	return output, retry.MaxRetriesExceeded{Description: actionDescription, MaxRetries: maxRetries}
}

func createGslbWithHealthyApp(t *testing.T, options *k8s.KubectlOptions, kubeResourcePath string, gslbName string, hostName string) {
	k8s.KubectlApply(t, options, kubeResourcePath)

	k8s.WaitUntilIngressAvailable(t, options, gslbName, 60, 1*time.Second)
	ingress := k8s.GetIngress(t, options, gslbName)
	require.Equal(t, ingress.Name, gslbName)

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

	serviceHealthStatus := fmt.Sprintf("%s:Healthy", hostName)
	assertGslbStatus(t, options, gslbName, serviceHealthStatus)
}

func assertGslbStatus(t *testing.T, options *k8s.KubectlOptions, gslbName string, serviceStatus string) {

	t.Helper()

	actualHealthStatus := func() ([]string, error) {
		k8gbServiceHealth, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "gslb", gslbName, "-o", "jsonpath='{.status.serviceHealth}'")
		if err != nil {
			t.Errorf("Failed to get k8gb status with kubectl (%s)", err)
		}
		return []string{k8gbServiceHealth}, nil
	}
	expectedHealthStatus := []string{fmt.Sprintf("'map[%s]'", serviceStatus)}
	_, err := DoWithRetryWaitingForValueE(
		t,
		"Wait for expected ServiceHealth status...",
		60,
		1*time.Second,
		actualHealthStatus,
		expectedHealthStatus)
	require.NoError(t, err)
}
