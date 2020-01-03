package gslb

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"regexp"
	"testing"

	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
	yamlConv "github.com/ghodss/yaml"
	externaldns "github.com/kubernetes-incubator/external-dns/endpoint"
	corev1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	zap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var corednsConfigMapYaml = []byte(`
apiVersion: v1
data:
  Corefile: |-
    .:53 {
        cache 30
        errors
        health
        ready
        kubernetes cluster.local
        loadbalance round_robin
        prometheus 0.0.0.0:9153
        forward . /etc/resolv.conf
        hosts /etc/coredns/gslb.hosts {
            ttl 30
            reload "300ms"
        }
    }
  gslb.hosts: |
    127.0.0.1       localhost
    192.168.1.10    absa.oss.org            absa
kind: ConfigMap
metadata:
  creationTimestamp: null
  labels:
    app.kubernetes.io/instance: gslb-coredns
  name: gslb-coredns-coredns
  namespace: ohmyglb
`)

var crSampleYaml = "../../../deploy/crds/ohmyglb.absa.oss_v1beta1_gslb_cr.yaml"

func TestGslbController(t *testing.T) {
	gslbYaml, err := ioutil.ReadFile(crSampleYaml)
	if err != nil {
		t.Fatalf("Can't open example CR file: %s", crSampleYaml)
	}
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(zap.Logger(true))

	cmName := "gslb-coredns-coredns"

	gslb, err := yamlToGslb(gslbYaml)
	if err != nil {
		t.Fatal(err)
	}

	corednsConfigMap, err := yamlToConfigMap(corednsConfigMapYaml)
	if err != nil {
		t.Fatal(err)
	}

	objs := []runtime.Object{
		gslb,
		corednsConfigMap,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(ohmyglbv1beta1.SchemeGroupVersion, gslb)
	// Register external-dns DNSEndpoint CRD
	s.AddKnownTypes(schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}, &externaldns.DNSEndpoint{})
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileGslb object with the scheme and fake client.
	r := &ReconcileGslb{client: cl, scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      gslb.Name,
			Namespace: gslb.Namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("requeue expected")
	}
	ingress := &v1beta1.Ingress{}
	err = cl.Get(context.TODO(), req.NamespacedName, ingress)
	if err != nil {
		t.Fatalf("Failed to get expected ingress: (%v)", err)
	}

	// Reconcile again so Reconcile() checks services and updates the Gslb
	// resources' Status.
	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	if res != (reconcile.Result{}) {
		t.Error("reconcile did not return an empty Result")
	}

	t.Run("ManagedHosts status", func(t *testing.T) {
		err = cl.Get(context.TODO(), req.NamespacedName, gslb)
		if err != nil {
			t.Fatalf("Failed to get expected gslb: (%v)", err)
		}

		expectedHosts := []string{"app.cloud.absa.internal", "app2.cloud.absa.internal", "app3.cloud.absa.internal"}
		actualHosts := gslb.Status.ManagedHosts
		if !reflect.DeepEqual(expectedHosts, actualHosts) {
			t.Errorf("expected %v managed hosts, but got %v", expectedHosts, actualHosts)
		}
	})

	t.Run("NotFound service status", func(t *testing.T) {
		expectedServiceStatus := "NotFound"
		notFoundHost := "app.cloud.absa.internal"
		actualServiceStatus := gslb.Status.ServiceHealth[notFoundHost]
		if expectedServiceStatus != actualServiceStatus {
			t.Errorf("expected %s service status to be %s, but got %s", notFoundHost, expectedServiceStatus, actualServiceStatus)
		}
	})

	t.Run("Unhealthy service status", func(t *testing.T) {
		serviceName := "unhealthy-nginx"
		unhealthyHost := "app2.cloud.absa.internal"
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: gslb.Namespace,
			},
		}

		err = cl.Create(context.TODO(), service)
		if err != nil {
			t.Fatalf("Failed to create testing service: (%v)", err)
		}

		endpoint := &corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: gslb.Namespace,
			},
		}

		err = cl.Create(context.TODO(), endpoint)
		if err != nil {
			t.Fatalf("Failed to create testing endpoint: (%v)", err)
		}

		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		expectedServiceStatus := "Unhealthy"
		actualServiceStatus := gslb.Status.ServiceHealth[unhealthyHost]
		if expectedServiceStatus != actualServiceStatus {
			t.Errorf("expected %s service status to be %s, but got %s", unhealthyHost, expectedServiceStatus, actualServiceStatus)
		}
	})

	t.Run("Healthy service status", func(t *testing.T) {
		serviceName := "healthy-nginx"
		labels := map[string]string{"app": "nginx"}
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: gslb.Namespace,
				Labels:    labels,
			},
		}

		err = cl.Create(context.TODO(), service)
		if err != nil {
			t.Fatalf("Failed to create testing service: (%v)", err)
		}

		// Create fake endpoint with populated address slice
		endpoint := &corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: gslb.Namespace,
				Labels:    labels,
			},
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{{IP: "1.2.3.4"}},
				},
			},
		}

		err = cl.Create(context.TODO(), endpoint)
		if err != nil {
			t.Fatalf("Failed to create testing endpoint: (%v)", err)
		}

		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		expectedServiceStatus := "Healthy"
		healthyHost := "app3.cloud.absa.internal"
		actualServiceStatus := gslb.Status.ServiceHealth[healthyHost]
		if expectedServiceStatus != actualServiceStatus {
			t.Errorf("expected %s service status to be %s, but got %s", healthyHost, expectedServiceStatus, actualServiceStatus)
		}
	})

	t.Run("Healthy workers status", func(t *testing.T) {
		nodeName := "test-node"
		nodeAddress := "10.0.0.1"
		nodeLabels := map[string]string{"node-role.kubernetes.io/worker": ""}
		node := &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: gslb.Namespace,
				Name:      nodeName,
				Labels:    nodeLabels,
			},
		}
		err = cl.Create(context.TODO(), node)
		if err != nil {
			t.Fatalf("Failed to create testing node: (%v)", err)
		}
		node.Status.Addresses = []corev1.NodeAddress{
			{
				Type:    "InternalIP",
				Address: nodeAddress},
		}
		err = r.client.Status().Update(context.TODO(), node)
		if err != nil {
			t.Fatalf("Failed to update Node status:(%v)", err)
		}

		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		want := make(map[string]string)
		want[nodeName] = nodeAddress

		got := gslb.Status.HealthyWorkers

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected HealthyWorkers status to be %s, but got %s", want, got)
		}
	})

	t.Run("CoreDNS configmap was updated with healthy service references", func(t *testing.T) {

		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		nn := types.NamespacedName{
			Name:      cmName,
			Namespace: gslbOperatorNamespace,
		}

		err = cl.Get(context.TODO(), nn, corednsConfigMap)
		if err != nil {
			t.Fatalf("Failed to get expected corednsConfigMap: (%v)", err)
		}

		want := `10\.0\.0\.1.*app3.cloud.absa.internal`
		got := corednsConfigMap.Data["gslb.hosts"]
		matched, _ := regexp.MatchString(want, got)
		if !matched {
			t.Errorf("Expected to have '%s' in the host configmap '%s'", want, got)
		}
	})

	t.Run("CoreDNS configmap contains no uhealthy service references", func(t *testing.T) {

		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		nn := types.NamespacedName{
			Name:      cmName,
			Namespace: gslbOperatorNamespace,
		}

		err = cl.Get(context.TODO(), nn, corednsConfigMap)
		if err != nil {
			t.Fatalf("Failed to get expected corednsConfigMap: (%v)", err)
		}

		wants := []string{`10\.0\.0\.1.*app.cloud.absa.internal`, `10\.0\.0\.1.*app2.cloud.absa.internal`}
		for _, want := range wants {
			got := corednsConfigMap.Data["gslb.hosts"]
			matched, _ := regexp.MatchString(want, got)
			if matched {
				t.Errorf("Expected to NOT to have non healthy '%s' in the host configmap '%s' ", want, got)
			}
		}
	})

	t.Run("Gslb creates DNSEndpoint CR for healthy ingress hosts", func(t *testing.T) {

		ingressIP := corev1.LoadBalancerIngress{
			IP: "10.0.0.1",
		}
		ingress.Status.LoadBalancer.Ingress = append(ingress.Status.LoadBalancer.Ingress, ingressIP)
		err := cl.Status().Update(context.TODO(), ingress)
		if err != nil {
			t.Fatalf("Failed to update gslb Ingress Address: (%v)", err)
		}

		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		dnsEndpoint := &externaldns.DNSEndpoint{}
		err = cl.Get(context.TODO(), req.NamespacedName, dnsEndpoint)
		if err != nil {
			t.Fatalf("Failed to get expected DNSEndpoint: (%v)", err)
		}

		got := dnsEndpoint.Spec.Endpoints

		want := []*externaldns.Endpoint{{
			DNSName:    "app3.cloud.absa.internal",
			RecordTTL:  30,
			RecordType: "A",
			Targets:    externaldns.Targets{"10.0.0.1"}},
		}

		prettyGot := prettyPrint(got)
		prettyWant := prettyPrint(want)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
		}
	})
}

func reconcileAndUpdateGslb(t *testing.T,
	r *ReconcileGslb,
	req reconcile.Request,
	cl client.Client,
	gslb *ohmyglbv1beta1.Gslb,
) {
	t.Helper()
	// Reconcile again so Reconcile() checks services and updates the Gslb
	// resources' Status.
	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	if res != (reconcile.Result{}) {
		t.Error("reconcile did not return an empty Result")
	}

	err = cl.Get(context.TODO(), req.NamespacedName, gslb)
	if err != nil {
		t.Fatalf("Failed to get expected gslb: (%v)", err)
	}
}

func yamlToGslb(yaml []byte) (*ohmyglbv1beta1.Gslb, error) {
	// yamlBytes contains a []byte of my yaml job spec
	// convert the yaml to json
	jsonBytes, err := yamlConv.YAMLToJSON(yaml)
	if err != nil {
		return &ohmyglbv1beta1.Gslb{}, err
	}
	// unmarshal the json into the kube struct
	var gslb = &ohmyglbv1beta1.Gslb{}
	err = json.Unmarshal(jsonBytes, &gslb)
	if err != nil {
		return &ohmyglbv1beta1.Gslb{}, err
	}
	return gslb, nil
}

func yamlToConfigMap(yaml []byte) (*corev1.ConfigMap, error) {
	// yamlBytes contains a []byte of my yaml job spec
	// convert the yaml to json
	jsonBytes, err := yamlConv.YAMLToJSON(yaml)
	if err != nil {
		return &corev1.ConfigMap{}, err
	}
	// unmarshal the json into the kube struct
	var cm = &corev1.ConfigMap{}
	err = json.Unmarshal(jsonBytes, &cm)
	if err != nil {
		return &corev1.ConfigMap{}, err
	}
	return cm, nil
}
