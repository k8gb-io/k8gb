package gslb

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
	yamlConv "github.com/ghodss/yaml"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
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
    - host: app.cloud.absa.external
      http:
        paths:
        - backend:
            serviceName: nginx
            servicePort: http
          path: /
  strategy: roundRobin
`)

func TestGslbController(t *testing.T) {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

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
			Name:      "test-gslb",
			Namespace: "test-gslb",
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

	// Reconcile again so Reconcile() checks pods and updates the Gslb
	// resources' Status.
	res, err = r.Reconcile(req)
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

	expectedHosts := []string{"app.cloud.absa.internal", "app.cloud.absa.external"}
	actualHosts := gslb.Status.ManagedHosts
	if !reflect.DeepEqual(expectedHosts, actualHosts) {
		t.Errorf("expected %v managed hosts, but got %v", expectedHosts, actualHosts)
	}

	expectedServiceStatus := "NotFound"
	actualServiceStatus := gslb.Status.ServiceHealth["app"]
	if expectedServiceStatus != actualServiceStatus {
		t.Errorf("expected App service status to be %s, but got %s", expectedServiceStatus, actualServiceStatus)
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
