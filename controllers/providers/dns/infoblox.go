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
	"reflect"
	"time"

	ibcl "github.com/infobloxopen/infoblox-go-client"
	"github.com/k8gb-io/k8gb/controllers/depresolver"
	"github.com/k8gb-io/k8gb/controllers/providers/metrics"
	"github.com/rs/zerolog"
)

type InfobloxProvider struct {
	config depresolver.Config
	client InfobloxClient
	logger *zerolog.Logger
}

var m = metrics.Metrics()

func NewInfobloxDNS(config depresolver.Config, client InfobloxClient, logger *zerolog.Logger) *InfobloxProvider {
	return &InfobloxProvider{
		client: client,
		config: config,
		logger: logger,
	}
}

// current IP list is up to date, so we remove it from delegatedTo.
func (p *InfobloxProvider) sanitizeDelegateZone(local, upstream []ibcl.NameServer, zoneInfo *depresolver.DelegationZoneInfo) []ibcl.NameServer {
	// Drop own records for straight away update
	// And ensure local entries are up to date
	// And final list is sorted
	final := local
	remote := p.filterOutDelegateTo(upstream, zoneInfo.ClusterNSName)
	final = append(final, remote...)
	sortZones(final)

	return final
}

func (p *InfobloxProvider) CreateZoneDelegation(zoneInfo *depresolver.DelegationZoneInfo) error {
	objMgr, err := p.client.GetObjectManager()
	if err != nil {
		return err
	}
	findZone, err := p.getZoneDelegated(objMgr, zoneInfo.Domain)
	if err != nil {
		return err
	}

	var delegateTo []ibcl.NameServer
	for _, address := range zoneInfo.GetSortedIPs() {
		nameServer := ibcl.NameServer{Address: address, Name: zoneInfo.ClusterNSName}
		delegateTo = append(delegateTo, nameServer)
	}

	if findZone == nil {
		p.logger.Info().
			Str("DNSZone", zoneInfo.Domain).
			Msg("Creating delegated zone")
		p.logger.Debug().
			Interface("records", delegateTo).
			Msg("Delegated records")
		_, err = p.createZoneDelegated(objMgr, zoneInfo.Domain, delegateTo)
		if err != nil {
			return err
		}
		return nil
	}

	// if zone exists
	if len(findZone.Ref) > 0 {
		sortZones(findZone.DelegateTo)
		currentList := p.sanitizeDelegateZone(delegateTo, findZone.DelegateTo, zoneInfo)
		if !reflect.DeepEqual(findZone.DelegateTo, currentList) {
			p.logger.Info().
				Interface("records", findZone.DelegateTo).
				Msg("Found delegated zone records")
			p.logger.Info().
				Str("DNSZone", zoneInfo.Domain).
				Interface("serverList", currentList).
				Msg("Updating delegated zone with the server list")
			_, err = p.updateZoneDelegated(objMgr, findZone.Ref, currentList)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *InfobloxProvider) Finalize(zoneInfo *depresolver.DelegationZoneInfo) error {
	p.logger.Info().Msgf("Domain %s must deleted by manually in Infoblox", zoneInfo.Domain)
	return nil
}

func (p *InfobloxProvider) String() string {
	return "Infoblox"
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

// func (p *InfobloxProvider) deleteZoneDelegated(o *ibcl.ObjectManager, fqdn string) (res string, err error) {
//	start := time.Now()
//	res, err = o.DeleteZoneDelegated(fqdn)
//	m.InfobloxObserveRequestDuration(start, metrics.DeleteZoneDelegated, err == nil)
//	return
// }
