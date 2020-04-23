package gslb

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	ohmyglbNamespace = "ohmyglb"
	gslbSubsystem    = "gslb"
)
const (
	healthyStatus   = "Healthy"
	unhealthyStatus = "Unhealthy"
	notFoundStatus  = "NotFound"
)

// Custom gslb prometheus metrics
var (
	healthyRecordsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: ohmyglbNamespace,
			Subsystem: gslbSubsystem,
			Name:      "healthy_records_total",
			Help:      "Number of healthy records observed by OhMyGLB.",
		},
		[]string{"namespace", "name"},
	)
	managedHostsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: ohmyglbNamespace,
			Subsystem: gslbSubsystem,
			Name:      "managed_hosts_total",
			Help:      "Number of managed hosts observed by OhMyGLB.",
		},
		[]string{"namespace", "name", "status"},
	)
)

func init() {
	metrics.Registry.MustRegister(healthyRecordsTotal)
	metrics.Registry.MustRegister(managedHostsTotal)
}
