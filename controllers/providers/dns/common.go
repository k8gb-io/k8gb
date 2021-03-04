/*
Copyright 2021 Absa Group Limited

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dns

import (
	"fmt"
	"strings"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"

	"github.com/AbsaOSS/k8gb/controllers/depresolver"
)

func nsServerName(config depresolver.Config) string {
	dnsZoneIntoNS := strings.ReplaceAll(config.DNSZone, ".", "-")
	return fmt.Sprintf("gslb-ns-%s-%s.%s", dnsZoneIntoNS, config.ClusterGeoTag, config.EdgeDNSZone)
}

func nsServerNameExt(config depresolver.Config) (extNSServers []string) {
	dnsZoneIntoNS := strings.ReplaceAll(config.DNSZone, ".", "-")
	extNSServers = []string{}
	for _, clusterGeoTag := range config.ExtClustersGeoTags {

		extNSServers = append(extNSServers,
			fmt.Sprintf("gslb-ns-%s-%s.%s", dnsZoneIntoNS, clusterGeoTag, config.EdgeDNSZone))
	}
	return extNSServers
}

func getExternalClusterHeartbeatFQDNs(gslb *k8gbv1beta1.Gslb, config depresolver.Config) (extGslbClusters []string) {
	for _, geoTag := range config.ExtClustersGeoTags {
		extGslbClusters = append(extGslbClusters, fmt.Sprintf("%s-heartbeat-%s.%s", gslb.Name, geoTag, config.EdgeDNSZone))
	}
	return
}
