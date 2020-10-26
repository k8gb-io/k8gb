package depresolver

import (
	"fmt"

	"github.com/AbsaOSS/k8gb/controllers/internal/env"
)

// Environment variables keys
const (
	ReconcileRequeueSecondsKey = "RECONCILE_REQUEUE_SECONDS"
	ClusterGeoTagKey           = "CLUSTER_GEO_TAG"
	ExtClustersGeoTagsKey      = "EXT_GSLB_CLUSTERS_GEO_TAGS"
	Route53EnabledKey          = "ROUTE53_ENABLED"
	EdgeDNSServerKey           = "EDGE_DNS_SERVER"
	EdgeDNSZoneKey             = "EDGE_DNS_ZONE"
	DNSZoneKey                 = "DNS_ZONE"
	InfobloxGridHostKey        = "INFOBLOX_GRID_HOST"
	InfobloxVersionKey         = "INFOBLOX_WAPI_VERSION"
	InfobloxPortKey            = "INFOBLOX_WAPI_PORT"
	InfobloxUsernameKey        = "EXTERNAL_DNS_INFOBLOX_WAPI_USERNAME"
	// #nosec G101; ignore false positive gosec; see: https://securego.io/docs/rules/g101.html
	InfobloxPasswordKey     = "EXTERNAL_DNS_INFOBLOX_WAPI_PASSWORD"
	OverrideWithFakeDNSKey  = "OVERRIDE_WITH_FAKE_EXT_DNS"
	OverrideFakeInfobloxKey = "FAKE_INFOBLOX"
)

// ResolveOperatorConfig executes once. It reads operator's configuration
// from environment variables into &Config and validates
func (dr *DependencyResolver) ResolveOperatorConfig() (*Config, error) {
	dr.onceConfig.Do(func() {
		dr.config = &Config{}
		dr.config.ReconcileRequeueSeconds, _ = env.GetEnvAsIntOrFallback(ReconcileRequeueSecondsKey, 30)
		dr.config.ClusterGeoTag = env.GetEnvAsStringOrFallback(ClusterGeoTagKey, "")
		dr.config.ExtClustersGeoTags = env.GetEnvAsArrayOfStringsOrFallback(ExtClustersGeoTagsKey, []string{})
		dr.config.Route53Enabled = env.GetEnvAsBoolOrFallback(Route53EnabledKey, false)
		dr.config.EdgeDNSServer = env.GetEnvAsStringOrFallback(EdgeDNSServerKey, "")
		dr.config.EdgeDNSZone = env.GetEnvAsStringOrFallback(EdgeDNSZoneKey, "")
		dr.config.DNSZone = env.GetEnvAsStringOrFallback(DNSZoneKey, "")
		dr.config.Infoblox.Host = env.GetEnvAsStringOrFallback(InfobloxGridHostKey, "")
		dr.config.Infoblox.Version = env.GetEnvAsStringOrFallback(InfobloxVersionKey, "")
		dr.config.Infoblox.Port, _ = env.GetEnvAsIntOrFallback(InfobloxPortKey, 0)
		dr.config.Infoblox.Username = env.GetEnvAsStringOrFallback(InfobloxUsernameKey, "")
		dr.config.Infoblox.Password = env.GetEnvAsStringOrFallback(InfobloxPasswordKey, "")
		dr.config.Override.FakeDNSEnabled = env.GetEnvAsBoolOrFallback(OverrideWithFakeDNSKey, false)
		dr.config.Override.FakeInfobloxEnabled = env.GetEnvAsBoolOrFallback(OverrideFakeInfobloxKey, false)
		dr.errorConfig = dr.validateConfig(dr.config)
		dr.config.EdgeDNSType = getEdgeDNSType(dr.config)
	})
	return dr.config, dr.errorConfig
}

func (dr *DependencyResolver) validateConfig(config *Config) (err error) {
	err = field("reconcileRequeueSeconds", config.ReconcileRequeueSeconds).isHigherThanZero().err
	if err != nil {
		return err
	}
	err = field("clusterGeoTag", config.ClusterGeoTag).isNotEmpty().matchRegexp(geoTagRegex).err
	if err != nil {
		return err
	}
	err = field("extClusterGeoTags", config.ExtClustersGeoTags).hasItems().hasUniqueItems().err
	if err != nil {
		return err
	}
	for i, geoTag := range config.ExtClustersGeoTags {
		err = field(fmt.Sprintf("extClustersGeoTags[%v]", i), geoTag).isNotEmpty().matchRegexp(geoTagRegex).isNotEqualTo(config.ClusterGeoTag).err
		if err != nil {
			return err
		}
	}
	err = field("edgeDNSServer", config.EdgeDNSServer).isNotEmpty().matchRegexps(hostNameRegex, ipAddressRegex).err
	if err != nil {
		return err
	}
	err = field("edgeDNSZone", config.EdgeDNSZone).isNotEmpty().matchRegexp(hostNameRegex).err
	if err != nil {
		return err
	}
	err = field("DNSZone", config.DNSZone).isNotEmpty().matchRegexp(hostNameRegex).err
	if err != nil {
		return err
	}
	// do full Infoblox validation only in case that Host exists
	if isNotEmpty(config.Infoblox.Host) {
		err = field("InfobloxGridHost", config.Infoblox.Host).matchRegexps(hostNameRegex, ipAddressRegex).err
		if err != nil {
			return err
		}
		err = field("InfobloxVersion", config.Infoblox.Version).isNotEmpty().matchRegexp(versionNumberRegex).err
		if err != nil {
			return err
		}
		err = field("InfobloxPort", config.Infoblox.Port).isHigherThanZero().isLessOrEqualTo(65535).err
		if err != nil {
			return err
		}
		err = field("InfobloxUsername", config.Infoblox.Username).isNotEmpty().err
		if err != nil {
			return err
		}
		err = field("InfobloxPassword", config.Infoblox.Password).isNotEmpty().err
		if err != nil {
			return err
		}
	}
	return nil
}

// getEdgeDNSType contains logic retrieving EdgeDNSType
func getEdgeDNSType(config *Config) EdgeDNSType {
	var t = DNSTypeNoEdgeDNS
	if config.Route53Enabled {
		t = t | DNSTypeRoute53
	}
	if isNotEmpty(config.Infoblox.Host) {
		t = t | DNSTypeInfoblox
	}
	if t > DNSTypeNoEdgeDNS {
		t = t - DNSTypeNoEdgeDNS
	}
	return t
}
