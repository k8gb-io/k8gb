package dns

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
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/k8gb-io/k8gb/controllers/logging"

	assistant2 "github.com/k8gb-io/k8gb/controllers/providers/assistant"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	"github.com/k8gb-io/k8gb/controllers/depresolver"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

const externalDNSTypeCommon = "extdns"

type ExternalDNSProvider struct {
	assistant    assistant2.Assistant
	config       depresolver.Config
	endpointName string
}

var log = logging.Logger()

func NewExternalDNS(config depresolver.Config, assistant assistant2.Assistant) *ExternalDNSProvider {
	return &ExternalDNSProvider{
		assistant:    assistant,
		config:       config,
		endpointName: fmt.Sprintf("k8gb-ns-%s", externalDNSTypeCommon),
	}
}

func (p *ExternalDNSProvider) CreateZoneDelegationForExternalDNS(gslb *k8gbv1beta1.Gslb) error {
	ttl := externaldns.TTL(p.config.NSRecordTTL)
	log.Info().
		Interface("provider", p).
		Msg("Creating/Updating DNSEndpoint CRDs")
	NSServerList := []string{p.config.GetClusterNSName()}
	for _, v := range p.config.GetExternalClusterNSNames() {
		NSServerList = append(NSServerList, v)
	}
	sort.Strings(NSServerList)
	var NSServerIPs []string
	var err error
	if p.config.CoreDNSExposed {
		if p.config.CoreDNSClusterIPs {
			NSServerIPs, err = p.assistant.CoreDNSClusterIPs()
		} else {
			NSServerIPs, err = p.assistant.CoreDNSExposedIPs()
		}
	} else {
		if len(gslb.Status.LoadBalancer.ExposedIPs) == 0 {
			// do not update DNS Endpoint for External DNS if no IPs are exposed
			// new GSLB resources may have this field empty
			return nil
		}
		NSServerIPs = gslb.Status.LoadBalancer.ExposedIPs
	}
	if err != nil {
		return err
	}
	NSRecord := &externaldns.DNSEndpoint{
		ObjectMeta: metav1.ObjectMeta{
			Name:        p.endpointName,
			Namespace:   p.config.K8gbNamespace,
			Annotations: map[string]string{"k8gb.absa.oss/dnstype": externalDNSTypeCommon},
		},
		Spec: externaldns.DNSEndpointSpec{
			Endpoints: []*externaldns.Endpoint{
				{
					DNSName:    p.config.DNSZone,
					RecordTTL:  ttl,
					RecordType: "NS",
					Targets:    NSServerList,
				},
				{
					DNSName:    p.config.GetClusterNSName(),
					RecordTTL:  ttl,
					RecordType: "A",
					Targets:    NSServerIPs,
				},
			},
		},
	}
	err = p.assistant.SaveDNSEndpoint(p.config.K8gbNamespace, NSRecord)
	if err != nil {
		return err
	}
	return nil
}

func (p *ExternalDNSProvider) Finalize(_ *k8gbv1beta1.Gslb, k8sClient client.Client) error {
	gslbList := &k8gbv1beta1.GslbList{}

	err := k8sClient.List(context.TODO(), gslbList)
	if err != nil {
		return err
	}

	// only remove the DNSEndpoint if there are no more GSLB resources
	if len(gslbList.Items) > 1 {
		return nil
	}

	return p.assistant.RemoveEndpoint(p.endpointName)
}

func (p *ExternalDNSProvider) GetExternalTargets(host string) (targets assistant2.Targets) {
	return p.assistant.GetExternalTargets(host, p.config.GetExternalClusterNSNames())
}

func (p *ExternalDNSProvider) SaveDNSEndpoint(gslb *k8gbv1beta1.Gslb, i *externaldns.DNSEndpoint) error {
	return p.assistant.SaveDNSEndpoint(gslb.Namespace, i)
}

func (p *ExternalDNSProvider) String() string {
	return strings.ToUpper(externalDNSTypeCommon)
}
