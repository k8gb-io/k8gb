package depresolver

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
	"sort"
	"strings"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
)

type DelegationZones []DelegationZoneInfo

type DelegationZoneInfo struct {
	Domain            string
	Zone              string
	ClusterNSName     string
	ExtClusterNSNames map[string]string
}

func parseDelegationZones(config *Config) []DelegationZoneInfo {

	zones := config.dnsZones
	edgeDNSZone := config.edgeDNSZone
	dnsZone := config.dnsZone

	getNsName := func(tag, edgeDNSServer, zone, edge string) string {
		if edgeDNSServer == localhost || edgeDNSServer == localhostIPv4 {
			return edgeDNSServer
		}
		const prefix = "gslb-ns"
		d := strings.TrimSuffix(zone, "."+edge)
		domainX := strings.ReplaceAll(d, ".", "-")
		return fmt.Sprintf("%s-%s-%s.%s", prefix, tag, domainX, edge)
	}

	// parse example.com:cloud.example.com;example.io:cloud.example.io into map[string]string
	getEnvAsArrayOfPairsOrFallback := func(zones string, fallback map[string]string) map[string]string {
		pairs := make(map[string]string)
		slice := strings.Split(zones, ";")
		if len(slice) == 0 {
			return fallback
		}
		for _, z := range slice {
			pair := strings.Split(z, ":")
			if len(pair) != 2 {
				return fallback
			}
			pairs[strings.Trim(pair[0], " ")] = strings.Trim(pair[1], " ")
		}
		for k, v := range fallback {
			if _, found := pairs[k]; !found {
				pairs[k] = v
			}
		}
		return pairs
	}
	var dzi []DelegationZoneInfo
	zones = strings.TrimSuffix(strings.TrimSuffix(zones, ";"), " ")
	fallbackDNSZone := map[string]string{}
	if !(edgeDNSZone == "" && dnsZone == "") {
		fallbackDNSZone[edgeDNSZone] = dnsZone
	}
	di := getEnvAsArrayOfPairsOrFallback(zones, fallbackDNSZone)

	for edge, zone := range di {
		zoneInfo := DelegationZoneInfo{
			Domain:        zone,
			Zone:          edge,
			ClusterNSName: getNsName(config.ClusterGeoTag, config.EdgeDNSServers[0].Host, zone, edge),
			ExtClusterNSNames: func(zone, edge string) map[string]string {
				m := map[string]string{}
				for _, tag := range config.extClustersGeoTags {
					m[tag] = getNsName(tag, config.EdgeDNSServers[0].Host, zone, edge)
				}
				return m
			}(zone, edge),
		}
		dzi = append(dzi, zoneInfo)
	}
	return dzi
}

// GetNSServerList returns a sorted list of all NS servers for the delegation zone
func (z *DelegationZoneInfo) GetNSServerList() []string {
	list := []string{z.ClusterNSName}
	for _, v := range z.ExtClusterNSNames {
		list = append(list, v)
	}
	sort.Strings(list)
	return list
}

// GetExternalDNSEndpointName returns name of endpoint sitting in k8gb namespace
func (z *DelegationZoneInfo) GetExternalDNSEndpointName() string {
	var suffix = strings.Trim(strings.ReplaceAll(z.Domain, ".", "-"), " ")
	return fmt.Sprintf("k8gb-ns-extdns-%s", suffix)
}

// FindByGslbStatusHostname returns DelegationZoneInfo for the hostname
func (d *DelegationZones) FindByGslbStatusHostname(gslb *k8gbv1beta1.Gslb) *DelegationZoneInfo {
	if len(gslb.Status.Servers) == 0 {
		return nil
	}
	for _, z := range *d {
		if strings.HasSuffix(gslb.Status.Servers[0].Host, z.Domain) {
			return &z
		}
	}
	return nil
}

func (d *DelegationZones) GetClusterNSNameByGslb(gslb *k8gbv1beta1.Gslb) string {
	z := d.FindByGslbStatusHostname(gslb)
	if z != nil {
		return z.ClusterNSName
	}
	return ""
}

func (d *DelegationZones) GetExternalClusterNSNamesByHostname(host string) map[string]string {
	z := d.getZone(host)
	if z != nil {
		return z.ExtClusterNSNames
	}
	return map[string]string{}
}

func (d *DelegationZones) ContainsZone(host string) bool {
	return d.getZone(host) != nil
}

func (d *DelegationZones) ListZones() []string {
	var zones []string
	for _, z := range *d {
		zones = append(zones, z.Zone)
	}
	return zones
}

func (d *DelegationZones) getZone(host string) *DelegationZoneInfo {
	for _, z := range *d {
		if strings.Contains(host, z.Zone) {
			return &z
		}
	}
	return nil
}
