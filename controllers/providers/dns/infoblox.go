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
	"fmt"
	"reflect"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"
	externaldns "sigs.k8s.io/external-dns/endpoint"

	ibcl "github.com/infobloxopen/infoblox-go-client"
	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	"github.com/k8gb-io/k8gb/controllers/depresolver"
	"github.com/k8gb-io/k8gb/controllers/providers/assistant"
	"github.com/k8gb-io/k8gb/controllers/providers/metrics"
)

type InfobloxProvider struct {
	assistant assistant.Assistant
	config    depresolver.Config
	client    InfobloxClient
}

var m = metrics.Metrics()

func NewInfobloxDNS(config depresolver.Config, assistant assistant.Assistant, client InfobloxClient) *InfobloxProvider {
	return &InfobloxProvider{
		client:    client,
		assistant: assistant,
		config:    config,
	}
}

func (p *InfobloxProvider) sanitizeDelegateZone(local, upstream []ibcl.NameServer) []ibcl.NameServer {
	// Drop own records for straight away update
	// And ensure local entries are up to date
	// And final list is sorted
	final := local
	remote := p.filterOutDelegateTo(upstream, p.config.GetClusterNSName())
	final = append(final, remote...)
	sortZones(final)

	return final
}

func (p *InfobloxProvider) CreateZoneDelegationForExternalDNS() error {
	objMgr, err := p.client.GetObjectManager()
	if err != nil {
		return err
	}

	var addresses []string
	if p.config.CoreDNSExposed {
		addresses, err = p.assistant.GetCoreDNSLoadBalancerServiceIPs()
	} else {
		addresses, err = p.assistant.GetIngressStatusIPs()
	}
	if err != nil {
		return err
	}
	var delegateTo []ibcl.NameServer

	for _, address := range addresses {
		nameServer := ibcl.NameServer{Address: address, Name: p.config.GetClusterNSName()}
		delegateTo = append(delegateTo, nameServer)
	}

	findZone, err := p.getZoneDelegated(objMgr, p.config.DNSZone)
	if err != nil {
		return err
	}

	if findZone != nil {
		err = p.checkZoneDelegated(findZone)
		if err != nil {
			return err
		}

		if len(findZone.Ref) > 0 {

			sortZones(findZone.DelegateTo)
			currentList := p.sanitizeDelegateZone(delegateTo, findZone.DelegateTo)

			if !reflect.DeepEqual(findZone.DelegateTo, currentList) {
				log.Info().
					Interface("records", findZone.DelegateTo).
					Msg("Found delegated zone records")
				log.Info().
					Str("DNSZone", p.config.DNSZone).
					Interface("serverList", currentList).
					Msg("Updating delegated zone with the server list")
				_, err = p.updateZoneDelegated(objMgr, findZone.Ref, currentList)
				if err != nil {
					return err
				}
			}
		}
	} else {
		log.Info().
			Str("DNSZone", p.config.DNSZone).
			Msg("Creating delegated zone")
		sortZones(delegateTo)
		log.Debug().
			Interface("records", delegateTo).
			Msg("Delegated records")
		_, err = p.createZoneDelegated(objMgr, p.config.DNSZone, delegateTo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *InfobloxProvider) Finalize(gslb *k8gbv1beta1.Gslb, _ client.Client) error {
	objMgr, err := p.client.GetObjectManager()
	if err != nil {
		return err
	}
	findZone, err := p.getZoneDelegated(objMgr, p.config.DNSZone)
	if err != nil {
		return err
	}

	if findZone != nil {
		err = p.checkZoneDelegated(findZone)
		if err != nil {
			return err
		}
		if len(findZone.Ref) > 0 {
			log.Info().
				Str("DNSZone", p.config.DNSZone).
				Msg("Deleting delegated zone")
			_, err := p.deleteZoneDelegated(objMgr, findZone.Ref)
			if err != nil {
				return err
			}
		}
	}

	heartbeatTXTName := p.config.GetClusterHeartbeatFQDN(gslb.Name)
	findTXT, err := p.getTXTRecord(objMgr, heartbeatTXTName)
	if err != nil {
		return err
	}

	if findTXT != nil {
		if len(findTXT.Ref) > 0 {
			log.Info().
				Str("TXTRecords", heartbeatTXTName).
				Msg("Deleting split brain TXT record")
			_, err := p.deleteTXTRecord(objMgr, findTXT.Ref)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *InfobloxProvider) GetExternalTargets(host string) (targets assistant.Targets) {
	return p.assistant.GetExternalTargets(host, p.config.GetExternalClusterNSNames())
}

func (p *InfobloxProvider) SaveDNSEndpoint(gslb *k8gbv1beta1.Gslb, i *externaldns.DNSEndpoint) error {
	return p.assistant.SaveDNSEndpoint(gslb.Namespace, i)
}

func (p *InfobloxProvider) String() string {
	return "Infoblox"
}

func (p *InfobloxProvider) checkZoneDelegated(findZone *ibcl.ZoneDelegated) error {
	if findZone.Fqdn != p.config.DNSZone {
		err := fmt.Errorf("delegated zone returned from infoblox(%s) does not match requested gslb zone(%s)", findZone.Fqdn, p.config.DNSZone)
		return err
	}
	return nil
}

func (p *InfobloxProvider) filterOutDelegateTo(delegateTo []ibcl.NameServer, fqdn string) (result []ibcl.NameServer) {
	result = make([]ibcl.NameServer, 0)

	for _, v := range delegateTo {
		if v.Name != fqdn {
			result = append(result, v)
		}
	}
	return
}

func (p *InfobloxProvider) createZoneDelegated(o *ibcl.ObjectManager, fqdn string, d []ibcl.NameServer) (res *ibcl.ZoneDelegated, err error) {
	start := time.Now()
	res, err = o.CreateZoneDelegated(fqdn, d)
	m.InfobloxObserveRequestDuration(start, metrics.CreateZoneDelegated, err == nil)
	return
}

func (p *InfobloxProvider) getZoneDelegated(o *ibcl.ObjectManager, fqdn string) (res *ibcl.ZoneDelegated, err error) {
	start := time.Now()
	res, err = o.GetZoneDelegated(fqdn)
	m.InfobloxObserveRequestDuration(start, metrics.GetZoneDelegated, err == nil)
	return
}

func (p *InfobloxProvider) updateZoneDelegated(o *ibcl.ObjectManager, fqdn string, d []ibcl.NameServer) (res *ibcl.ZoneDelegated, err error) {
	start := time.Now()
	res, err = o.UpdateZoneDelegated(fqdn, d)
	m.InfobloxObserveRequestDuration(start, metrics.UpdateZoneDelegated, err == nil)
	return
}

func (p *InfobloxProvider) deleteZoneDelegated(o *ibcl.ObjectManager, fqdn string) (res string, err error) {
	start := time.Now()
	res, err = o.DeleteZoneDelegated(fqdn)
	m.InfobloxObserveRequestDuration(start, metrics.DeleteZoneDelegated, err == nil)
	return
}

func (p *InfobloxProvider) getTXTRecord(o *ibcl.ObjectManager, name string) (res *ibcl.RecordTXT, err error) {
	start := time.Now()
	res, err = o.GetTXTRecord(name)
	m.InfobloxObserveRequestDuration(start, metrics.GetTXTRecord, err == nil)
	return
}

func (p *InfobloxProvider) deleteTXTRecord(o *ibcl.ObjectManager, name string) (res string, err error) {
	start := time.Now()
	res, err = o.DeleteTXTRecord(name)
	m.InfobloxObserveRequestDuration(start, metrics.DeleteTXTRecord, err == nil)
	return
}
