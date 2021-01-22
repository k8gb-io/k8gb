package dns

import (
	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

type IDnsProvider interface {
	// CreateZoneDelegationForExternalDNS handles delegated zone in Edge DNS
	CreateZoneDelegationForExternalDNS(*k8gbv1beta1.Gslb) (*reconcile.Result, error)
	// GslbIngressExposedIPs retrieves list of IP's exposed by all GSLB ingresses
	GslbIngressExposedIPs(*k8gbv1beta1.Gslb) ([]string, error)
	// GetExternalTargets retrieves list of external targets for specified host
	GetExternalTargets(string) []string
	// SaveDNSEndpoint update DNS endpoint in gslb or create new one if doesn't exist
	SaveDNSEndpoint(*k8gbv1beta1.Gslb, *externaldns.DNSEndpoint) (*reconcile.Result, error)
	// Finalize finalize gslb in k8gbNamespace
	Finalize(*k8gbv1beta1.Gslb) error
}
