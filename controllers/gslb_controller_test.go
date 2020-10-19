package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	"github.com/AbsaOSS/k8gb/controllers/depresolver"
	"github.com/AbsaOSS/k8gb/controllers/internal/utils"
	ibclient "github.com/infobloxopen/infoblox-go-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	corev1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	zap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

var crSampleYaml = "../deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr.yaml"

func TestGslbController(t *testing.T) {
	// Start fakedns server for external dns tests
	fakedns()
	// Isolate the unit tests from interaction with real infoblox grid

	err := os.Setenv("FAKE_INFOBLOX", "true")
	if err != nil {
		t.Fatalf("Can't setup env var: (%v)", err)
	}

	predefinedConfig := depresolver.Config{
		ReconcileRequeueSeconds: 30,
		ClusterGeoTag:           "us-west-1",
		ExtClustersGeoTags:      []string{"us-east-1"},
		EdgeDNSServer:           "8.8.8.8",
		EdgeDNSZone:             "example.com",
		DNSZone:                 "cloud.example.com",
		Route53Enabled:          false,
		Infoblox: depresolver.Infoblox{
			Host:     "fakeinfoblox.example.com",
			Username: "foo",
			Password: "blah",
			Port:     443,
			Version:  "0.0.0",
		},
	}
	configureEnvVar(predefinedConfig)

	_, err = os.Stat(crSampleYaml)
	if os.IsNotExist(err) {
		t.Fatalf("Sample CR yaml file not found at: %s", crSampleYaml)
	}
	gslbYaml, err := ioutil.ReadFile(crSampleYaml)
	if err != nil {
		t.Fatalf("Can't open example CR file: %s", crSampleYaml)
	}
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(zap.Logger(true))

	gslb, err := utils.YamlToGslb(gslbYaml)
	if err != nil {
		t.Fatal(err)
	}

	objs := []runtime.Object{
		gslb,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(k8gbv1beta1.GroupVersion, gslb)
	// Register external-dns DNSEndpoint CRD
	s.AddKnownTypes(schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}, &externaldns.DNSEndpoint{})
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	// Create config
	config, err := depresolver.NewDependencyResolver(context.TODO(), cl).ResolveOperatorConfig()
	if err != nil {
		t.Fatalf("config error: (%v)", err)
	}
	// Create a GslbReconciler object with the scheme and fake client.
	r := &GslbReconciler{
		Client: cl,
		Log:    ctrl.Log.WithName("setup"),
		Scheme: s,
	}
	r.DepResolver = depresolver.NewDependencyResolver(context.TODO(), r.Client)
	r.Config = config
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
	reconcileAndUpdateGslb(t, r, req, cl, gslb)

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

		createUnhealthyService(t, serviceName, cl, gslb)
		defer deleteUnhealthyService(t, serviceName, cl, gslb)
		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		expectedServiceStatus := "Unhealthy"
		actualServiceStatus := gslb.Status.ServiceHealth[unhealthyHost]
		if expectedServiceStatus != actualServiceStatus {
			t.Errorf("expected %s service status to be %s, but got %s", unhealthyHost, expectedServiceStatus, actualServiceStatus)
		}

		reconcileAndUpdateGslb(t, r, req, cl, gslb)

	})

	t.Run("Healthy service status", func(t *testing.T) {
		serviceName := "frontend-podinfo"

		createHealthyService(t, serviceName, cl, gslb)
		defer deleteHealthyService(t, serviceName, cl, gslb)
		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		expectedServiceStatus := "Healthy"
		healthyHost := "app3.cloud.example.com"
		actualServiceStatus := gslb.Status.ServiceHealth[healthyHost]
		if expectedServiceStatus != actualServiceStatus {
			t.Errorf("expected %s service status to be %s, but got %s", healthyHost, expectedServiceStatus, actualServiceStatus)
		}
		reconcileAndUpdateGslb(t, r, req, cl, gslb)

	})

	t.Run("ingress_hosts_per_status metric", func(t *testing.T) {
		err = cl.Get(context.TODO(), req.NamespacedName, gslb)
		if err != nil {
			t.Fatalf("Failed to get expected gslb: (%v)", err)
		}
		expectedHostsMetricCount := 3
		actualHostsMetricCount := testutil.CollectAndCount(ingressHostsPerStatusMetric)
		if !reflect.DeepEqual(expectedHostsMetricCount, actualHostsMetricCount) {
			t.Errorf("expected %v managed hosts, but got %v", expectedHostsMetricCount, actualHostsMetricCount)
		}
	})

	t.Run("ingress_hosts_per_status metric reflection for Healthy status", func(t *testing.T) {
		err = cl.Get(context.TODO(), req.NamespacedName, gslb)
		if err != nil {
			t.Fatalf("Failed to get expected gslb: (%v)", err)
		}
		healthyHosts := ingressHostsPerStatusMetric.With(prometheus.Labels{"namespace": gslb.Namespace, "name": gslb.Name, "status": healthyStatus})
		expectedHostsMetric := 1.
		actualHostsMetric := testutil.ToFloat64(healthyHosts)
		if !reflect.DeepEqual(expectedHostsMetric, actualHostsMetric) {
			t.Errorf("expected %v managed hosts with Healthy status, but got %v", expectedHostsMetric, actualHostsMetric)
		}

		serviceName := "frontend-podinfo"
		createHealthyService(t, serviceName, cl, gslb)
		defer deleteHealthyService(t, serviceName, cl, gslb)
		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		healthyHosts = ingressHostsPerStatusMetric.With(prometheus.Labels{"namespace": gslb.Namespace, "name": gslb.Name, "status": healthyStatus})
		expectedHostsMetric = 1.0
		actualHostsMetric = testutil.ToFloat64(healthyHosts)
		if !reflect.DeepEqual(expectedHostsMetric, actualHostsMetric) {
			t.Errorf("expected %v managed hosts with Healthy status, but got %v", expectedHostsMetric, actualHostsMetric)
		}
		reconcileAndUpdateGslb(t, r, req, cl, gslb)
	})

	t.Run("ingress_hosts_per_status metric reflection for Unhealthy status", func(t *testing.T) {
		err = cl.Get(context.TODO(), req.NamespacedName, gslb)
		if err != nil {
			t.Fatalf("Failed to get expected gslb: (%v)", err)
		}
		unhealthyHosts := ingressHostsPerStatusMetric.With(prometheus.Labels{"namespace": gslb.Namespace, "name": gslb.Name, "status": unhealthyStatus})
		expectedHostsMetricCount := 0.0
		actualHostsMetricCount := testutil.ToFloat64(unhealthyHosts)
		if !reflect.DeepEqual(expectedHostsMetricCount, actualHostsMetricCount) {
			t.Errorf("expected %v managed hosts, but got %v", expectedHostsMetricCount, actualHostsMetricCount)
		}

		serviceName := "unhealthy-app"
		createUnhealthyService(t, serviceName, cl, gslb)
		defer deleteUnhealthyService(t, serviceName, cl, gslb)
		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		unhealthyHosts = ingressHostsPerStatusMetric.With(prometheus.Labels{"namespace": gslb.Namespace, "name": gslb.Name, "status": unhealthyStatus})
		expectedHostsMetricCount = 1.0
		actualHostsMetricCount = testutil.ToFloat64(unhealthyHosts)
		if !reflect.DeepEqual(expectedHostsMetricCount, actualHostsMetricCount) {
			t.Errorf("expected %v managed hosts with Unhealthy status, but got %v", expectedHostsMetricCount, actualHostsMetricCount)
		}
		reconcileAndUpdateGslb(t, r, req, cl, gslb)
	})

	t.Run("ingress_hosts_per_status metric reflection for NotFound status", func(t *testing.T) {
		err = cl.Get(context.TODO(), req.NamespacedName, gslb)
		if err != nil {
			t.Fatalf("Failed to get expected gslb: (%v)", err)
		}
		unknownHosts, _ := ingressHostsPerStatusMetric.GetMetricWith(prometheus.Labels{"namespace": gslb.Namespace, "name": gslb.Name, "status": notFoundStatus})
		expectedHostsMetricCount := 2.0
		actualHostsMetricCount := testutil.ToFloat64(unknownHosts)
		if !reflect.DeepEqual(expectedHostsMetricCount, actualHostsMetricCount) {
			t.Errorf("expected %v managed hosts with NotFound status, but got %v", expectedHostsMetricCount, actualHostsMetricCount)
		}
	})

	t.Run("healthy_records metric", func(t *testing.T) {
		err = cl.Get(context.TODO(), req.NamespacedName, gslb)
		if err != nil {
			t.Fatalf("Failed to get expected gslb: (%v)", err)
		}
		serviceName := "frontend-podinfo"
		createHealthyService(t, serviceName, cl, gslb)
		defer deleteHealthyService(t, serviceName, cl, gslb)
		ingressIPs := []corev1.LoadBalancerIngress{
			{IP: "10.0.0.1"},
			{IP: "10.0.0.2"},
			{IP: "10.0.0.3"},
		}
		err = cl.Get(context.TODO(), req.NamespacedName, ingress)
		if err != nil {
			t.Fatalf("Failed to get expected ingress: (%v)", err)
		}
		ingress.Status.LoadBalancer.Ingress = append(ingress.Status.LoadBalancer.Ingress, ingressIPs...)
		err := cl.Status().Update(context.TODO(), ingress)
		if err != nil {
			t.Fatalf("Failed to update gslb Ingress Address: (%v)", err)
		}
		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		expectedHealthyRecordsMetricCount := 3.0
		actualHealthyRecordsMetricCount := testutil.ToFloat64(healthyRecordsMetric)
		if !reflect.DeepEqual(expectedHealthyRecordsMetricCount, actualHealthyRecordsMetricCount) {
			t.Errorf("expected %v healthy records, but got %v", expectedHealthyRecordsMetricCount, actualHealthyRecordsMetricCount)
		}
		reconcileAndUpdateGslb(t, r, req, cl, gslb)
		ingress.Status.LoadBalancer.Ingress = nil
	})

	scenarios := map[string]prometheus.Collector{
		"healthy_records":          healthyRecordsMetric,
		"ingress_hosts_per_status": ingressHostsPerStatusMetric,
	}

	for name, scenario := range scenarios {
		t.Run(name+" metric linter check", func(t *testing.T) {
			defer func() {
				lintErrors, _ := testutil.CollectAndLint(scenario)
				if len(lintErrors) > 0 {
					t.Errorf("Metric linting error(s): %s", lintErrors)
				}
			}()
		})
	}

	t.Run("Gslb creates DNSEndpoint CR for healthy ingress hosts", func(t *testing.T) {
		serviceName := "frontend-podinfo"

		createHealthyService(t, serviceName, cl, gslb)
		defer deleteHealthyService(t, serviceName, cl, gslb)
		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		dnsEndpoint := &externaldns.DNSEndpoint{}
		err = cl.Get(context.TODO(), req.NamespacedName, dnsEndpoint)
		if err != nil {
			t.Fatalf("Failed to get expected DNSEndpoint: (%v)", err)
		}

		got := dnsEndpoint.Spec.Endpoints

		want := []*externaldns.Endpoint{
			{
				DNSName:    "localtargets-app3.cloud.example.com",
				RecordTTL:  30,
				RecordType: "A",
				Targets:    externaldns.Targets{"10.0.0.1", "10.0.0.2", "10.0.0.3"}},
			{
				DNSName:    "app3.cloud.example.com",
				RecordTTL:  30,
				RecordType: "A",
				Targets:    externaldns.Targets{"10.0.0.1", "10.0.0.2", "10.0.0.3"}},
		}

		prettyGot := prettyPrint(got)
		prettyWant := prettyPrint(want)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
		}
		reconcileAndUpdateGslb(t, r, req, cl, gslb)
	})

	// Test is dependant on fixtures created in other tests which is
	// kind of antipattern. OTOH we avoid a lot of fixture creation
	// code so I will keep it this way for a time being
	t.Run("DNS Record reflection in status", func(t *testing.T) {
		got := gslb.Status.HealthyRecords
		want := map[string][]string{"app3.cloud.example.com": {"10.0.0.1", "10.0.0.2", "10.0.0.3"}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got:\n %s healthyRecordsMetric status,\n\n want:\n %s", got, want)
		}
	})

	t.Run("Local DNS records has special annotation", func(t *testing.T) {
		dnsEndpoint := &externaldns.DNSEndpoint{}
		err = cl.Get(context.TODO(), req.NamespacedName, dnsEndpoint)
		if err != nil {
			t.Fatalf("Failed to get expected DNSEndpoint: (%v)", err)
		}

		got := dnsEndpoint.Annotations["k8gb.absa.oss/dnstype"]

		want := "local"
		if got != want {
			t.Errorf("got:\n %q annotation value,\n\n want:\n %q", got, want)
		}
	})

	t.Run("Generates proper external NS target FQDNs according to the geo tags", func(t *testing.T) {
		err := os.Setenv(depresolver.EdgeDNSZoneKey, "example.com")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}
		err = os.Setenv(depresolver.ExtClustersGeoTagsKey, "za")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}

		got := getExternalClusterFQDNs(gslb)

		want := []string{"gslb-ns-za.example.com"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got:\n %q externalGslb NS records,\n\n want:\n %q", got, want)
		}
	})

	t.Run("Can get external targets from k8gb in another location", func(t *testing.T) {
		serviceName := "frontend-podinfo"
		createHealthyService(t, serviceName, cl, gslb)
		defer deleteHealthyService(t, serviceName, cl, gslb)
		reconcileAndUpdateGslb(t, r, req, cl, gslb)
		err := os.Setenv("OVERRIDE_WITH_FAKE_EXT_DNS", "true")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
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
				DNSName:    "localtargets-app3.cloud.example.com",
				RecordTTL:  30,
				RecordType: "A",
				Targets:    externaldns.Targets{"10.0.0.1", "10.0.0.2", "10.0.0.3"}},
			{
				DNSName:    "app3.cloud.example.com",
				RecordTTL:  30,
				RecordType: "A",
				Targets:    externaldns.Targets{"10.0.0.1", "10.0.0.2", "10.0.0.3", "10.1.0.1", "10.1.0.2", "10.1.0.3"}},
		}

		prettyGot := prettyPrint(got)
		prettyWant := prettyPrint(want)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
		}

		hrGot := gslb.Status.HealthyRecords
		hrWant := map[string][]string{"app3.cloud.example.com": {"10.0.0.1", "10.0.0.2", "10.0.0.3", "10.1.0.1", "10.1.0.2", "10.1.0.3"}}

		if !reflect.DeepEqual(hrGot, hrWant) {
			t.Errorf("got:\n %s Gslb Records status,\n\n want:\n %s", hrGot, hrWant)
		}

		err = os.Setenv("OVERRIDE_WITH_FAKE_EXT_DNS", "false")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}
		reconcileAndUpdateGslb(t, r, req, cl, gslb)
	})

	t.Run("Can check external Gslb TXT record for validity and fail if it is expired", func(t *testing.T) {
		err = os.Setenv("OVERRIDE_WITH_FAKE_EXT_DNS", "true")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}

		got := checkAliveFromTXT("fake", "test-gslb-heartbeat-eu.example.com", time.Minute*5)

		want := errors.NewGone("Split brain TXT record expired the time threshold: (5m0s)")

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got:\n %s from TXT split brain check,\n\n want error:\n %v", got, want)
		}

	})

	t.Run("Can check external Gslb TXT record for validity and pass if it is not expired", func(t *testing.T) {
		err = os.Setenv("OVERRIDE_WITH_FAKE_EXT_DNS", "true")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}

		err := checkAliveFromTXT("fake", "test-gslb-heartbeat-za.example.com", time.Minute*5)

		if err != nil {
			t.Errorf("got:\n %s from TXT split brain check,\n\n want error:\n %v", err, nil)
		}

	})

	t.Run("Can filter out delegated zone entry according FQDN provided", func(t *testing.T) {
		extClusters := getExternalClusterFQDNs(gslb)

		delegateTo := []ibclient.NameServer{
			{Address: "10.0.0.1", Name: "gslb-ns-eu.example.com"},
			{Address: "10.0.0.2", Name: "gslb-ns-eu.example.com"},
			{Address: "10.0.0.3", Name: "gslb-ns-eu.example.com"},
			{Address: "10.1.0.1", Name: "gslb-ns-za.example.com"},
			{Address: "10.1.0.2", Name: "gslb-ns-za.example.com"},
			{Address: "10.1.0.3", Name: "gslb-ns-za.example.com"},
		}

		got := filterOutDelegateTo(delegateTo, extClusters[0])

		want := []ibclient.NameServer{
			{Address: "10.0.0.1", Name: "gslb-ns-eu.example.com"},
			{Address: "10.0.0.2", Name: "gslb-ns-eu.example.com"},
			{Address: "10.0.0.3", Name: "gslb-ns-eu.example.com"},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got:\n %q filtered out delegation records,\n\n want:\n %q", got, want)
		}
	})

	t.Run("Can generate external heartbeat FQDNs", func(t *testing.T) {
		got := getExternalClusterHeartbeatFQDNs(gslb)
		want := []string{"test-gslb-heartbeat-za.example.com"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got:\n %s unexpected heartbeat records,\n\n want:\n %s", got, want)
		}
	})

	t.Run("Returns own records using Failover strategy when Primary", func(t *testing.T) {
		serviceName := "frontend-podinfo"
		createHealthyService(t, serviceName, cl, gslb)
		defer deleteHealthyService(t, serviceName, cl, gslb)
		reconcileAndUpdateGslb(t, r, req, cl, gslb)
		err := os.Setenv("OVERRIDE_WITH_FAKE_EXT_DNS", "true")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}
		err = os.Setenv(depresolver.ClusterGeoTagKey, "eu")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}

		// Enable failover strategy
		gslb.Spec.Strategy.Type = "failover"
		gslb.Spec.Strategy.PrimaryGeoTag = "eu"
		err = cl.Update(context.TODO(), gslb)
		if err != nil {
			t.Fatalf("Can't update gslb: (%v)", err)
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
				DNSName:    "localtargets-app3.cloud.example.com",
				RecordTTL:  30,
				RecordType: "A",
				Targets:    externaldns.Targets{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
			},
			{
				DNSName:    "app3.cloud.example.com",
				RecordTTL:  30,
				RecordType: "A",
				Targets:    externaldns.Targets{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
			},
		}

		prettyGot := prettyPrint(got)
		prettyWant := prettyPrint(want)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
		}

		err = os.Setenv("OVERRIDE_WITH_FAKE_EXT_DNS", "false")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}
		reconcileAndUpdateGslb(t, r, req, cl, gslb)
	})

	t.Run("Returns external records using Failover strategy when Secondary", func(t *testing.T) {
		serviceName := "frontend-podinfo"
		createHealthyService(t, serviceName, cl, gslb)
		defer deleteHealthyService(t, serviceName, cl, gslb)
		reconcileAndUpdateGslb(t, r, req, cl, gslb)
		err := os.Setenv("OVERRIDE_WITH_FAKE_EXT_DNS", "true")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}
		err = os.Setenv(depresolver.ClusterGeoTagKey, "za")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}

		// Enable failover strategy
		gslb.Spec.Strategy.Type = "failover"
		gslb.Spec.Strategy.PrimaryGeoTag = "eu"
		err = cl.Update(context.TODO(), gslb)
		if err != nil {
			t.Fatalf("Can't update gslb: (%v)", err)
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
				DNSName:    "localtargets-app3.cloud.example.com",
				RecordTTL:  30,
				RecordType: "A",
				Targets:    externaldns.Targets{"10.0.0.1", "10.0.0.2", "10.0.0.3"},
			},
			{
				DNSName:    "app3.cloud.example.com",
				RecordTTL:  30,
				RecordType: "A",
				Targets:    externaldns.Targets{"10.1.0.1", "10.1.0.2", "10.1.0.3"},
			},
		}

		prettyGot := prettyPrint(got)
		prettyWant := prettyPrint(want)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
		}

		err = os.Setenv("OVERRIDE_WITH_FAKE_EXT_DNS", "false")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}
		reconcileAndUpdateGslb(t, r, req, cl, gslb)

	})

	t.Run("Gslb properly propagates annotation down to Ingress", func(t *testing.T) {

		expectedAnnotations := map[string]string{"annotation": "test"}
		gslb.Annotations = expectedAnnotations
		err = cl.Update(context.TODO(), gslb)
		if err != nil {
			t.Fatalf("Can't update gslb: (%v)", err)
		}

		reconcileAndUpdateGslb(t, r, req, cl, gslb)

		err = cl.Get(context.TODO(), req.NamespacedName, ingress)
		if err != nil {
			t.Fatalf("Failed to get expected ingress: (%v)", err)
		}

		if !reflect.DeepEqual(ingress.Annotations, expectedAnnotations) {
			t.Errorf("got:\n %s Gslb ingress annotations,\n\n want:\n %s", ingress.Annotations, expectedAnnotations)
		}
	})

	t.Run("Reflect GeoTag in the Status as unset by default", func(t *testing.T) {
		reconcileAndUpdateGslb(t, r, req, cl, gslb)
		got := gslb.Status.GeoTag
		want := "us-west-1"

		if got != want {
			t.Errorf("got: '%s' GeoTag status, want:'%s'", got, want)
		}
	})

	t.Run("Reflect GeoTag in the Status", func(t *testing.T) {
		defer func() {
			err = os.Unsetenv(depresolver.ClusterGeoTagKey)
			if err != nil {
				t.Fatalf("Can't unset env var: (%v)", err)
			}
		}()

		err = os.Setenv(depresolver.ClusterGeoTagKey, "eu")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}
		resolver := depresolver.NewDependencyResolver(context.TODO(), cl)
		r.Config, err = resolver.ResolveOperatorConfig()
		if err != nil {
			t.Fatalf("config error: (%v)", err)
		}
		reconcileAndUpdateGslb(t, r, req, cl, gslb)
		got := gslb.Status.GeoTag
		want := "eu"

		if got != want {
			t.Errorf("got: '%s' GeoTag status, want:'%s'", got, want)
		}
	})

	t.Run("Detects Ingress hostname mismatch", func(t *testing.T) {
		defer func() {
			err := os.Setenv(depresolver.EdgeDNSZoneKey, "example.com")
			if err != nil {
				t.Fatalf("Can't set env var: (%v)", err)
			}
		}()
		err := os.Setenv(depresolver.EdgeDNSZoneKey, "otherdnszone.com")
		if err != nil {
			t.Fatalf("Can't set env var: (%v)", err)
		}
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      gslb.Name,
				Namespace: gslb.Namespace,
			},
		}

		_, err = r.Reconcile(req)
		log.Info(fmt.Sprintf("got an error from controller: %s", err))
		if err == nil {
			t.Errorf("expected controller to detect Ingress hostname and edgeDNSZone mismatch")
		}
	})

	t.Run("Creates NS DNS records for route53", func(t *testing.T) {
		err = os.Setenv(depresolver.EdgeDNSServerKey, "1.1.1.1")
		if err != nil {
			t.Fatalf("Can't setup env var: (%v)", err)
		}
		coreDNSLBServiceName := "k8gb-coredns-lb"
		coreDNSService := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      coreDNSLBServiceName,
				Namespace: k8gbNamespace,
			},
		}
		err = cl.Create(context.TODO(), coreDNSService)
		if err != nil {
			t.Fatalf("Failed to create testing %s service: (%v)", coreDNSLBServiceName, err)
		}
		serviceIPs := []corev1.LoadBalancerIngress{
			{Hostname: "one.one.one.one"}, // rely on 1.1.1.1 response from Cloudflare
		}
		coreDNSService.Status.LoadBalancer.Ingress = append(coreDNSService.Status.LoadBalancer.Ingress, serviceIPs...)
		err = cl.Status().Update(context.TODO(), coreDNSService)
		if err != nil {
			t.Fatalf("Failed to update coredns service lb hostname: (%v)", err)
		}

		err := os.Setenv(depresolver.Route53EnabledKey, "true")
		if err != nil {
			t.Fatalf("Can't set env var: (%v)", err)
		}
		err = os.Setenv(depresolver.ClusterGeoTagKey, "eu")
		if err != nil {
			t.Fatalf("Can't set env var: (%v)", err)
		}
		err = os.Setenv(depresolver.ExtClustersGeoTagsKey, "za,us")
		if err != nil {
			t.Fatalf("Can't set env var: (%v)", err)
		}
		err = os.Setenv(depresolver.DNSZoneKey, "cloud.example.com")
		if err != nil {
			t.Fatalf("Can't set env var: (%v)", err)
		}
		resolver := depresolver.NewDependencyResolver(context.TODO(), cl)
		r.Config, err = resolver.ResolveOperatorConfig()
		if err != nil {
			t.Fatalf("config error: (%v)", err)
		}

		reconcileAndUpdateGslb(t, r, req, cl, gslb)
		dnsEndpointRoute53 := &externaldns.DNSEndpoint{}
		err = cl.Get(context.TODO(), client.ObjectKey{Namespace: k8gbNamespace, Name: "k8gb-ns-route53"}, dnsEndpointRoute53)
		if err != nil {
			t.Fatalf("Failed to get expected DNSEndpoint: (%v)", err)
		}

		got := dnsEndpointRoute53.Annotations["k8gb.absa.oss/dnstype"]

		want := "route53"
		if got != want {
			t.Errorf("got:\n %q annotation value,\n\n want:\n %q", got, want)
		}

		gotEp := dnsEndpointRoute53.Spec.Endpoints

		wantEp := []*externaldns.Endpoint{
			{
				DNSName:    os.Getenv(depresolver.DNSZoneKey),
				RecordTTL:  30,
				RecordType: "NS",
				Targets: externaldns.Targets{
					"gslb-ns-eu.example.com",
					"gslb-ns-us.example.com",
					"gslb-ns-za.example.com",
				},
			},
			{
				DNSName:    "gslb-ns-eu.example.com",
				RecordTTL:  30,
				RecordType: "A",
				Targets: externaldns.Targets{
					"1.0.0.1",
					"1.1.1.1",
				},
			},
		}

		prettyGot := prettyPrint(gotEp)
		prettyWant := prettyPrint(wantEp)

		if !reflect.DeepEqual(gotEp, wantEp) {
			t.Errorf("got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
		}
	})

	t.Run("Resolves LoadBalancer hostname from Ingress status", func(t *testing.T) {
		serviceName := "frontend-podinfo"
		createHealthyService(t, serviceName, cl, gslb)
		defer deleteHealthyService(t, serviceName, cl, gslb)
		err = cl.Get(context.TODO(), req.NamespacedName, ingress)
		if err != nil {
			t.Fatalf("Failed to get expected ingress: (%v)", err)
		}
		ingress.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{Hostname: "one.one.one.one"}}
		err := cl.Status().Update(context.TODO(), ingress)
		if err != nil {
			t.Fatalf("Failed to update gslb Ingress Address: (%v)", err)
		}
		dnsEndpoint := &externaldns.DNSEndpoint{ObjectMeta: metav1.ObjectMeta{Namespace: gslb.Namespace, Name: gslb.Name}}
		err = cl.Delete(context.Background(), dnsEndpoint)
		if err != nil {
			t.Fatalf("Failed to update DNSEndpoint: (%v)", err)
		}

		reconcileAndUpdateGslb(t, r, req, cl, gslb)
		err = cl.Get(context.TODO(), req.NamespacedName, dnsEndpoint)
		if err != nil {
			t.Fatalf("Failed to get expected DNSEndpoint: (%v)", err)
		}

		got := dnsEndpoint.Spec.Endpoints

		want := []*externaldns.Endpoint{
			{
				DNSName:    "localtargets-app3.cloud.example.com",
				RecordTTL:  30,
				RecordType: "A",
				Targets:    externaldns.Targets{"1.0.0.1", "1.1.1.1"}},
			{
				DNSName:    "app3.cloud.example.com",
				RecordTTL:  30,
				RecordType: "A",
				Targets:    externaldns.Targets{"1.0.0.1", "1.1.1.1"}},
		}

		prettyGot := prettyPrint(got)
		prettyWant := prettyPrint(want)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
		}
	})

}

func createHealthyService(t *testing.T, serviceName string, cl client.Client, gslb *k8gbv1beta1.Gslb) {
	t.Helper()
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: gslb.Namespace,
		},
	}
	err := cl.Create(context.TODO(), service)
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
}

func deleteHealthyService(t *testing.T, serviceName string, cl client.Client, gslb *k8gbv1beta1.Gslb) {
	t.Helper()
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: gslb.Namespace,
		},
	}
	err := cl.Delete(context.TODO(), service)
	if err != nil {
		t.Fatalf("Failed to delete testing service: (%v)", err)
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

	err = cl.Delete(context.TODO(), endpoint)
	if err != nil {
		t.Fatalf("Failed to delete testing endpoint: (%v)", err)
	}
}

func createUnhealthyService(t *testing.T, serviceName string, cl client.Client, gslb *k8gbv1beta1.Gslb) {
	t.Helper()
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: gslb.Namespace,
		},
	}

	err := cl.Create(context.TODO(), service)
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

}

func deleteUnhealthyService(t *testing.T, serviceName string, cl client.Client, gslb *k8gbv1beta1.Gslb) {
	t.Helper()
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: gslb.Namespace,
		},
	}

	err := cl.Delete(context.TODO(), service)
	if err != nil {
		t.Fatalf("Failed to delete testing service: (%v)", err)
	}

	endpoint := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: gslb.Namespace,
		},
	}

	err = cl.Delete(context.TODO(), endpoint)
	if err != nil {
		t.Fatalf("Failed to delete testing endpoint: (%v)", err)
	}

}

func reconcileAndUpdateGslb(t *testing.T,
	r *GslbReconciler,
	req reconcile.Request,
	cl client.Client,
	gslb *k8gbv1beta1.Gslb,
) {
	t.Helper()
	// Reconcile again so Reconcile() checks services and updates the Gslb
	// resources' Status.
	res, err := r.Reconcile(req)
	if err != nil {
		return
	}
	if res != (reconcile.Result{RequeueAfter: time.Second * 30}) {
		t.Error("reconcile did not return Result with Requeue")
	}

	err = cl.Get(context.TODO(), req.NamespacedName, gslb)
	if err != nil {
		t.Fatalf("Failed to get expected gslb: (%v)", err)
	}
}

func configureEnvVar(config depresolver.Config) {
	_ = os.Setenv(depresolver.ReconcileRequeueSecondsKey, strconv.Itoa(config.ReconcileRequeueSeconds))
	_ = os.Setenv(depresolver.ClusterGeoTagKey, config.ClusterGeoTag)
	_ = os.Setenv(depresolver.ExtClustersGeoTagsKey, strings.Join(config.ExtClustersGeoTags, ","))
	_ = os.Setenv(depresolver.EdgeDNSServerKey, config.EdgeDNSServer)
	_ = os.Setenv(depresolver.EdgeDNSZoneKey, config.EdgeDNSZone)
	_ = os.Setenv(depresolver.DNSZoneKey, config.DNSZone)
	_ = os.Setenv(depresolver.Route53EnabledKey, strconv.FormatBool(config.Route53Enabled))
	_ = os.Setenv(depresolver.InfobloxGridHostKey, config.Infoblox.Host)
	_ = os.Setenv(depresolver.InfobloxVersionKey, config.Infoblox.Version)
	_ = os.Setenv(depresolver.InfobloxPortKey, strconv.Itoa(config.Infoblox.Port))
	_ = os.Setenv(depresolver.InfobloxUsernameKey, config.Infoblox.Username)
	_ = os.Setenv(depresolver.InfobloxPasswordKey, config.Infoblox.Password)
}
