package gslb

import (
	kgbv1beta1 "github.com/AbsaOSS/kgb/pkg/apis/kgb/v1beta1"
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	kgbNamespace = "kgb"
	gslbSubsystem    = "gslb"
)
const (
	healthyStatus   = "Healthy"
	unhealthyStatus = "Unhealthy"
	notFoundStatus  = "NotFound"
)

// Custom gslb prometheus metrics
var (
	healthyRecordsMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: kgbNamespace,
			Subsystem: gslbSubsystem,
			Name:      "healthy_records",
			Help:      "Number of healthy records observed by KGB.",
		},
		[]string{"namespace", "name"},
	)
	ingressHostsPerStatusMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: kgbNamespace,
			Subsystem: gslbSubsystem,
			Name:      "ingress_hosts_per_status",
			Help:      "Number of managed hosts observed by KGB.",
		},
		[]string{"namespace", "name", "status"},
	)
)

func (r *ReconcileGslb) updateIngressHostsPerStatusMetric(gslb *kgbv1beta1.Gslb, serviceHealth map[string]string) error {
	var healthyHostsCount, unhealthyHostsCount, notFoundHostsCount int
	for _, hs := range serviceHealth {
		switch hs {
		case healthyStatus:
			healthyHostsCount++
		case unhealthyStatus:
			unhealthyHostsCount++
		default:
			notFoundHostsCount++
		}
	}
	ingressHostsPerStatusMetric.With(prometheus.Labels{"namespace": gslb.Namespace, "name": gslb.Name, "status": healthyStatus}).Set(float64(healthyHostsCount))
	ingressHostsPerStatusMetric.With(prometheus.Labels{"namespace": gslb.Namespace, "name": gslb.Name, "status": unhealthyStatus}).Set(float64(unhealthyHostsCount))
	ingressHostsPerStatusMetric.With(prometheus.Labels{"namespace": gslb.Namespace, "name": gslb.Name, "status": notFoundStatus}).Set(float64(notFoundHostsCount))

	return nil
}

func (r *ReconcileGslb) updateHealthyRecordsMetric(gslb *kgbv1beta1.Gslb, healthyRecords map[string][]string) error {
	var hrsCount int
	for _, hrs := range healthyRecords {
		hrsCount += len(hrs)
	}
	healthyRecordsMetric.With(prometheus.Labels{"namespace": gslb.Namespace, "name": gslb.Name}).Set(float64(hrsCount))
	return nil
}

func init() {
	metrics.Registry.MustRegister(healthyRecordsMetric)
	metrics.Registry.MustRegister(ingressHostsPerStatusMetric)
}
