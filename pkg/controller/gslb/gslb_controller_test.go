package gslb

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
	yamlConv "github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	zap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var gslbYaml = []byte(`
apiVersion: ohmyglb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: test-gslb
  namespace: test-gslb
spec:
  ingress:
    rules:
    - host: app.cloud.absa.internal
      http:
        paths:
        - backend:
            serviceName: app
            servicePort: http
          path: /
    - host: app1.cloud.absa.external
      http:
        paths:
        - backend:
            serviceName: unhealthy-nginx
            servicePort: http
          path: /
    - host: app2.cloud.absa.external
      http:
        paths:
        - backend:
            serviceName: healthy-nginx
            servicePort: http
          path: /
  strategy: roundRobin
`)

func TestGslbController(t *testing.T) {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(zap.Logger(true))

	gslbName := "test-gslb"
	namespace := "test-gslb"

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
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileGslb object with the scheme and fake client.
	r := &ReconcileGslb{client: cl, scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      gslbName,
			Namespace: namespace,
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

		expectedHosts := []string{"app.cloud.absa.internal", "app1.cloud.absa.external", "app2.cloud.absa.external"}
		actualHosts := gslb.Status.ManagedHosts
		if !reflect.DeepEqual(expectedHosts, actualHosts) {
			t.Errorf("expected %v managed hosts, but got %v", expectedHosts, actualHosts)
		}
	})

	t.Run("NotFound service status", func(t *testing.T) {
		expectedServiceStatus := "NotFound"
		actualServiceStatus := gslb.Status.ServiceHealth["app"]
		if expectedServiceStatus != actualServiceStatus {
			t.Errorf("expected App service status to be %s, but got %s", expectedServiceStatus, actualServiceStatus)
		}
	})

	t.Run("Unhealthy service status", func(t *testing.T) {
		serviceName := "unhealthy-nginx"
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: namespace,
			},
		}

		err = cl.Create(context.TODO(), service)
		if err != nil {
			t.Fatalf("Failed to create testing service: (%v)", err)
		}

		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		expectedServiceStatus := "Unhealthy"
		actualServiceStatus := gslb.Status.ServiceHealth[serviceName]
		if expectedServiceStatus != actualServiceStatus {
			t.Errorf("expected App service status to be %s, but got %s", expectedServiceStatus, actualServiceStatus)
		}
	})

	t.Run("Healthy service status", func(t *testing.T) {
		serviceName := "healthy-nginx"
		labels := map[string]string{"app": "nginx"}
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: namespace,
				Labels:    labels,
			},
		}

		err = cl.Create(context.TODO(), service)
		if err != nil {
			t.Fatalf("Failed to create testing service: (%v)", err)
		}

		endpoint := &corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: namespace,
				Labels:    labels,
			},
		}

		err = cl.Create(context.TODO(), endpoint)
		if err != nil {
			t.Fatalf("Failed to create testing endpoint: (%v)", err)
		}

		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		expectedServiceStatus := "Healthy"
		actualServiceStatus := gslb.Status.ServiceHealth[serviceName]
		if expectedServiceStatus != actualServiceStatus {
			t.Errorf("expected App service status to be %s, but got %s", expectedServiceStatus, actualServiceStatus)
		}
	})

	t.Run("Healthy workers status", func(t *testing.T) {
		nodeName := "test-node"
		nodeAddress := "10.0.0.1"
		nodeLabels := map[string]string{"node-role.kubernetes.io/worker": ""}
		node := &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
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

func YamlToGslb(yaml []byte) (*ohmyglbv1beta1.Gslb, error) {
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
