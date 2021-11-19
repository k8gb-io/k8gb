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
package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/AbsaOSS/gopkg/dns"
	gopkgstr "github.com/AbsaOSS/gopkg/strings"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Workflow struct {
	error      error
	namespace  string
	cluster    string
	k8sOptions *k8s.KubectlOptions
	t          *testing.T
	settings   struct {
		ingressResourcePath string
		gslbResourcePath    string
		ingressName         string
	}
	state struct {
		namespaceCreated bool
		testApp          struct {
			name        string
			message     string
			isRunning   bool
			isInstalled bool
		}
		gslb struct {
			name        string
			host        string
			port        int
			isInstalled bool
		}
	}
}

type Instance struct {
	w *Workflow
}

type TestAppResult struct {
	Message string `json:"message"`
	Version string `json:"version"`
	Color   string `json:"color"`
	Pod     string `json:"hostname"`
	Body    string
}

// InstanceStatus provides a simplified overview of the instance status
type InstanceStatus struct {
	Annotation       string   `json:"annotation"`
	AppMessage       string   `json:"app-msg"`
	AppRunning       bool     `json:"podinfo-running"`
	AppReplicas      string   `json:"podinfo-replicas"`
	LocalTargets     []string `json:"local-targets-ip"`
	Ingresses        []string `json:"ingress-ip"`
	Dig              []string `json:"dig-result"`
	CoreDNS          string   `json:"coredns-ip"`
	GslbHealthStatus string   `json:"gslb-status"`
	Cluster          string   `json:"cluster"`
	Namespace        string   `json:"namespace"`
	Endpoint0DNSName string   `json:"ep0-dns-name"`
	Endpoint0Targets string   `json:"ep0-dns-targets"`
	Endpoint1DNSName string   `json:"ep1-dns-name"`
	Endpoint1Targets string   `json:"ep1-dns-targets"`
}

func NewWorkflow(t *testing.T, cluster string, port int) *Workflow {
	var err error
	if cluster == "" {
		err = fmt.Errorf("empty cluster")
	}
	if port < 1000 {
		err = fmt.Errorf("invalid port")
	}
	w := new(Workflow)
	w.cluster = cluster
	w.namespace = fmt.Sprintf("k8gb-test-%s", strings.ToLower(random.UniqueId()))
	w.k8sOptions = k8s.NewKubectlOptions(cluster, "", w.namespace)
	w.t = t
	w.state.gslb.port = port
	w.error = err
	return w
}

func (w *Workflow) WithIngress(path string) *Workflow {
	if path == "" {
		w.error = fmt.Errorf("empty ingress resource path")
	}
	w.settings.ingressResourcePath = path
	return w
}

func (w *Workflow) WithGslb(path, host string) *Workflow {
	var err error
	if host == "" {
		w.error = fmt.Errorf("empty gslb host")
	}
	if path == "" {
		w.error = fmt.Errorf("empty gslb resource path")
	}
	w.settings.gslbResourcePath, err = filepath.Abs(path)
	if err != nil {
		w.error = fmt.Errorf("reading %s; %s", path, err)
	}
	w.state.gslb.name, err = w.getManifestName(w.settings.gslbResourcePath)
	if err != nil {
		w.error = err
	}
	w.state.gslb.host = host
	if err != nil {
		w.error = err
	}
	return w
}

func (w *Workflow) WithTestApp(uiMessage string) *Workflow {
	w.state.testApp.isInstalled = true
	w.state.testApp.name = "frontend-podinfo"
	w.state.testApp.message = uiMessage
	return w
}

func (w *Workflow) Start() (*Instance, error) {
	if w.error != nil {
		return nil, w.error
	}

	// namespace
	w.t.Logf("Create namespace %s", w.namespace)
	k8s.CreateNamespace(w.t, w.k8sOptions, w.namespace)
	w.state.namespaceCreated = true

	// gslb
	if w.settings.gslbResourcePath != "" {
		w.t.Logf("Create ingress %s from %s", w.state.gslb.name, w.settings.gslbResourcePath)
		k8s.KubectlApply(w.t, w.k8sOptions, w.settings.gslbResourcePath)
		k8s.WaitUntilIngressAvailable(w.t, w.k8sOptions, w.state.gslb.name, 60, 1*time.Second)
		ingress := k8s.GetIngress(w.t, w.k8sOptions, w.state.gslb.name)
		require.Equal(w.t, ingress.Name, w.state.gslb.name)
		w.settings.ingressName = w.state.gslb.name
	}

	// ingress
	if w.settings.ingressResourcePath != "" {

	}

	// app
	if w.state.testApp.isInstalled {
		const app = "https://stefanprodan.github.io/podinfo"
		w.t.Logf("Create test application %s", app)
		helmRepoAdd := shell.Command{
			Command: "helm",
			Args:    []string{"repo", "add", "--force-update", "podinfo", app},
		}
		helmRepoUpdate := shell.Command{
			Command: "helm",
			Args:    []string{"repo", "update"},
		}
		shell.RunCommand(w.t, helmRepoAdd)
		shell.RunCommand(w.t, helmRepoUpdate)
		helmOptions := helm.Options{
			KubectlOptions: w.k8sOptions,
			Version:        "5.1.1",
			SetValues:      map[string]string{"ui.message": w.state.testApp.message},
		}
		helm.Install(w.t, &helmOptions, "podinfo/podinfo", "frontend")
		testAppFilter := metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=" + w.state.testApp.name,
		}
		k8s.WaitUntilNumPodsCreated(w.t, w.k8sOptions, testAppFilter, 1, 60, 1*time.Second)
		var testAppPods []corev1.Pod
		testAppPods = k8s.ListPods(w.t, w.k8sOptions, testAppFilter)
		for _, pod := range testAppPods {
			k8s.WaitUntilPodAvailable(w.t, w.k8sOptions, pod.Name, 60, 1*time.Second)
		}
		k8s.WaitUntilServiceAvailable(w.t, w.k8sOptions, w.state.testApp.name, 60, 1*time.Second)
		w.state.testApp.isRunning = true
	}
	return &Instance{
		w: w,
	}, nil
}

func (w *Workflow) getManifestName(path string) (string, error) {
	m := struct {
		Metadata struct {
			Name string `yaml:"name"`
		} `yaml:"metadata"`
	}{}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("parse %s; %s", path, err)
	}
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		return "", fmt.Errorf("unmarshall %s; %s", path, err)
	}
	return m.Metadata.Name, nil
}

func (i *Instance) Kill() {
	i.w.t.Logf("killing %s", i)
	if i.w.state.namespaceCreated {
		k8s.DeleteNamespace(i.w.t, i.w.k8sOptions, i.w.namespace)
	}
}

// GetCoreDNSIP gets core DNS IP address
func (i *Instance) GetCoreDNSIP() string {
	cmd := shell.Command{
		Command: "kubectl",
		Args:    []string{"--context", i.w.k8sOptions.ContextName, "-n", "k8gb", "get", "svc", "k8gb-coredns", "--no-headers", "-o", "custom-columns=IP:spec.clusterIPs[0]"},
		Env:     i.w.k8sOptions.Env,
	}
	out, err := shell.RunCommandAndGetOutputE(i.w.t, cmd)
	require.NoError(i.w.t, err)
	require.NotEqual(i.w.t, "<none>", out)
	return out
}

func (i *Instance) GetIngressIPs() []string {
	var ingressIPs []string
	ingress := k8s.GetIngress(i.w.t, i.w.k8sOptions, i.w.settings.ingressName)
	for _, ip := range ingress.Status.LoadBalancer.Ingress {
		ingressIPs = append(ingressIPs, ip.IP)
	}
	return ingressIPs
}

func (i *Instance) StopTestApp() {
	k8s.RunKubectl(i.w.t, i.w.k8sOptions, "scale", "deploy", i.w.state.testApp.name, "--replicas=0")
	AssertGslbStatus(i.w.t, i.w.k8sOptions, i.w.state.gslb.name, i.w.state.gslb.host+":Unhealthy")
	i.w.state.testApp.isRunning = false
}

func (i *Instance) StartTestApp() {
	k8s.RunKubectl(i.w.t, i.w.k8sOptions, "scale", "deploy", i.w.state.testApp.name, "--replicas=1")
	AssertGslbStatus(i.w.t, i.w.k8sOptions, i.w.state.gslb.name, i.w.state.gslb.host+":Healthy")
	i.w.state.testApp.isRunning = true
}

// WaitForGSLB waits until GSLB contains desired IP address list and if it is, the desired list is returned.
// Desired IP list is LocalTargetsIP combination of all instances with running app.
// e.g.: instance1.WaitForGSLB(instance2, instance3) produces:
// desiredIPList := instance1.GetLocalTargets() + instance2.GetLocalTargets() + instance3.GetLocalTargets()
// If app is stopped, the IP addresses are excluded from desired list.
func (i *Instance) WaitForGSLB(instances ...*Instance) ([]string, error) {
	var expectedIPs []string
	instances = append(instances, i)
	for _, in := range instances {
		// add expected IP's only if app is running
		if in.w.state.testApp.isRunning {
			ip := in.GetLocalTargets()
			expectedIPs = append(expectedIPs, ip...)
		}
	}
	return waitForLocalGSLBNew(i.w.t, i.w.state.gslb.host, i.w.state.gslb.port, expectedIPs)
}

// WaitForExpected waits until GSLB dig doesnt return list of expected IP's
func (i *Instance) WaitForExpected(expectedIPs []string) (err error) {
	_, err = waitForLocalGSLBNew(i.w.t, i.w.state.gslb.host, i.w.state.gslb.port, expectedIPs)
	if err != nil {
		fmt.Println(i.GetStatus(fmt.Sprintf("expected IPs: %s", expectedIPs)).String())
	}
	return
}

// WaitForAppIsRunning waits until app has one pod running
func (i *Instance) WaitForAppIsRunning() (err error) {
	return i.waitForApp(func(instances int) bool { return instances > 0 })
}

// WaitForAppIsStopped waits until app has 0 pods running
func (i *Instance) WaitForAppIsStopped() (err error) {
	return i.waitForApp(func(instances int) bool { return instances == 0 })
}

// WaitForAppIsRunning waits until app has one pod running
func (i *Instance) waitForApp(action func(instances int) bool) (err error) {
	const description = "Wait for app is running"
	for n := 0; n < 25; n++ {
		var r int
		rep := i.GetStatus("").AppReplicas
		// r == 0
		if rep != "<none>" {
			r, err = strconv.Atoi(rep)
			if err != nil {
				i.w.t.Logf("%s returned an error: %s.", description, err)
				return err
			}
		}
		if action(r) {
			i.w.t.Logf("%s found match: Expected:(%v)", description, r)
			return nil
		}
		i.w.t.Logf("Application %s is not in expected state. Waiting...", i.w.state.testApp.name)
	}
	return
}

// String retrieves rough information about cluster
func (i *Instance) String() (out string) {
	return fmt.Sprintf("Instance: %s:%s", i.w.cluster, i.w.namespace)
}

// Dig  returns a list of IP addresses from CoreDNS that belong to the instance
func (i *Instance) Dig() []string {
	dig, err := dns.Dig("localhost:"+strconv.Itoa(i.w.state.gslb.port), i.w.state.gslb.host)
	require.NoError(i.w.t, err)
	return dig
}

// GetLocalTargets returns instance local targets
func (i *Instance) GetLocalTargets() []string {
	dnsName := fmt.Sprintf("localtargets-%s", i.w.state.gslb.host)
	dig, err := dns.Dig("localhost:"+strconv.Itoa(i.w.state.gslb.port), dnsName)
	require.NoError(i.w.t, err)
	return dig
}

// HitTestApp makes HTTP GET to TestApp when installed otherwise panics.
// If the function successfully hits the TestApp, it returns the TestAppResult.
func (i *Instance) HitTestApp() (result *TestAppResult) {
	require.True(i.w.t, i.w.state.testApp.isInstalled)
	var err error
	result = new(TestAppResult)
	coreDNSIP := i.GetCoreDNSIP()
	command := []string{"sh", "-c", fmt.Sprintf("wget -qO - %s", i.w.state.gslb.host)}
	for t := 0; t < 3; t++ {
		result.Body, err = RunBusyBoxCommand(i.w.t, i.w.k8sOptions, coreDNSIP, command)
		require.NoError(i.w.t, err, "busybox", command, result.Body)
		if strings.HasPrefix(result.Body, "{") {
			break
		}
		i.w.t.Log("podinfo didn't start yet, waiting....")
		time.Sleep(time.Second * 1)
	}
	// unwrap json from busybox messages
	parsedJson := strings.Split(result.Body, "}")[0]
	s := strings.Split(parsedJson, "{")
	require.Len(i.w.t, s, 2, "invalid busybox response", result.Body)
	parsedJson = s[1]

	err = json.Unmarshal([]byte("{"+parsedJson+"}"), result)
	require.NoError(i.w.t, err, "unmarshall json", result.Body)
	return
}

// GetStatus reads overall status about instance. Status can be used for assertion as well as printed to test output
// Annotation argument is just free text which might be used in various test scenarios
func (i *Instance) GetStatus(annotation string) (s *InstanceStatus) {
	const na = "n/a"
	var err error
	s = new(InstanceStatus)
	s.Annotation = annotation
	s.Cluster = i.w.cluster
	s.Namespace = i.w.namespace
	s.Dig = i.Dig()
	s.LocalTargets = i.GetLocalTargets()
	s.Ingresses = i.GetIngressIPs()
	s.CoreDNS = i.GetCoreDNSIP()
	s.AppMessage = i.w.state.testApp.message
	s.AppRunning = i.w.state.testApp.isRunning
	s.AppReplicas, err = k8s.RunKubectlAndGetOutputE(i.w.t, i.w.k8sOptions, "get", "deployments", "frontend-podinfo",
		"-o", "custom-columns=STATUS:.status.readyReplicas", "--no-headers")
	if err != nil {
		s.AppReplicas = na
	}
	s.GslbHealthStatus, err = k8s.RunKubectlAndGetOutputE(i.w.t, i.w.k8sOptions, "get", "gslb", i.w.state.gslb.name, "-o",
		"custom-columns=SERVICESTATUS:.status.serviceHealth", "--no-headers")
	if err != nil {
		s.GslbHealthStatus = na
	}
	s.Endpoint0DNSName, err = k8s.RunKubectlAndGetOutputE(i.w.t, i.w.k8sOptions, "get", "dnsendpoints.externaldns.k8s.io", "test-gslb", "-o",
		"custom-columns=SERVICESTATUS:.spec.endpoints[0].dnsName", "--no-headers")
	if err != nil {
		s.Endpoint0DNSName = na
	}
	s.Endpoint0Targets, err = k8s.RunKubectlAndGetOutputE(i.w.t, i.w.k8sOptions, "get", "dnsendpoints.externaldns.k8s.io", "test-gslb", "-o",
		"custom-columns=SERVICESTATUS:.spec.endpoints[0].targets", "--no-headers")
	if err != nil {
		s.Endpoint0Targets = na
	}
	s.Endpoint1DNSName, err = k8s.RunKubectlAndGetOutputE(i.w.t, i.w.k8sOptions, "get", "dnsendpoints.externaldns.k8s.io", "test-gslb", "-o",
		"custom-columns=SERVICESTATUS:.spec.endpoints[1].dnsName", "--no-headers")
	if err != nil {
		s.Endpoint1DNSName = na
	}
	s.Endpoint1Targets, err = k8s.RunKubectlAndGetOutputE(i.w.t, i.w.k8sOptions, "get", "dnsendpoints.externaldns.k8s.io", "test-gslb", "-o",
		"custom-columns=SERVICESTATUS:.spec.endpoints[1].targets", "--no-headers")
	if err != nil {
		s.Endpoint1Targets = na
	}
	return
}

func (s *InstanceStatus) String() string {
	return gopkgstr.ToString(s)
}

func waitForLocalGSLBNew(t *testing.T, host string, port int, expectedResult []string) (output []string, err error) {
	return DoWithRetryWaitingForValueE(
		t,
		"Wait for failover to happen and coredns to pickup new values...",
		100,
		time.Second*1,
		func() ([]string, error) { return dns.Dig("localhost:"+strconv.Itoa(port), host) },
		expectedResult)
}
