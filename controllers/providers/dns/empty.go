package dns

import (
	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	"github.com/AbsaOSS/k8gb/controllers/depresolver"
	"github.com/AbsaOSS/k8gb/controllers/providers/assistant"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

// EmptyDNSProvider is executed when fakeDNSEnabled is true.
type EmptyDNSProvider struct {
	assistant assistant.IAssistant
	config    depresolver.Config
}

func NewEmptyDNS(config depresolver.Config, assistant assistant.IAssistant) *EmptyDNSProvider {
	return &EmptyDNSProvider{
		config:    config,
		assistant: assistant,
	}
}

func (p *EmptyDNSProvider) CreateZoneDelegationForExternalDNS(*k8gbv1beta1.Gslb) (err error) {
	return
}

func (p *EmptyDNSProvider) GslbIngressExposedIPs(gslb *k8gbv1beta1.Gslb) (r []string, err error) {
	return p.assistant.GslbIngressExposedIPs(gslb)
}

func (p *EmptyDNSProvider) GetExternalTargets(host string) (targets []string) {
	return p.assistant.GetExternalTargets(host, p.config.Override.FakeDNSEnabled, nsServerNameExt(p.config))
}

func (p *EmptyDNSProvider) SaveDNSEndpoint(gslb *k8gbv1beta1.Gslb, i *externaldns.DNSEndpoint) error {
	return p.assistant.SaveDNSEndpoint(gslb.Namespace, i)
}

func (p *EmptyDNSProvider) Finalize(gslb *k8gbv1beta1.Gslb) (err error) {
	return p.assistant.RemoveEndpoint(gslb.Name)
}

func (p *EmptyDNSProvider) String() string {
	return "EMPTY"
}
