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
