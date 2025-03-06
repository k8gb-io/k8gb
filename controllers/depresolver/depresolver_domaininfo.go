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
	"net"
	"sort"
	"strconv"
	"strings"
)

type DelegationZones []*DelegationZoneInfo

type DelegationZoneInfo struct {
	Domain            string // cloud.example.com
	Zone              string // example.com
	NegativeTTL       int
	ClusterNSName     string
	ExtClusterNSNames map[string]string
	IPs               []string
}

func parseDelegationZones(config *Config) ([]*DelegationZoneInfo, error) {
	type info struct {
		domain string
		zone   string
		negTTL string
	}

	zones := config.dnsZones

	getNsName := func(tag, zone, edge string) string {
		const prefix = "gslb-ns"
		d := strings.TrimSuffix(zone, "."+edge)
		domainX := strings.ReplaceAll(d, ".", "-")
		return fmt.Sprintf("%s-%s-%s.%s", prefix, tag, domainX, edge)
	}

	// parse example.com:cloud.example.com:30;example.io:cloud.example.io:50
	getEnvAsArrayOfPairsOrFallback := func(zones string) ([]info, error) {
		tuples := make([]info, 0)
		slice := strings.Split(zones, ";")
		for _, z := range slice {
			touple := strings.Split(z, ":")
			if len(touple) != 3 {
				return tuples, fmt.Errorf("invalid format of delegation zones: %s", z)
			}
			tuples = append(tuples, info{zone: strings.Trim(touple[0], " "), domain: strings.Trim(touple[1], " "), negTTL: strings.Trim(touple[2], " ")})
		}
		return tuples, nil
	}
	var dzi []*DelegationZoneInfo
	zones = strings.TrimSuffix(strings.TrimSuffix(zones, ";"), " ")
	di, err := getEnvAsArrayOfPairsOrFallback(zones)
	if err != nil {
		return dzi, err
	}

	for _, inf := range di {
		negTTL, err := strconv.Atoi(inf.negTTL)
		if err != nil {
			return dzi, fmt.Errorf("invalid value of delegation zones: %s", zones)
		}
		zoneInfo := &DelegationZoneInfo{
			Domain:        inf.domain,
			Zone:          inf.zone,
			NegativeTTL:   negTTL,
			ClusterNSName: getNsName(config.ClusterGeoTag, inf.domain, inf.zone),
			ExtClusterNSNames: func(zone, edge string) map[string]string {
				m := map[string]string{}
				for _, tag := range config.extClustersGeoTags {
					m[tag] = getNsName(tag, zone, edge)
				}
				return m
			}(inf.domain, inf.zone),
		}
		dzi = append(dzi, zoneInfo)
	}
	return dzi, nil
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

func (z *DelegationZoneInfo) GetSortedIPs() []string {
	// ipLess compares two IPs and returns true if ip1 < ip2
	ipLess := func(ip1, ip2 net.IP) bool {
		if ip1 == nil || ip2 == nil {
			return false
		}
		return ip1.To16() != nil && ip2.To16() != nil && ip1.String() < ip2.String()
	}

	sort.Slice(z.IPs, func(i, j int) bool {
		return ipLess(net.ParseIP(z.IPs[i]), net.ParseIP(z.IPs[j]))
	})

	return z.IPs
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
			return z
		}
	}
	return nil
}

func (d *DelegationZones) SetIPs(ips []string) {
	for _, z := range *d {
		z.IPs = ips
	}
}
