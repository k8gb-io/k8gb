package assistant

import (
	"time"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

type IAssistant interface {
	// CoreDNSExposedIPs retrieves list of exposed IP by CoreDNS
	CoreDNSExposedIPs() ([]string, error)
	// GslbIngressExposedIPs retrieves list of IP's exposed by all GSLB ingresses
	GslbIngressExposedIPs(gslb *k8gbv1beta1.Gslb) ([]string, error)
	// GetExternalTargets retrieves slice of targets from external clusters
	GetExternalTargets(host string, fakeDNSEnabled bool, extGslbClusters []string) (targets []string)
	// SaveDNSEndpoint update DNS endpoint or create new one if doesnt exist
	SaveDNSEndpoint(namespace string, i *externaldns.DNSEndpoint) error
	// RemoveEndpoint removes endpoint
	RemoveEndpoint(endpointName string) error
	// Info wraps private logger and provides log.Error()
	// TODO: extract logging functions outside
	Info(msg string, args ...interface{})
	// Error wraps private logger and provides log.Info()
	// TODO: extract logging functions outside
	Error(err error, msg string, args ...interface{})
	// InspectTXTThreshold inspects fqdn TXT record from edgeDNSServer. If record doesn't exists or timestamp is greater than
	// splitBrainThreshold the error is returned. In case fakeDNSEnabled is true, 127.0.0.1:7753 is used as edgeDNSServer
	InspectTXTThreshold(fqdn string, fakeDNSEnabled bool, splitBrainThreshold time.Duration) error
}
