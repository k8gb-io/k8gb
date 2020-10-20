package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ibclient "github.com/infobloxopen/infoblox-go-client"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	"github.com/AbsaOSS/k8gb/controllers/depresolver"
	"github.com/AbsaOSS/k8gb/controllers/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

type testSettings struct {
	gslb       *k8gbv1beta1.Gslb
	reconciler *GslbReconciler
	request    reconcile.Request
	config     depresolver.Config
	client     client.Client
	ingress    *v1beta1.Ingress
}

var crSampleYaml = "../deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr.yaml"

var predefinedConfig = depresolver.Config{
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

func TestNotFoundServiceStatus(t *testing.T) {
	// "NotFound service status" independent
	// arrange
	defer cleanup()
	settings := provideSettings(t, predefinedConfig)
	expectedServiceStatus := "NotFound"
	notFoundHost := "app1.cloud.example.com"
	// act
	actualServiceStatus := settings.gslb.Status.ServiceHealth[notFoundHost]
	// assert
	assert.Equal(t, expectedServiceStatus, actualServiceStatus, "expected %s service status to be %s, but got %s", notFoundHost, expectedServiceStatus, actualServiceStatus)
}

func TestUnhealthyServiceStatus(t *testing.T) {
	// "Unhealthy service status" independent
	// arrange
	defer cleanup()
	settings := provideSettings(t, predefinedConfig)
	serviceName := "unhealthy-app"
	unhealthyHost := "app2.cloud.example.com"
	expectedServiceStatus := "Unhealthy"
	defer deleteUnhealthyService(t, serviceName, &settings)
	// act
	createUnhealthyService(t, serviceName, &settings)
	reconcileAndUpdateGslb(t, settings)
	// assert
	actualServiceStatus := settings.gslb.Status.ServiceHealth[unhealthyHost]
	assert.Equal(t, expectedServiceStatus, actualServiceStatus, "expected %s service status to be %s, but got %s", unhealthyHost, expectedServiceStatus, actualServiceStatus)
}

func TestHealthyServiceStatus(t *testing.T) {
	// "Healthy service status" independent
	// arrange
	defer cleanup()
	settings := provideSettings(t, predefinedConfig)
	serviceName := "frontend-podinfo"
	expectedServiceStatus := "Healthy"
	healthyHost := "app3.cloud.example.com"
	defer deleteHealthyService(t, serviceName, &settings)
	createHealthyService(t, serviceName, &settings)
	reconcileAndUpdateGslb(t, settings)
	// act
	actualServiceStatus := settings.gslb.Status.ServiceHealth[healthyHost]
	// assert
	assert.Equal(t, expectedServiceStatus, actualServiceStatus, "expected %s service status to be %s, but got %s", healthyHost, expectedServiceStatus, actualServiceStatus)
}

func TestIngressHostsPerStatusMetric(t *testing.T) {
	// "ingress_hosts_per_status metric"
	// arrange
	defer cleanup()
	settings := provideSettings(t, predefinedConfig)
	expectedHostsMetricCount := 3
	// act
	err := settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.gslb)
	actualHostsMetricCount := testutil.CollectAndCount(ingressHostsPerStatusMetric)
	// assert
	assert.NoError(t, err, "Failed to get expected gslb")
	assert.Equal(t, expectedHostsMetricCount, actualHostsMetricCount, "expected %v managed hosts, but got %v", expectedHostsMetricCount, actualHostsMetricCount)
}

func TestIngressHostsPerStatusMetricReflectionForHealthyStatus(t *testing.T) {
	// "ingress_hosts_per_status metric reflection for Healthy status" was executed twice.
	// Originally it reuse service "frontend-podinfo" from "Healthy service status" in and then it recreated the new instance of "frontend-podinfo"
	// I'm running test multiple times to check that it work properly when healthy service is up and down multiple times
	defer cleanup()
	for i := 0; i < 4; i++ {
		func() {
			// arrange
			settings := provideSettings(t, predefinedConfig)
			serviceName := "frontend-podinfo"
			defer deleteHealthyService(t, serviceName, &settings)
			expectedHostsMetric := 1.
			createHealthyService(t, serviceName, &settings)
			reconcileAndUpdateGslb(t, settings)
			// act
			err := settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.gslb)
			healthyHosts := ingressHostsPerStatusMetric.With(prometheus.Labels{"namespace": settings.gslb.Namespace, "name": settings.gslb.Name, "status": healthyStatus})
			actualHostsMetric := testutil.ToFloat64(healthyHosts)
			// assert
			assert.NoError(t, err, "Failed to get expected gslb")
			assert.Equal(t, expectedHostsMetric, actualHostsMetric, "expected %v managed hosts with Healthy status, but got %v", expectedHostsMetric, actualHostsMetric)
		}()
	}
}

func TestIngressHostsPerStatusMetricReflectionForUnhealthyStatus(t *testing.T) {
	//originally "ingress_hosts_per_status metric reflection for Unhealthy status"
	// arrange
	defer cleanup()
	settings := provideSettings(t, predefinedConfig)
	err := settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.gslb)
	expectedHostsMetricCount := 0.
	// act
	unhealthyHosts := ingressHostsPerStatusMetric.With(prometheus.Labels{"namespace": settings.gslb.Namespace, "name": settings.gslb.Name, "status": unhealthyStatus})
	actualHostsMetricCount := testutil.ToFloat64(unhealthyHosts)
	// assert
	assert.NoError(t, err, "Failed to get expected gslb")
	assert.Equal(t, expectedHostsMetricCount, actualHostsMetricCount, "expected %v managed hosts, but got %v", expectedHostsMetricCount, actualHostsMetricCount)

	// arrange
	serviceName := "unhealthy-app"
	createUnhealthyService(t, serviceName, &settings)
	defer deleteUnhealthyService(t, serviceName, &settings)
	reconcileAndUpdateGslb(t, settings)
	expectedHostsMetricCount = 1
	// act
	unhealthyHosts = ingressHostsPerStatusMetric.With(prometheus.Labels{"namespace": settings.gslb.Namespace, "name": settings.gslb.Name, "status": unhealthyStatus})
	actualHostsMetricCount = testutil.ToFloat64(unhealthyHosts)
	// assert
	assert.Equal(t, expectedHostsMetricCount, actualHostsMetricCount, "expected %v managed hosts with Healthy status, but got %v", expectedHostsMetricCount, actualHostsMetricCount)
}

func TestIngressHostsPerStatusMetricReflectionForNotFoundStatus(t *testing.T) {
	// originally "ingress_hosts_per_status metric reflection for NotFound status"
	// dependent on "ingress_hosts_per_status metric reflection for Unhealthy status"
	// arrange
	defer cleanup()
	settings := provideSettings(t, predefinedConfig)
	expectedHostsMetricCount := 2.0

	// TODO: Ask Yury, what happens in this test. This block is legacy from dependent test "ingress_hosts_per_status metric reflection for Unhealthy status"
	// If I remove it, there are 3 NotFound services. Why 3 ? By adding new services I'm reducing NotFound ?How does it work?
	serviceName := "unhealthy-app"
	createUnhealthyService(t, serviceName, &settings)
	reconcileAndUpdateGslb(t, settings)
	deleteUnhealthyService(t, serviceName, &settings)

	// act
	err := settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.gslb)
	require.NoError(t, err, "Failed to get expected gslb")
	unknownHosts, err := ingressHostsPerStatusMetric.GetMetricWith(prometheus.Labels{"namespace": settings.gslb.Namespace, "name": settings.gslb.Name, "status": notFoundStatus})
	require.NoError(t, err, "Failed to get ingress metrics")
	actualHostsMetricCount := testutil.ToFloat64(unknownHosts)
	// assert
	assert.Equal(t, expectedHostsMetricCount, actualHostsMetricCount, "expected %v managed hosts with NotFound status, but got %v", expectedHostsMetricCount, actualHostsMetricCount)
}

func TestHealthyRecordMetric(t *testing.T) {
	// originally "healthy_records metric"; independent
	// arrange
	defer cleanup()
	expectedHealthyRecordsMetricCount := 3.0
	ingressIPs := []corev1.LoadBalancerIngress{
		{IP: "10.0.0.1"},
		{IP: "10.0.0.2"},
		{IP: "10.0.0.3"},
	}
	serviceName := "frontend-podinfo"
	settings := provideSettings(t, predefinedConfig)
	err := settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.gslb)
	require.NoError(t, err, "Failed to get expected gslb")
	defer deleteHealthyService(t, serviceName, &settings)
	createHealthyService(t, serviceName, &settings)
	err = settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.ingress)
	require.NoError(t, err, "Failed to get expected ingress")
	settings.ingress.Status.LoadBalancer.Ingress = append(settings.ingress.Status.LoadBalancer.Ingress, ingressIPs...)
	err = settings.client.Status().Update(context.TODO(), settings.ingress)
	require.NoError(t, err, "Failed to update gslb Ingress Address")
	reconcileAndUpdateGslb(t, settings)
	// act
	actualHealthyRecordsMetricCount := testutil.ToFloat64(healthyRecordsMetric)
	reconcileAndUpdateGslb(t, settings)
	// assert
	assert.Equal(t, expectedHealthyRecordsMetricCount, actualHealthyRecordsMetricCount, "expected %v healthy records, but got %v", expectedHealthyRecordsMetricCount, actualHealthyRecordsMetricCount)
}

func TestMetricLinterCheck(t *testing.T) {
	// originally name+" metric linter check"
	// TODO: ask Yury what is this test good for
	// arrange
	for name, scenario := range map[string]prometheus.Collector{
		"healthy_records":          healthyRecordsMetric,
		"ingress_hosts_per_status": ingressHostsPerStatusMetric,
	} {
		// act
		// assert
		lintErrors, err := testutil.CollectAndLint(scenario)
		assert.NoError(t, err)
		assert.True(t, len(lintErrors) == 0, "Metric linting error(s): %s - %s", name, lintErrors)
	}
}

func TestGslbCreatesDNSEndpointCRForHealthyIngressHosts(t *testing.T) {
	// "Gslb creates DNSEndpoint CR for healthy ingress hosts" was depending on "healthy_records metric" ingress
	// arrange
	defer cleanup()
	serviceName := "frontend-podinfo"
	dnsEndpoint := &externaldns.DNSEndpoint{}
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
	ingressIPs := []corev1.LoadBalancerIngress{
		{IP: "10.0.0.1"},
		{IP: "10.0.0.2"},
		{IP: "10.0.0.3"},
	}
	settings := provideSettings(t, predefinedConfig)
	err := settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.ingress)
	require.NoError(t, err, "Failed to get expected ingress")
	settings.ingress.Status.LoadBalancer.Ingress = append(settings.ingress.Status.LoadBalancer.Ingress, ingressIPs...)
	err = settings.client.Status().Update(context.TODO(), settings.ingress)
	require.NoError(t, err, "Failed to update gslb Ingress Address")
	createHealthyService(t, serviceName, &settings)
	defer deleteHealthyService(t, serviceName, &settings)
	reconcileAndUpdateGslb(t, settings)

	// act
	err = settings.client.Get(context.TODO(), settings.request.NamespacedName, dnsEndpoint)
	require.NoError(t, err, "Failed to load DNS endpoint")
	got := dnsEndpoint.Spec.Endpoints
	prettyGot := prettyPrint(got)
	prettyWant := prettyPrint(want)

	// assert
	assert.Equal(t, want, got, "got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
}

func TestDNSRecordReflectionInStatus(t *testing.T) {
	// "DNS Record reflection in status" was depending on "Gslb creates DNSEndpoint CR for healthy ingress hosts"
	// arrange
	defer cleanup()
	serviceName := "frontend-podinfo"
	dnsEndpoint := &externaldns.DNSEndpoint{}
	want := map[string][]string{"app3.cloud.example.com": {"10.0.0.1", "10.0.0.2", "10.0.0.3"}}
	ingressIPs := []corev1.LoadBalancerIngress{
		{IP: "10.0.0.1"},
		{IP: "10.0.0.2"},
		{IP: "10.0.0.3"},
	}
	settings := provideSettings(t, predefinedConfig)
	err := settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.ingress)
	require.NoError(t, err, "Failed to get expected ingress")
	settings.ingress.Status.LoadBalancer.Ingress = append(settings.ingress.Status.LoadBalancer.Ingress, ingressIPs...)
	err = settings.client.Status().Update(context.TODO(), settings.ingress)
	require.NoError(t, err, "Failed to update gslb Ingress Address")

	// act
	createHealthyService(t, serviceName, &settings)
	defer deleteHealthyService(t, serviceName, &settings)
	reconcileAndUpdateGslb(t, settings)
	err = settings.client.Get(context.TODO(), settings.request.NamespacedName, dnsEndpoint)
	require.NoError(t, err, "Failed to load DNS endpoint")
	got := settings.gslb.Status.HealthyRecords

	// assert
	assert.Equal(t, got, want, "got:\n %s healthyRecordsMetric status,\n\n want:\n %s", got, want)
}

func TestLocalDNSRecordsHasSpecialAnnotation(t *testing.T) {
	// "Local DNS records has special annotation" was depending on "Gslb creates DNSEndpoint CR for healthy ingress hosts"
	// arrange
	defer cleanup()
	serviceName := "frontend-podinfo"
	dnsEndpoint := &externaldns.DNSEndpoint{}
	want := "local"
	ingressIPs := []corev1.LoadBalancerIngress{
		{IP: "10.0.0.1"},
		{IP: "10.0.0.2"},
		{IP: "10.0.0.3"},
	}
	settings := provideSettings(t, predefinedConfig)
	err := settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.ingress)
	require.NoError(t, err, "Failed to get expected ingress")
	settings.ingress.Status.LoadBalancer.Ingress = append(settings.ingress.Status.LoadBalancer.Ingress, ingressIPs...)
	err = settings.client.Status().Update(context.TODO(), settings.ingress)
	require.NoError(t, err, "Failed to update gslb Ingress Address")

	// act
	createHealthyService(t, serviceName, &settings)
	defer deleteHealthyService(t, serviceName, &settings)
	reconcileAndUpdateGslb(t, settings)
	err = settings.client.Get(context.TODO(), settings.request.NamespacedName, dnsEndpoint)
	require.NoError(t, err, "Failed to load DNS endpoint")
	got := dnsEndpoint.Annotations["k8gb.absa.oss/dnstype"]

	// assert
	assert.Equal(t, got, want, "got:\n %q annotation value,\n\n want:\n %q", got, want, got, want)
}

func TestGeneratesProperExternalNSTargetFQDNsAccordingToTheGeoTags(t *testing.T) {
	//"Generates proper external NS target FQDNs according to the geo tags" independent
	// arrange
	defer cleanup()
	want := []string{"gslb-ns-za.example.com"}
	customConfig := predefinedConfig
	customConfig.EdgeDNSZone = "example.com"
	customConfig.ExtClustersGeoTags = []string{"za"}
	settings := provideSettings(t, customConfig)
	// act
	got := getExternalClusterFQDNs(settings.gslb)
	// assert
	assert.Equal(t, want, got, "got:\n %q externalGslb NS records,\n\n want:\n %q", got, want)
}

func TestCanGetExternalTargetsFromK8gbInAnotherLocation(t *testing.T) {
	// "Can get external targets from k8gb in another location" was depending on "healthy_records metric"
	// arrange
	defer cleanup()
	serviceName := "frontend-podinfo"
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
	hrWant := map[string][]string{"app3.cloud.example.com": {"10.0.0.1", "10.0.0.2", "10.0.0.3", "10.1.0.1", "10.1.0.2", "10.1.0.3"}}
	ingressIPs := []corev1.LoadBalancerIngress{
		{IP: "10.0.0.1"},
		{IP: "10.0.0.2"},
		{IP: "10.0.0.3"},
	}
	dnsEndpoint := &externaldns.DNSEndpoint{}
	settings := provideSettings(t, predefinedConfig)
	err := os.Setenv(depresolver.OverrideWithFakeDNSKey, "true")
	require.NoError(t, err, "Can't setup env var: (%v)", depresolver.OverrideWithFakeDNSKey)

	err = settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.ingress)
	require.NoError(t, err, "Failed to get expected ingress")
	settings.ingress.Status.LoadBalancer.Ingress = append(settings.ingress.Status.LoadBalancer.Ingress, ingressIPs...)
	err = settings.client.Status().Update(context.TODO(), settings.ingress)
	require.NoError(t, err, "Failed to update gslb Ingress Address")

	// act
	createHealthyService(t, serviceName, &settings)
	defer deleteHealthyService(t, serviceName, &settings)
	reconcileAndUpdateGslb(t, settings)
	err = settings.client.Get(context.TODO(), settings.request.NamespacedName, dnsEndpoint)
	require.NoError(t, err, "Failed to get expected DNSEndpoint")

	got := dnsEndpoint.Spec.Endpoints
	hrGot := settings.gslb.Status.HealthyRecords
	prettyGot := prettyPrint(got)
	prettyWant := prettyPrint(want)

	// assert
	assert.Equal(t, want, got, "got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
	assert.Equal(t, hrGot, hrWant, "got:\n %s Gslb Records status,\n\n want:\n %s", hrGot, hrWant)
}

func TestCanCheckExternalGslbTXTRecordForValidityAndFailIfItIsExpired(t *testing.T) {
	// "Can check external Gslb TXT record for validity and fail if it is expired" independent
	// arrange
	defer cleanup()
	err := os.Setenv(depresolver.OverrideWithFakeDNSKey, "true")
	require.NoError(t, err, "Can't setup env var: (%v)", depresolver.OverrideWithFakeDNSKey)
	// act
	got := checkAliveFromTXT("fake", "test-gslb-heartbeat-eu.example.com", time.Minute*5)
	want := errors.NewGone("Split brain TXT record expired the time threshold: (5m0s)")
	// assert
	assert.Equal(t, want, got, "got:\n %s from TXT split brain check,\n\n want error:\n %v", got, want)
}

func TestCanFilterOutDelegatedZoneEntryAccordingFQDNProvided(t *testing.T) {
	// "Can filter out delegated zone entry according FQDN provided" depending
	// on "Generates proper external NS target FQDNs according to the geo tags"
	// arrange
	defer cleanup()
	delegateTo := []ibclient.NameServer{
		{Address: "10.0.0.1", Name: "gslb-ns-eu.example.com"},
		{Address: "10.0.0.2", Name: "gslb-ns-eu.example.com"},
		{Address: "10.0.0.3", Name: "gslb-ns-eu.example.com"},
		{Address: "10.1.0.1", Name: "gslb-ns-za.example.com"},
		{Address: "10.1.0.2", Name: "gslb-ns-za.example.com"},
		{Address: "10.1.0.3", Name: "gslb-ns-za.example.com"},
	}
	want := []ibclient.NameServer{
		{Address: "10.0.0.1", Name: "gslb-ns-eu.example.com"},
		{Address: "10.0.0.2", Name: "gslb-ns-eu.example.com"},
		{Address: "10.0.0.3", Name: "gslb-ns-eu.example.com"},
	}
	customConfig := predefinedConfig
	customConfig.EdgeDNSZone = "example.com"
	customConfig.ExtClustersGeoTags = []string{"za"}
	settings := provideSettings(t, customConfig)
	// act
	extClusters := getExternalClusterFQDNs(settings.gslb)
	got := filterOutDelegateTo(delegateTo, extClusters[0])
	// assert
	assert.Equal(t, want, got, "got:\n %q filtered out delegation records,\n\n want:\n %q", got, want)
}

func TestCanGenerateExternalHeartbeatFQDNs(t *testing.T) {
	// "Can generate external heartbeat FQDNs" depending
	// on "Generates proper external NS target FQDNs according to the geo tags"
	// arrange
	defer cleanup()
	want := []string{"test-gslb-heartbeat-za.example.com"}
	customConfig := predefinedConfig
	customConfig.EdgeDNSZone = "example.com"
	customConfig.ExtClustersGeoTags = []string{"za"}
	settings := provideSettings(t, customConfig)
	// act
	got := getExternalClusterHeartbeatFQDNs(settings.gslb)
	// assert
	assert.Equal(t, want, got, "got:\n %s unexpected heartbeat records,\n\n want:\n %s", got, want)
}

func TestCanCheckExternalGslbTXTRecordForValidityAndPAssIfItISNotExpired(t *testing.T) {
	// "Can check external Gslb TXT record for validity and pass if it is not expired"
	// arrange
	err := os.Setenv(depresolver.OverrideWithFakeDNSKey, "true")
	require.NoError(t, err, "Can't setup env var: (%v)", depresolver.OverrideWithFakeDNSKey)
	// act
	err2 := checkAliveFromTXT("fake", "test-gslb-heartbeat-za.example.com", time.Minute*5)
	// assert
	assert.NoError(t, err2, "got:\n %s from TXT split brain check,\n\n want error:\n %v", err2, nil)
}

func TestReturnsOwnRecordsUsingFailoverStrategyWhenPrimary(t *testing.T) {
	//"Returns own records using Failover strategy when Primary" depends on "healthy_records metric"
	// arrange
	defer cleanup()
	serviceName := "frontend-podinfo"
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
	ingressIPs := []corev1.LoadBalancerIngress{
		{IP: "10.0.0.1"},
		{IP: "10.0.0.2"},
		{IP: "10.0.0.3"},
	}
	dnsEndpoint := &externaldns.DNSEndpoint{}
	customConfig := predefinedConfig
	customConfig.ClusterGeoTag = "eu"
	settings := provideSettings(t, customConfig)
	err := os.Setenv(depresolver.OverrideWithFakeDNSKey, "true")
	require.NoError(t, err, "Can't setup env var: (%v)", depresolver.OverrideWithFakeDNSKey)

	// ingress
	err = settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.ingress)
	require.NoError(t, err, "Failed to get expected ingress")
	settings.ingress.Status.LoadBalancer.Ingress = append(settings.ingress.Status.LoadBalancer.Ingress, ingressIPs...)
	err = settings.client.Status().Update(context.TODO(), settings.ingress)
	require.NoError(t, err, "Failed to update gslb Ingress Address")

	// enable failover strategy
	settings.gslb.Spec.Strategy.Type = "failover"
	settings.gslb.Spec.Strategy.PrimaryGeoTag = "eu"
	err = settings.client.Update(context.TODO(), settings.gslb)
	require.NoError(t, err, "Can't update gslb")

	// act
	createHealthyService(t, serviceName, &settings)
	defer deleteHealthyService(t, serviceName, &settings)
	reconcileAndUpdateGslb(t, settings)
	err = settings.client.Get(context.TODO(), settings.request.NamespacedName, dnsEndpoint)
	require.NoError(t, err, "Failed to get expected DNSEndpoint")
	got := dnsEndpoint.Spec.Endpoints
	prettyGot := prettyPrint(got)
	prettyWant := prettyPrint(want)

	// assert
	assert.Equal(t, want, got, "got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
}

func TestReturnsExternalRecordsUsingFailoverStrategy(t *testing.T) {
	// "Returns external records using Failover strategy when Secondary" depends on "healthy_records metric"
	// arrange
	serviceName := "frontend-podinfo"
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
	ingressIPs := []corev1.LoadBalancerIngress{
		{IP: "10.0.0.1"},
		{IP: "10.0.0.2"},
		{IP: "10.0.0.3"},
	}
	dnsEndpoint := &externaldns.DNSEndpoint{}
	customConfig := predefinedConfig
	customConfig.ClusterGeoTag = "za"
	settings := provideSettings(t, customConfig)
	err := os.Setenv(depresolver.OverrideWithFakeDNSKey, "true")
	require.NoError(t, err, "Can't setup env var: (%v)", depresolver.OverrideWithFakeDNSKey)

	// ingress
	err = settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.ingress)
	require.NoError(t, err, "Failed to get expected ingress")
	settings.ingress.Status.LoadBalancer.Ingress = append(settings.ingress.Status.LoadBalancer.Ingress, ingressIPs...)
	err = settings.client.Status().Update(context.TODO(), settings.ingress)
	require.NoError(t, err, "Failed to update gslb Ingress Address")

	// enable failover strategy
	settings.gslb.Spec.Strategy.Type = "failover"
	settings.gslb.Spec.Strategy.PrimaryGeoTag = "eu"
	err = settings.client.Update(context.TODO(), settings.gslb)
	require.NoError(t, err, "Can't update gslb")

	// act
	createHealthyService(t, serviceName, &settings)
	defer deleteHealthyService(t, serviceName, &settings)
	reconcileAndUpdateGslb(t, settings)
	err = settings.client.Get(context.TODO(), settings.request.NamespacedName, dnsEndpoint)
	require.NoError(t, err, "Failed to get expected DNSEndpoint")
	got := dnsEndpoint.Spec.Endpoints
	prettyGot := prettyPrint(got)
	prettyWant := prettyPrint(want)

	// assert
	assert.Equal(t, want, got, "got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
}

func TestGslbProperlyPropagatesAnnotationDownToIngress(t *testing.T) {
	// "Gslb properly propagates annotation down to Ingress" independent
	// arrange
	defer cleanup()
	settings := provideSettings(t, predefinedConfig)
	expectedAnnotations := map[string]string{"annotation": "test"}
	settings.gslb.Annotations = expectedAnnotations
	err := settings.client.Update(context.TODO(), settings.gslb)
	require.NoError(t, err, "Can't update gslb")
	// act
	reconcileAndUpdateGslb(t, settings)
	err2 := settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.ingress)
	// assert
	assert.NoError(t, err2, "Failed to get expected ingress")
	assert.Equal(t, expectedAnnotations, settings.ingress.Annotations)
}

func TestReflectGeoTagInStatusAsUnsetByDefault(t *testing.T) {
	// "Reflect GeoTag in the Status as unset by default" independent
	// arrange
	defer cleanup()
	want := "us-west-1"
	settings := provideSettings(t, predefinedConfig)
	// act
	reconcileAndUpdateGslb(t, settings)
	got := settings.gslb.Status.GeoTag
	// assert
	assert.Equal(t, want, got, "got: '%s' GeoTag status, want:'%s'", got, want)
}

func TestReflectGeoTagInTheStatus(t *testing.T) {
	// "Reflect GeoTag in the Status" independent
	// arrange
	defer cleanup()
	want := "eu"
	customConfig := predefinedConfig
	customConfig.ClusterGeoTag = "eu"
	settings := provideSettings(t, customConfig)
	// act
	reconcileAndUpdateGslb(t, settings)
	got := settings.gslb.Status.GeoTag
	// assert
	assert.Equal(t, want, got, "got: '%s' GeoTag status, want:'%s'", got, want)
}

func TestDetectsIngressHostnameMismatch(t *testing.T) {
	//"Detects Ingress hostname mismatch" independent
	// arrange
	defer cleanup()
	//getting Gslb and Reconciler
	settings := provideSettings(t, predefinedConfig)
	err := os.Setenv(depresolver.EdgeDNSZoneKey, "otherdnszone.com")
	require.NoError(t, err, "Can't set env var: (%v)", depresolver.EdgeDNSZoneKey)
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      settings.gslb.Name,
			Namespace: settings.gslb.Namespace,
		},
	}
	// act
	_, err = settings.reconciler.Reconcile(req)
	log.Info(fmt.Sprintf("got an error from controller: %s", err))
	// assert
	assert.Error(t, err, "expected controller to detect Ingress hostname and edgeDNSZone mismatch")
	assert.True(t, strings.HasSuffix(err.Error(), "cloud.example.com does not match delegated zone otherdnszone.com"))
}

func TestCreatesNSDNSRecordsForRoute53(t *testing.T) {
	// "Creates NS DNS records for route53" independent
	// arrange
	defer cleanup()
	const dnsZone = "cloud.example.com"
	const want = "route53"
	const coreDNSLBServiceName = "k8gb-coredns-lb"
	wantEp := []*externaldns.Endpoint{
		{
			DNSName:    dnsZone,
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
	dnsEndpointRoute53 := &externaldns.DNSEndpoint{}
	customConfig := predefinedConfig
	customConfig.EdgeDNSServer = "1.1.1.1"
	coreDNSService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      coreDNSLBServiceName,
			Namespace: k8gbNamespace,
		},
	}
	serviceIPs := []corev1.LoadBalancerIngress{
		{Hostname: "one.one.one.one"}, // rely on 1.1.1.1 response from Cloudflare
	}
	settings := provideSettings(t, customConfig)
	err := settings.client.Create(context.TODO(), coreDNSService)
	require.NoError(t, err, "Failed to create testing %s service", coreDNSLBServiceName)
	coreDNSService.Status.LoadBalancer.Ingress = append(coreDNSService.Status.LoadBalancer.Ingress, serviceIPs...)
	err = settings.client.Status().Update(context.TODO(), coreDNSService)
	require.NoError(t, err, "Failed to update coredns service lb hostname")

	// act
	customConfig.Route53Enabled = true
	customConfig.ClusterGeoTag = "eu"
	customConfig.ExtClustersGeoTags = []string{"za", "us"}
	customConfig.DNSZone = dnsZone
	// apply new environment variables and update config only
	configureEnvVar(customConfig)
	settings.reconciler.Config = &customConfig

	reconcileAndUpdateGslb(t, settings)
	err = settings.client.Get(context.TODO(), client.ObjectKey{Namespace: k8gbNamespace, Name: "k8gb-ns-route53"}, dnsEndpointRoute53)
	require.NoError(t, err, "Failed to get expected DNSEndpoint")
	got := dnsEndpointRoute53.Annotations["k8gb.absa.oss/dnstype"]
	gotEp := dnsEndpointRoute53.Spec.Endpoints
	prettyGot := prettyPrint(gotEp)
	prettyWant := prettyPrint(wantEp)

	// assert
	assert.Equal(t, want, got, "got:\n %q annotation value,\n\n want:\n %q", got, want)
	assert.Equal(t, wantEp, gotEp, "got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
}

func TestResolvesLoadBalancerHostnameFromIngressStatus(t *testing.T) {
	// "Resolves LoadBalancer hostname from Ingress status" independent
	// arrange
	defer cleanup()
	serviceName := "frontend-podinfo"
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
	settings := provideSettings(t, predefinedConfig)
	dnsEndpoint := &externaldns.DNSEndpoint{ObjectMeta: metav1.ObjectMeta{Namespace: settings.gslb.Namespace, Name: settings.gslb.Name}}
	createHealthyService(t, serviceName, &settings)
	defer deleteHealthyService(t, serviceName, &settings)
	err := settings.client.Get(context.TODO(), settings.request.NamespacedName, settings.ingress)
	require.NoError(t, err, "Failed to get expected ingress")

	settings.ingress.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{Hostname: "one.one.one.one"}}
	err = settings.client.Status().Update(context.TODO(), settings.ingress)
	require.NoError(t, err, "Failed to update gslb Ingress Address")

	// act
	err = settings.client.Delete(context.Background(), dnsEndpoint)
	require.NoError(t, err, "Failed to update DNSEndpoint")
	reconcileAndUpdateGslb(t, settings)
	err = settings.client.Get(context.TODO(), settings.request.NamespacedName, dnsEndpoint)
	require.NoError(t, err, "Failed to get expected DNSEndpoint")
	got := dnsEndpoint.Spec.Endpoints
	prettyGot := prettyPrint(got)
	prettyWant := prettyPrint(want)

	// assert
	assert.Equal(t, want, got, "got:\n %s DNSEndpoint,\n\n want:\n %s", prettyGot, prettyWant)
}

func TestMain(m *testing.M) {
	// setup tests
	fakedns()
	// run tests
	exitVal := m.Run()
	// teardown
	os.Exit(exitVal)
}

func createHealthyService(t *testing.T, serviceName string, s *testSettings) {
	t.Helper()
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: s.gslb.Namespace,
		},
	}
	err := s.client.Create(context.TODO(), service)
	if err != nil {
		t.Fatalf("Failed to create testing service: (%v)", err)
	}

	// Create fake endpoint with populated address slice
	endpoint := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: s.gslb.Namespace,
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{{IP: "1.2.3.4"}},
			},
		},
	}

	err = s.client.Create(context.TODO(), endpoint)
	if err != nil {
		t.Fatalf("Failed to create testing endpoint: (%v)", err)
	}
}

func deleteHealthyService(t *testing.T, serviceName string, s *testSettings) {
	t.Helper()
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: s.gslb.Namespace,
		},
	}
	err := s.client.Delete(context.TODO(), service)
	if err != nil {
		t.Fatalf("Failed to delete testing service: (%v)", err)
	}

	// Create fake endpoint with populated address slice
	endpoint := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: s.gslb.Namespace,
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{{IP: "1.2.3.4"}},
			},
		},
	}

	err = s.client.Delete(context.TODO(), endpoint)
	if err != nil {
		t.Fatalf("Failed to delete testing endpoint: (%v)", err)
	}
}

func createUnhealthyService(t *testing.T, serviceName string, s *testSettings) {
	t.Helper()
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: s.gslb.Namespace,
		},
	}

	err := s.client.Create(context.TODO(), service)
	if err != nil {
		t.Fatalf("Failed to create testing service: (%v)", err)
	}

	endpoint := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: s.gslb.Namespace,
		},
	}

	err = s.client.Create(context.TODO(), endpoint)
	if err != nil {
		t.Fatalf("Failed to create testing endpoint: (%v)", err)
	}

}

// TODO: refactor this to accept settings struct
func deleteUnhealthyService(t *testing.T, serviceName string, s *testSettings) {
	t.Helper()
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: s.gslb.Namespace,
		},
	}

	err := s.client.Delete(context.TODO(), service)
	if err != nil {
		t.Fatalf("Failed to delete testing service: (%v)", err)
	}

	endpoint := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: s.gslb.Namespace,
		},
	}

	err = s.client.Delete(context.TODO(), endpoint)
	if err != nil {
		t.Fatalf("Failed to delete testing endpoint: (%v)", err)
	}

}

// TODO: refactor this to accept settings struct
func reconcileAndUpdateGslb(t *testing.T, s testSettings) {
	t.Helper()
	// Reconcile again so Reconcile() checks services and updates the Gslb
	// resources' Status.
	res, err := s.reconciler.Reconcile(s.request)
	if err != nil {
		return
	}
	if res != (reconcile.Result{RequeueAfter: time.Second * 30}) {
		t.Error("reconcile did not return Result with Requeue")
	}

	err = s.client.Get(context.TODO(), s.request.NamespacedName, s.gslb)
	if err != nil {
		t.Fatalf("Failed to get expected gslb: (%v)", err)
	}
}

func provideSettings(t *testing.T, expected depresolver.Config) (settings testSettings) {
	configureEnvVar(expected)
	_, err := os.Stat(crSampleYaml)
	if os.IsNotExist(err) {
		t.Fatalf("Sample CR yaml file not found at: %s", crSampleYaml)
	}
	gslbYaml, err := ioutil.ReadFile(crSampleYaml)
	if err != nil {
		t.Fatalf("Can't open example CR file: %s", crSampleYaml)
	}
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(zap.New(zap.UseDevMode(true)))
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
	cl := fake.NewFakeClientWithScheme(s, objs...)
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
	settings = testSettings{
		gslb:       gslb,
		config:     expected,
		reconciler: r,
		request:    req,
		client:     cl,
		ingress:    ingress,
	}
	reconcileAndUpdateGslb(t, settings)
	return
}

func cleanup() {
	for _, s := range []string{depresolver.ReconcileRequeueSecondsKey, depresolver.ClusterGeoTagKey, depresolver.ExtClustersGeoTagsKey,
		depresolver.EdgeDNSZoneKey, depresolver.DNSZoneKey, depresolver.EdgeDNSServerKey,
		depresolver.Route53EnabledKey, depresolver.InfobloxGridHostKey, depresolver.InfobloxVersionKey, depresolver.InfobloxPortKey,
		depresolver.InfobloxUsernameKey, depresolver.InfobloxPasswordKey, depresolver.OverrideWithFakeDNSKey, depresolver.FakeInfoblox} {
		if os.Unsetenv(s) != nil {
			panic(fmt.Errorf("cleanup %s", s))
		}
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
	_ = os.Setenv(depresolver.OverrideWithFakeDNSKey, "false")
	_ = os.Setenv(depresolver.FakeInfoblox, "true")
}
