package e2e

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/AbsaOSS/ohmyglb/pkg/apis"
	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
	ohmyhelpers "github.com/AbsaOSS/ohmyglb/pkg/controller/gslb"
	externaldns "github.com/kubernetes-incubator/external-dns/endpoint"
	"k8s.io/api/extensions/v1beta1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/scheme"

	goctx "context"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
	crSampleYaml         = "./deploy/crds/ohmyglb.absa.oss_v1beta1_gslb_cr.yaml"
)

func TestGslb(t *testing.T) {
	gslbList := &ohmyglbv1beta1.GslbList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, gslbList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	schemeBuilder := &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}}
	schemeBuilder.Register(&externaldns.DNSEndpoint{}, &externaldns.DNSEndpointList{})
	err = framework.AddToFrameworkScheme(schemeBuilder.AddToScheme, &externaldns.DNSEndpointList{})
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()

	err = ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := framework.Global

	// Dynamically create ClusterRoleBinding for dynamically created namespaced gslb ServiceAccount
	// Otherwise tests will fail with rbac restrictions
	clusterRoleBinding := &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
		Subjects: []rbac.Subject{{
			Kind:      "ServiceAccount",
			Name:      "ohmyglb",
			Namespace: namespace,
		}},
		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "ohmyglb",
		},
	}

	err = f.Client.Create(goctx.TODO(), clusterRoleBinding, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatal(err)
	}

	// wait for ohmyglb-operator to be ready
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "ohmyglb", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	gslb, err := createGslb(t, f, ctx)
	if err != nil {
		t.Fatal(err)
	}

	testGslbIngress(t, f, ctx, gslb)
	testGslbDNSEndpoint(t, f, ctx, gslb)
}

func createGslb(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) (*ohmyglbv1beta1.Gslb, error) {

	gslbYaml, err := ioutil.ReadFile(crSampleYaml)
	if err != nil {
		t.Fatalf("Can't open example CR file: %s", crSampleYaml)
	}

	testGslb, err := ohmyhelpers.YamlToGslb(gslbYaml)
	if err != nil {
		t.Fatal(err)
	}

	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatalf("could not get namespace: %v", err)
	}

	testGslb.Namespace = namespace

	err = f.Client.Create(goctx.TODO(), testGslb, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return nil, err
	}

	return testGslb, nil

}

func testGslbIngress(t *testing.T, f *framework.Framework, ctx *framework.TestCtx, gslb *ohmyglbv1beta1.Gslb) {
	t.Run("Gslb creates associated Ingress", func(t *testing.T) {
		ingress := &v1beta1.Ingress{}
		nn := types.NamespacedName{Name: gslb.Name, Namespace: gslb.Namespace}
		err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
			err = f.Client.Get(context.TODO(), nn, ingress)
			if err != nil {
				return false, err
			}
			return true, nil
		})
		if err != nil {
			t.Fatalf("Failed to get expected ingress: (%v)", err)
		}
	})
}

func testGslbDNSEndpoint(t *testing.T, f *framework.Framework, ctx *framework.TestCtx, gslb *ohmyglbv1beta1.Gslb) {
	t.Run("Gslb creates DNSEndpoint", func(t *testing.T) {
		dnsEndpoint := &externaldns.DNSEndpoint{}
		nn := types.NamespacedName{Name: gslb.Name, Namespace: gslb.Namespace}
		err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
			err = f.Client.Get(context.TODO(), nn, dnsEndpoint)
			if err != nil {
				return false, err
			}
			return true, nil
		})
		if err != nil {
			t.Fatalf("Failed to get expected DNSEndpoint: (%v)", err)
		}
	})
}
