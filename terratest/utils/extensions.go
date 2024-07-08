package utils

/*
Copyright 2022 The k8gb Contributors.

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

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"

	"github.com/AbsaOSS/gopkg/dns"
	gopkgstr "github.com/AbsaOSS/gopkg/string"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
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
		digUsingUDP         bool
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
	i                     *Instance
	Annotation            string   `json:"annotation"`
	AppMessage            string   `json:"app-msg"`
	AppRunning            bool     `json:"podinfo-running"`
	AppReplicas           string   `json:"podinfo-replicas"`
	LocalTargets          []string `json:"local-targets-ip"`
	Ingresses             []string `json:"ingress-ip"`
	Dig                   []string `json:"dig-result"`
	CoreDNS               string   `json:"coredns-ip"`
	GslbHealthStatus      string   `json:"gslb-status"`
	Cluster               string   `json:"cluster"`
	Namespace             string   `json:"namespace"`
	EndpointLocalDNSName  string   `json:"ep0-dns-name"`
	EndpointLocalTargets  string   `json:"ep0-dns-targets"`
	EndpointGlobalDNSName string   `json:"ep1-dns-name"`
	EndpointGlobalTargets string   `json:"ep1-dns-targets"`
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
	var err error
	w.settings.ingressResourcePath, err = filepath.Abs(path)
	if err != nil {
		w.error = fmt.Errorf("reading %s; %s", path, err)
	}
	w.settings.ingressName, err = w.getManifestName(w.settings.ingressResourcePath)
	if err != nil {
		w.error = err
	}
	return w
}

// WithGslb
// TODO: consider taking host dynamically
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
	return w
}

func (w *Workflow) WithDigUsingUDP(digUsingUDP bool) *Workflow {
	w.settings.digUsingUDP = digUsingUDP
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
		k8s.WaitUntilNumPodsCreated(w.t, w.k8sOptions, testAppFilter, 1, DefaultRetries, 1*time.Second)
		var testAppPods []corev1.Pod
		testAppPods = k8s.ListPods(w.t, w.k8sOptions, testAppFilter)
		for _, pod := range testAppPods {
			k8s.WaitUntilPodAvailable(w.t, w.k8sOptions, pod.Name, DefaultRetries, 1*time.Second)
		}
		k8s.WaitUntilServiceAvailable(w.t, w.k8sOptions, w.state.testApp.name, DefaultRetries, 1*time.Second)
		w.state.testApp.isRunning = true
	}

	// gslb
	if w.settings.gslbResourcePath != "" {
		if w.settings.ingressResourcePath == "" {
			w.t.Logf("Create ingress %s from %s", w.state.gslb.name, w.settings.gslbResourcePath)
			k8s.KubectlApply(w.t, w.k8sOptions, w.settings.gslbResourcePath)
			k8s.WaitUntilIngressAvailable(w.t, w.k8sOptions, w.state.gslb.name, 100, 5*time.Second)
			ingress := k8s.GetIngress(w.t, w.k8sOptions, w.state.gslb.name)
			require.Equal(w.t, ingress.Name, w.state.gslb.name)
			w.settings.ingressName = w.state.gslb.name
		} else {
			w.t.Logf("Create ingress %s from %s", w.settings.ingressName, w.settings.ingressResourcePath)
			k8s.KubectlApply(w.t, w.k8sOptions, w.settings.ingressResourcePath)
			k8s.WaitUntilIngressAvailable(w.t, w.k8sOptions, w.settings.ingressName, 100, 5*time.Second)
			w.t.Logf("Create gslb %s from %s", w.state.gslb.name, w.settings.gslbResourcePath)
			k8s.KubectlApply(w.t, w.k8sOptions, w.settings.gslbResourcePath)
		}
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

	yamlFile, err := os.ReadFile(path)
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

func (i *Instance) ReapplyIngress(path string) {
	var err error
	i.w.t.Logf("reapplying %s", path)
	i.w.settings.ingressResourcePath = path
	i.w.settings.gslbResourcePath = ""
	i.w.state.gslb.name, err = i.w.getManifestName(i.w.settings.ingressResourcePath)
	require.NoError(i.w.t, err)
	k8s.KubectlApply(i.w.t, i.w.k8sOptions, i.w.settings.ingressResourcePath)
	// modifying inner state.gslb.name and ingress.Name has nothing to do with reading these values dynamically afterwards.
	k8s.WaitUntilIngressAvailable(i.w.t, i.w.k8sOptions, i.w.state.gslb.name, 60, 5*time.Second)
	ingress := k8s.GetIngress(i.w.t, i.w.k8sOptions, i.w.state.gslb.name)
	require.Equal(i.w.t, ingress.Name, i.w.state.gslb.name)
	i.w.settings.ingressName = i.w.state.gslb.name
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
	return waitForLocalGSLBNew(i.w.t, i.w.state.gslb.host, i.w.state.gslb.port, expectedIPs, i.w.settings.digUsingUDP)
}

func (i *Instance) WaitForLocalDNSEndpointExists() error {
	periodic := func() (result bool, err error) {
		lep := i.Resources().GetLocalDNSEndpoint()
		result = len(lep.Spec.Endpoints) > 0
		return result, err
	}
	return tickerWaiter(DefaultRetries, "LocalDNSEndpoint exists:", periodic)
}

func (i *Instance) WaitForExternalDNSEndpointExists() error {
	periodic := func() (result bool, err error) {
		lep := i.Resources().GetK8gbExternalDNSEndpoint()
		result = len(lep.Spec.Endpoints) > 0
		return result, err
	}
	return tickerWaiter(DefaultRetries, "ExternalDNSEndpoint exists:", periodic)
}

func (r *Resources) WaitForExternalDNSEndpointHasTargets(epName string) error {
	periodic := func() (result bool, err error) {
		epx, err := r.GetK8gbExternalDNSEndpoint().GetEndpointByName(epName)
		if err != nil {
			return false, nil
		}
		result = len(epx.Targets) > 0
		return result, nil
	}
	return tickerWaiter(DefaultRetries, "ExternalDNSEndpoint has targets:", periodic)
}

// WaitForLocalDNSEndpointHasTargets waits until LocalDNSEndpoint has expected targets
func (i *Instance) WaitForLocalDNSEndpointHasTargets(expectedIPs []string) error {
	periodic := func() (result bool, err error) {
		lep := i.Resources().GetLocalDNSEndpoint()
		ep, err := lep.GetEndpointByName(i.w.state.gslb.host)
		result = EqualStringSlices(ep.Targets, expectedIPs)
		return result, err
	}
	return tickerWaiter(160, "LocalDNSEndpoint targets:", periodic)
}

// WaitForExpected waits until GSLB dig doesn't return list of expected IP's
func (i *Instance) WaitForExpected(expectedIPs []string) (err error) {
	_, err = waitForLocalGSLBNew(i.w.t, i.w.state.gslb.host, i.w.state.gslb.port, expectedIPs, i.w.settings.digUsingUDP)
	if err != nil {
		fmt.Println(i.GetStatus(fmt.Sprintf("expected IPs: %s", expectedIPs)).String())
	}
	return
}

// WaitForAppIsRunning waits until app has one pod running
func (i *Instance) WaitForAppIsRunning() (err error) {
	return i.waitForApp(func(instances int) bool { return instances > 0 }, false)
}

// WaitForAppIsStopped waits until app has 0 pods running
func (i *Instance) WaitForAppIsStopped() (err error) {
	return i.waitForApp(func(instances int) bool { return instances == 0 }, true)
}

// WaitForAppIsRunning waits until app has one pod running
func (i *Instance) waitForApp(predicate func(instances int) bool, stop bool) (err error) {
	const (
		description = "Wait for app is running"
		maxRetries  = 50
		waitSeconds = 2
	)
	op := "stop"
	if !stop {
		op = "start"
	}
	// first condition is to have replicas
	for n := 0; n < maxRetries; n++ {
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
		if predicate(r) {
			i.w.t.Logf("%s found match: Expected:(%v)", description, r)
			break
		}
		i.w.t.Logf("Application %s is not in expected state. Waiting for %d seconds...", i.w.state.testApp.name, waitSeconds)
		time.Sleep(waitSeconds * time.Second)
	}
	i.w.t.Logf("Wait for ExternalDNSEndpoint %s.%s to be filled by targets %s", i.w.state.gslb.name, i.w.namespace, i.w.state.gslb.host)
	// second conditions
	for n := 0; n < maxRetries/2; n++ {
		ep, err := i.Resources().GetExternalDNSEndpointByName(i.w.state.gslb.name, i.w.namespace).GetEndpointByName(i.w.state.gslb.host)
		if err != nil {
			// app is already stopped and cant be found
			if stop && err.Error() == notFoundError {
				i.w.t.Logf("App is stopped %s", i.w.state.testApp.name)
				return nil
			}
			i.w.t.Logf("Error waiting for the app %s. %s", i.w.state.testApp.name, err)
			require.NoError(i.w.t, err)
		}
		// waiting for the start
		if !stop && len(ep.Targets) == 0 {
			i.w.t.Logf("Waiting for %s to be started. Waiting for %d seconds...", i.w.state.testApp.name, waitSeconds)
			time.Sleep(waitSeconds * time.Second)
			continue
		}
		// waiting for the stop
		if stop && len(ep.Targets) != 0 {
			i.w.t.Logf("Waiting for %s to be stopped. Waiting for %d seconds...", i.w.state.testApp.name, waitSeconds)
			time.Sleep(waitSeconds * time.Second)
			continue
		}
		return nil
	}
	return retry.MaxRetriesExceeded{Description: "Unable to " + op + " Podinfo app", MaxRetries: maxRetries}
}

// String retrieves rough information about cluster
func (i *Instance) String() (out string) {
	return fmt.Sprintf("Instance: %s:%s", i.w.cluster, i.w.namespace)
}

// Dig  returns a list of IP addresses from CoreDNS that belong to the instance
func (i *Instance) Dig() []string {
	dig, err := dns.Dig("localhost:"+strconv.Itoa(i.w.state.gslb.port), i.w.state.gslb.host, i.w.settings.digUsingUDP)
	require.NoError(i.w.t, err)
	return dig
}

// GetLocalTargets returns instance local targets
func (i *Instance) GetLocalTargets() []string {
	dnsName := fmt.Sprintf("localtargets-%s", i.w.state.gslb.host)
	dig, err := dns.Dig("localhost:"+strconv.Itoa(i.w.state.gslb.port), dnsName, i.w.settings.digUsingUDP)
	i.logIfError(err, "GetLocalTargets(), dig: %s", err)
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
	for t := 0; t < 60; t++ {
		result.Body, err = RunBusyBoxCommand(i.w.t, i.w.k8sOptions, coreDNSIP, command)
		if strings.Contains(result.Body, "503") {
			i.w.t.Log("podinfo returns 503, trying again....")
			time.Sleep(time.Second * 1)
			continue
		}
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
	s.EndpointLocalDNSName, err = k8s.RunKubectlAndGetOutputE(i.w.t, i.w.k8sOptions, "get", "dnsendpoints.externaldns.k8s.io", "test-gslb", "-o",
		"custom-columns=SERVICESTATUS:.spec.endpoints[0].dnsName", "--no-headers")
	if err != nil {
		s.EndpointLocalDNSName = na
	}
	s.EndpointLocalTargets, err = k8s.RunKubectlAndGetOutputE(i.w.t, i.w.k8sOptions, "get", "dnsendpoints.externaldns.k8s.io", "test-gslb", "-o",
		"custom-columns=SERVICESTATUS:.spec.endpoints[0].targets", "--no-headers")
	if err != nil {
		s.EndpointLocalTargets = na
	}
	s.EndpointGlobalDNSName, err = k8s.RunKubectlAndGetOutputE(i.w.t, i.w.k8sOptions, "get", "dnsendpoints.externaldns.k8s.io", "test-gslb", "-o",
		"custom-columns=SERVICESTATUS:.spec.endpoints[1].dnsName", "--no-headers")
	if err != nil {
		s.EndpointGlobalDNSName = na
	}
	s.EndpointGlobalTargets, err = k8s.RunKubectlAndGetOutputE(i.w.t, i.w.k8sOptions, "get", "dnsendpoints.externaldns.k8s.io", "test-gslb", "-o",
		"custom-columns=SERVICESTATUS:.spec.endpoints[1].targets", "--no-headers")
	if err != nil {
		s.EndpointGlobalTargets = na
	}
	return
}

func (s *InstanceStatus) String() string {
	return gopkgstr.ToString(s)
}

func waitForLocalGSLBNew(t *testing.T, host string, port int, expectedResult []string, isUdp bool) (output []string, err error) {
	return DoWithRetryWaitingForValueE(
		t,
		"Wait for failover to happen and coredns to pickup new values...",
		DefaultRetries,
		time.Second*1,
		func() ([]string, error) { return dns.Dig("localhost:"+strconv.Itoa(port), host, isUdp) },
		expectedResult)
}

// tickerWaiter periodically executes periodic function
func tickerWaiter(timeoutSeconds int, name string, periodic func() (result bool, err error)) error {
	const interval = 5
	var cycles = timeoutSeconds / interval
	t := 0
	for range time.NewTicker(interval * time.Second).C {
		result, err := periodic()
		if err != nil {
			return fmt.Errorf("%s: targets %v, error: %s", name, result, err)
		}
		if result {
			return nil
		}
		if t >= cycles {
			break
		}
		t++
	}
	return fmt.Errorf("%s: timeout", name)
}

func (i *Instance) Resources() (o *Resources) {
	return &Resources{
		i,
	}
}

func (i *Instance) continueIfK8sResourceNotFound(err error) {
	if err != nil && strings.HasSuffix(err.Error(), "not found") {
		return
	}
	require.NoError(i.w.t, err)
}

type Resources struct {
	i *Instance
}

// GslbSpecProperty returns actual value of one Spec property, e.g: `spec.ingress.rules[0].host`
func (g *Gslb) GslbSpecProperty(specPath string) string {
	actualValue, _ := k8s.RunKubectlAndGetOutputE(g.i.w.t, g.i.w.k8sOptions, "get", "gslb", g.i.w.state.gslb.name,
		"-o", fmt.Sprintf("custom-columns=SERVICESTATUS:%s", specPath), "--no-headers")
	return actualValue
}

func (r *Resources) GetLocalDNSEndpoint() DNSEndpoint {
	ep, err := r.getDNSEndpoint("test-gslb", r.i.w.namespace)
	r.i.continueIfK8sResourceNotFound(err)
	return ep
}

func (r *Resources) GetK8gbExternalDNSEndpoint() DNSEndpoint {
	return r.GetExternalDNSEndpointByName("k8gb-ns-extdns", "k8gb")
}

func (r *Resources) GetExternalDNSEndpointByName(name, namespace string) DNSEndpoint {
	ep, err := r.getDNSEndpoint(name, namespace)
	r.i.continueIfK8sResourceNotFound(err)
	return ep
}

func (i *Instance) logIfError(err error, message string, args ...any) {
	if err != nil {
		i.w.t.Logf(message, args)
	}
}

func (r *Resources) getDNSEndpoint(epName, ns string) (ep DNSEndpoint, err error) {
	ep = DNSEndpoint{}
	opts := k8s.NewKubectlOptions(r.i.w.k8sOptions.ContextName, r.i.w.k8sOptions.ConfigPath, ns)
	j, err := k8s.RunKubectlAndGetOutputE(r.i.w.t, opts, "get", "dnsendpoints.externaldns.k8s.io", epName, "-ojson")
	if err != nil {
		return ep, err
	}
	err = json.Unmarshal([]byte(j), &ep)
	return ep, err
}

type Gslb struct {
	i *Instance
}

type Ingress struct {
	*networkingv1.Ingress
	i *Instance
}

func (r *Resources) Ingress() *Ingress {
	var ing *networkingv1.Ingress
	ing = k8s.GetIngress(r.i.w.t, r.i.w.k8sOptions, r.i.w.settings.ingressName)
	return &Ingress{
		Ingress: ing,
		i:       r.i,
	}
}

func (r *Resources) Gslb() *Gslb {
	return &Gslb{
		i: r.i,
	}
}

func (g *Gslb) GetAnnotations() (a map[string]string) {
	m := struct {
		Metadata struct {
			Annotations map[string]string `yaml:"annotations"`
		} `yaml:"metadata"`
	}{}
	strValue, err := k8s.RunKubectlAndGetOutputE(g.i.w.t, g.i.w.k8sOptions, "get", "gslb", g.i.w.state.gslb.name, "-ojson")
	require.NoError(g.i.w.t, err)
	err = json.Unmarshal([]byte(strValue), &m)
	require.NoError(g.i.w.t, err)
	return m.Metadata.Annotations
}

func (g *Gslb) PatchAnnotations(a map[string]string) (err error) {
	return g.i.patchAnnotations(g.i.w.state.gslb.name, "gslb", a)
}

func (ing *Ingress) PatchAnnotations(a map[string]string) (err error) {
	return ing.i.patchAnnotations(ing.GetName(), "ingress", a)
}

func (i *Instance) patchAnnotations(name, ktype string, a map[string]string) (err error) {
	var data []string
	for k, v := range a {
		data = append(data, fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	annotations := fmt.Sprintf("{\"metadata\":{\"annotations\":{%s}}}", strings.Join(data, ","))
	_, err = k8s.RunKubectlAndGetOutputE(i.w.t, i.w.k8sOptions, "patch", ktype, name, "-p", annotations, "--type=merge")
	return err
}
