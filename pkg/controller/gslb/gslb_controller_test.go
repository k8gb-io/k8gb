package gslb

import (
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
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

var crSampleYaml = "../../../deploy/crds/ohmyglb.absa.oss_v1beta1_gslb_cr.yaml"

func TestGslbController(t *testing.T) {
	err := os.Setenv("DNS_ZONE", "example.com")
	if err != nil {
		t.Fatalf("Can't setup env var: (%v)", err)
	}
	gslbYaml, err := ioutil.ReadFile(crSampleYaml)
	if err != nil {
		t.Fatalf("Can't open example CR file: %s", crSampleYaml)
	}
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(zap.Logger(true))

	gslb, err := YamlToGslb(gslbYaml)
	if err != nil {
		t.Fatal(err)
	}

	objs := []runtime.Object{
		gslb,
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

		expectedHosts := []string{"app1.cloud.example.com", "app2.cloud.example.com", "app3.cloud.example.com"}
		actualHosts := gslb.Status.ManagedHosts
		if !reflect.DeepEqual(expectedHosts, actualHosts) {
			t.Errorf("expected %v managed hosts, but got %v", expectedHosts, actualHosts)
		}
	})

	t.Run("NotFound service status", func(t *testing.T) {
		expectedServiceStatus := "NotFound"
		notFoundHost := "app1.cloud.example.com"
		actualServiceStatus := gslb.Status.ServiceHealth[notFoundHost]
		if expectedServiceStatus != actualServiceStatus {
			t.Errorf("expected %s service status to be %s, but got %s", notFoundHost, expectedServiceStatus, actualServiceStatus)
		}
	})

	t.Run("Unhealthy service status", func(t *testing.T) {
		serviceName := "unhealthy-app"
		unhealthyHost := "app2.cloud.example.com"
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
		serviceName := "healthy-app"
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

		// Create fake endpoint with populated address slice
		endpoint := &corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: gslb.Namespace,
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
		healthyHost := "app3.cloud.example.com"
		actualServiceStatus := gslb.Status.ServiceHealth[healthyHost]
		if expectedServiceStatus != actualServiceStatus {
			t.Errorf("expected %s service status to be %s, but got %s", healthyHost, expectedServiceStatus, actualServiceStatus)
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

		want := []*externaldns.Endpoint{
			{
				DNSName:    "hostsz.test-gslb.example.com",
				RecordTTL:  30,
				RecordType: "A",
				Targets:    externaldns.Targets{"10.0.0.1"}},
			{
				DNSName:    "app3.cloud.example.com",
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

	// Test is dependant on fixtures created in other tests which is
	// kind of antipattern. OTOH we avoid a lot of fixture creation
	// code so I will keep it this way for a time being
	t.Run("DNS Record reflection in status", func(t *testing.T) {
		got := gslb.Status.HealthyRecords
		want := map[string][]string{"app3.cloud.example.com": {"10.0.0.1"}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got:\n %s healthyRecords status,\n\n want:\n %s", got, want)
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
