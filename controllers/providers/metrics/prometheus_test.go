package metrics

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
	"fmt"
	"os"
	"reflect"
	"runtime"
	"testing"

	externaldns "sigs.k8s.io/external-dns/endpoint"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/k8gb-io/k8gb/controllers/depresolver"
	"github.com/stretchr/testify/assert"
)

const (
	namespace        = "ns"
	gslbName         = "test-gslb"
	endpointName     = "test-gslb"
	localtargetsName = "localtargets.cloud.example.com"
	targetsName      = "cloud.example.com"
)

var (
	defaultGslb     = new(k8gbv1beta1.Gslb)
	defaultEndpoint = new(externaldns.DNSEndpoint)
	defaultConfig   = depresolver.Config{K8gbNamespace: namespace, DNSZone: "cloud.example.com"}
)

func TestMetricsSingletonIsNotNil(t *testing.T) {
	// arrange
	// act
	m := Metrics()
	// assert
	assert.NotNil(t, m)
	assert.Equal(t, DefaultMetricsNamespace, m.config.K8gbNamespace)
}

func TestMetricsSingletonInitTwice(t *testing.T) {
	// arrange
	c1 := &depresolver.Config{K8gbNamespace: "c1"}
	c2 := &depresolver.Config{K8gbNamespace: "c2"}
	// act
	Init(c1)
	Init(c2)
	m := Metrics()
	// assert
	assert.Equal(t, c1.K8gbNamespace, m.config.K8gbNamespace)
}

func TestPrometheusRegistry(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	fieldCnt := reflect.TypeOf(metrics.metrics).NumField()
	items := []string{K8gbGslbErrorsTotal, K8gbGslbHealthyRecords, K8gbGslbReconciliationLoopsTotal,
		K8gbGslbServiceStatusNum, K8gbGslbStatusCountForFailover, K8gbGslbStatusCountForRoundrobin,
		K8gbGslbStatusCountForGeoIP, K8gbInfobloxHeartbeatsTotal, K8gbInfobloxHeartbeatErrorsTotal,
		K8gbInfobloxRequestDuration, K8gbInfobloxZoneUpdatesTotal, K8gbInfobloxZoneUpdateErrorsTotal,
		K8gbEndpointStatusNum, K8gbRuntimeInfo}
	// act
	registry := m.registry()
	// assert
	for _, i := range items {
		assert.NotNil(t, registry[i], i, " is not in registry. Check init() function")
	}
	assert.Equal(t, len(registry), fieldCnt, "not all metrics are initialised, check init() function")
	assert.Equal(t, len(registry), len(items), "check that local items slice fits with collectors struct")
}

func TestMetricsRegister(t *testing.T) {
	// arrange
	m := Metrics()
	// act
	err := m.Register()
	m.Unregister()
	// assert
	assert.NoError(t, err)
}

func TestReconciliationTotal(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	cnt1 := testutil.ToFloat64(m.Get(K8gbGslbReconciliationLoopsTotal).AsCounterVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	// act
	m.IncrementReconciliation(defaultGslb)
	// assert
	cnt2 := testutil.ToFloat64(m.Get(K8gbGslbReconciliationLoopsTotal).AsCounterVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	assert.Equal(t, cnt1+1.0, cnt2)
}

func TestHealthyRecords(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	data := map[string][]string{
		"roundrobin.cloud.example.com":      {"10.0.0.1", "10.0.0.2", "10.0.0.3"},
		"roundrobin-test.cloud.example.com": {"10.0.0.4", "10.0.0.5", "10.0.0.6"},
	}
	// act
	cnt1 := testutil.ToFloat64(m.Get(K8gbGslbHealthyRecords).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	m.UpdateHealthyRecordsMetric(defaultGslb, data)
	cnt2 := testutil.ToFloat64(m.Get(K8gbGslbHealthyRecords).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	// assert
	assert.Equal(t, 0.0, cnt1)
	assert.Equal(t, 6.0, cnt2)
}

func TestEmptyHealthyRecords(t *testing.T) {
	// arrange
	var data map[string][]string
	m := newPrometheusMetrics(defaultConfig)
	// act
	cnt1 := testutil.ToFloat64(m.Get(K8gbGslbHealthyRecords).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	m.UpdateHealthyRecordsMetric(defaultGslb, data)
	cnt2 := testutil.ToFloat64(m.Get(K8gbGslbHealthyRecords).AsGaugeVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	// assert
	assert.Equal(t, 0.0, cnt1)
	assert.Equal(t, 0.0, cnt2)
}

func TestErrorIncrement(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	name := K8gbGslbErrorsTotal
	cnt1 := testutil.ToFloat64(m.Get(name).AsCounterVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	// act
	m.IncrementError(defaultGslb)
	// assert
	cnt2 := testutil.ToFloat64(m.Get(name).AsCounterVec().With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	assert.Equal(t, cnt1+1.0, cnt2)
}

func TestInfobloxZoneUpdateErrorIncrement(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	cnt1 := testutil.ToFloat64(m.Get(K8gbInfobloxZoneUpdateErrorsTotal).AsCounterVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	// act
	m.InfobloxIncrementZoneUpdateError(defaultGslb)
	// assert
	cnt2 := testutil.ToFloat64(m.Get(K8gbInfobloxZoneUpdateErrorsTotal).AsCounterVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	assert.Equal(t, cnt1+1.0, cnt2)
}

func TestInfobloxHeartbeatIncrement(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	cnt1 := testutil.ToFloat64(m.Get(K8gbInfobloxHeartbeatsTotal).AsCounterVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	// act
	m.InfobloxIncrementHeartbeat(defaultGslb)
	// assert
	cnt2 := testutil.ToFloat64(m.Get(K8gbInfobloxHeartbeatsTotal).AsCounterVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	assert.Equal(t, cnt1+1.0, cnt2)
}

func TestInfobloxHeartbeatErrorIncrement(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	cnt1 := testutil.ToFloat64(m.Get(K8gbInfobloxHeartbeatErrorsTotal).AsCounterVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	// act
	m.InfobloxIncrementHeartbeatError(defaultGslb)
	// assert
	cnt2 := testutil.ToFloat64(m.Get(K8gbInfobloxHeartbeatErrorsTotal).AsCounterVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	assert.Equal(t, cnt1+1.0, cnt2)
}

func TestInfobloxZoneUpdateIncrement(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	cnt1 := testutil.ToFloat64(m.Get(K8gbInfobloxZoneUpdatesTotal).AsCounterVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	// act
	m.InfobloxIncrementZoneUpdate(defaultGslb)
	// assert
	cnt2 := testutil.ToFloat64(m.Get(K8gbInfobloxZoneUpdatesTotal).AsCounterVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName}))
	assert.Equal(t, cnt1+1.0, cnt2)
}

func TestUpgradeIngressHost(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	var serviceHealth = map[string]k8gbv1beta1.HealthStatus{
		"roundrobin.cloud.example.com": k8gbv1beta1.Healthy,
		"failover.cloud.example.com":   k8gbv1beta1.Healthy,
		"unhealthy.cloud.example.com":  k8gbv1beta1.Unhealthy,
		"notfound.cloud.example.com":   k8gbv1beta1.NotFound,
	}
	// act
	cntHealthy1 := testutil.ToFloat64(m.Get(K8gbGslbServiceStatusNum).AsGaugeVec().With(
		prometheus.Labels{"namespace": namespace, "name": gslbName, "status": k8gbv1beta1.Healthy.String()}))
	cntUnhealthy1 := testutil.ToFloat64(m.Get(K8gbGslbServiceStatusNum).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": k8gbv1beta1.Unhealthy.String()}))
	cntNotFound1 := testutil.ToFloat64(m.Get(K8gbGslbServiceStatusNum).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": k8gbv1beta1.NotFound.String()}))
	m.UpdateIngressHostsPerStatusMetric(defaultGslb, serviceHealth)
	cntHealthy2 := testutil.ToFloat64(m.Get(K8gbGslbServiceStatusNum).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": k8gbv1beta1.Healthy.String()}))
	ctnUnhealthy2 := testutil.ToFloat64(m.Get(K8gbGslbServiceStatusNum).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": k8gbv1beta1.Unhealthy.String()}))
	cntNotFound2 := testutil.ToFloat64(m.Get(K8gbGslbServiceStatusNum).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": k8gbv1beta1.NotFound.String()}))
	// assert
	assert.Equal(t, .0, cntHealthy1)
	assert.Equal(t, .0, cntUnhealthy1)
	assert.Equal(t, .0, cntNotFound1)
	assert.Equal(t, 2., cntHealthy2)
	assert.Equal(t, 1., ctnUnhealthy2)
	assert.Equal(t, 1., cntNotFound2)
}

func TestUpdateFailover(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)

	// act
	m.UpdateFailoverStatus(defaultGslb, true, k8gbv1beta1.Healthy, []string{"10.0.0.1", "10.0.0.2"})
	// assert
	hp := testutil.ToFloat64(m.Get(K8gbGslbStatusCountForFailover).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": fmt.Sprintf("%s_%s", k8gbv1beta1.Healthy, primary)}))
	up := testutil.ToFloat64(m.Get(K8gbGslbStatusCountForFailover).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": fmt.Sprintf("%s_%s", k8gbv1beta1.Unhealthy, primary)}))
	fp := testutil.ToFloat64(m.Get(K8gbGslbStatusCountForFailover).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": fmt.Sprintf("%s_%s", k8gbv1beta1.NotFound, primary)}))
	assert.Equal(t, 2., hp)
	assert.Equal(t, 0., up)
	assert.Equal(t, 0., fp)

	// act
	m.UpdateFailoverStatus(defaultGslb, false, k8gbv1beta1.Unhealthy, []string{"10.0.1.1", "10.0.1.2"})
	// assert
	hs := testutil.ToFloat64(m.Get(K8gbGslbStatusCountForFailover).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": fmt.Sprintf("%s_%s", k8gbv1beta1.Healthy, secondary)}))
	us := testutil.ToFloat64(m.Get(K8gbGslbStatusCountForFailover).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": fmt.Sprintf("%s_%s", k8gbv1beta1.Unhealthy, secondary)}))
	fs := testutil.ToFloat64(m.Get(K8gbGslbStatusCountForFailover).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": fmt.Sprintf("%s_%s", k8gbv1beta1.NotFound, secondary)}))
	assert.Equal(t, 0., hs)
	assert.Equal(t, 2., us)
	assert.Equal(t, 0., fs)
}

func TestUpdateRoundRobin(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)

	// act
	m.UpdateRoundrobinStatus(defaultGslb, k8gbv1beta1.Healthy, []string{"10.0.0.1", "10.0.0.2"})
	// assert
	hp := testutil.ToFloat64(m.Get(K8gbGslbStatusCountForRoundrobin).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": k8gbv1beta1.Healthy.String()}))
	up := testutil.ToFloat64(m.Get(K8gbGslbStatusCountForRoundrobin).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": k8gbv1beta1.Unhealthy.String()}))
	fp := testutil.ToFloat64(m.Get(K8gbGslbStatusCountForRoundrobin).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": k8gbv1beta1.NotFound.String()}))
	assert.Equal(t, 2., hp)
	assert.Equal(t, 0., up)
	assert.Equal(t, 0., fp)

	// act
	m.UpdateRoundrobinStatus(defaultGslb, k8gbv1beta1.Unhealthy, []string{"10.0.1.1", "10.0.1.2"})
	// assert
	hs := testutil.ToFloat64(m.Get(K8gbGslbStatusCountForRoundrobin).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": k8gbv1beta1.Healthy.String()}))
	us := testutil.ToFloat64(m.Get(K8gbGslbStatusCountForRoundrobin).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": k8gbv1beta1.Unhealthy.String()}))
	fs := testutil.ToFloat64(m.Get(K8gbGslbStatusCountForRoundrobin).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": gslbName, "status": k8gbv1beta1.NotFound.String()}))
	assert.Equal(t, 0., hs)
	assert.Equal(t, 2., us)
	assert.Equal(t, 0., fs)
}

func TestEndpointStatus(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	ep := defaultEndpoint
	ep.Spec.Endpoints = []*externaldns.Endpoint{
		{
			DNSName: localtargetsName,
			Targets: []string{"2.2.2.4", "2.2.2.5"},
		},
		{
			DNSName: targetsName,
			Targets: []string{"2.2.2.4", "2.2.2.5", "3.3.3.4", "3.3.3.5"},
		}}
	// act
	m.UpdateEndpointStatus(ep)
	// assert
	cnt1 := testutil.ToFloat64(m.Get(K8gbEndpointStatusNum).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": endpointName, "dns_name": targetsName}))
	cnt2 := testutil.ToFloat64(m.Get(K8gbEndpointStatusNum).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": endpointName, "dns_name": localtargetsName}))
	cnt3 := testutil.ToFloat64(m.Get(K8gbEndpointStatusNum).AsGaugeVec().
		With(prometheus.Labels{"namespace": namespace, "name": endpointName, "dns_name": "nonexists"}))
	assert.Equal(t, 4., cnt1)
	assert.Equal(t, 2., cnt2)
	assert.Equal(t, 0., cnt3)
}

func TestRunMetricLinter(t *testing.T) {
	// arrange
	m := newPrometheusMetrics(defaultConfig)
	registry := m.registry()
	// act
	// assert
	for name, scenario := range registry {
		lintErrors, err := testutil.CollectAndLint(scenario)
		assert.NoError(t, err)
		assert.True(t, len(lintErrors) == 0, "Metric linting error(s): %s - %s", name, lintErrors)
	}
}

func TestRuntimeStatus(t *testing.T) {
	// arrange
	const version = "v0.8.1"
	const gitSHAShort = "74bf71b"
	const gitSHALarge = "74bf71b879d5326ebbf3e1172c1cb8c03b2e03a6"

	l := prometheus.Labels{
		"namespace":    namespace,
		"go_version":   runtime.Version(),
		"arch":         runtime.GOARCH,
		"os":           runtime.GOOS,
		"k8gb_version": version,
	}

	m := newPrometheusMetrics(defaultConfig)

	f := func(sha, expected string) {
		m.SetRuntimeInfo(version, sha)
		l["git_sha"] = expected
		cnt := testutil.ToFloat64(m.Get(K8gbRuntimeInfo).AsGaugeVec().With(l))
		assert.Equal(t, 1., cnt)
	}

	// act
	// assert
	f(gitSHALarge, gitSHAShort)
	for _, sha := range []string{gitSHAShort, "74bf", "none", ""} {
		f(sha, sha)
	}
}

func TestMain(m *testing.M) {
	defaultGslb.Name = gslbName
	defaultGslb.Namespace = namespace
	defaultEndpoint.Name = endpointName
	defaultEndpoint.Namespace = namespace
	os.Exit(m.Run())
}
