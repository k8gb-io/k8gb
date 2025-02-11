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
	"os"
	"strconv"
	"strings"

	"github.com/k8gb-io/k8gb/controllers/utils"

	"github.com/AbsaOSS/env-binder/env"
	"github.com/rs/zerolog"
)

// TODO: refactor with kong.CLI to read envvars into Config
// TODO: refactor with go-playground/validator to validate

// Environment variables keys
const (
	ReconcileRequeueSecondsKey = "RECONCILE_REQUEUE_SECONDS"
	NSRecordTTLKey             = "NS_RECORD_TTL"
	ClusterGeoTagKey           = "CLUSTER_GEO_TAG"
	ExtClustersGeoTagsKey      = "EXT_GSLB_CLUSTERS_GEO_TAGS"
	ExtDNSEnabledKey           = "EXTDNS_ENABLED"
	EdgeDNSServersKey          = "EDGE_DNS_SERVERS"
	EdgeDNSZoneKey             = "EDGE_DNS_ZONE"
	DNSZoneKey                 = "DNS_ZONE"
	DNSZonesKey                = "DNS_ZONES"
	InfobloxGridHostKey        = "INFOBLOX_GRID_HOST"
	InfobloxVersionKey         = "INFOBLOX_WAPI_VERSION"
	InfobloxPortKey            = "INFOBLOX_WAPI_PORT"
	InfobloxUsernameKey        = "INFOBLOX_WAPI_USERNAME"
	// #nosec G101; ignore false positive gosec; see: https://securego.io/docs/rules/g101.html
	InfobloxPasswordKey            = "INFOBLOX_WAPI_PASSWORD"
	InfobloxHTTPRequestTimeoutKey  = "INFOBLOX_HTTP_REQUEST_TIMEOUT"
	InfobloxHTTPPoolConnectionsKey = "INFOBLOX_HTTP_POOL_CONNECTIONS"
	K8gbNamespaceKey               = "POD_NAMESPACE"
	CoreDNSExposedKey              = "COREDNS_EXPOSED"
	LogLevelKey                    = "LOG_LEVEL"
	LogFormatKey                   = "LOG_FORMAT"
	LogNoColorKey                  = "NO_COLOR"
	TracingEnabled                 = "TRACING_ENABLED"
	OtelExporterOtlpEndpoint       = "OTEL_EXPORTER_OTLP_ENDPOINT"
	TracingSamplingRatio           = "TRACING_SAMPLING_RATIO"
	MetricsAddressKey              = "METRICS_ADDRESS"
)

// Deprecated environment variables keys
const (
	// Deprecated: Please use EDGE_DNS_SERVERS instead.
	EdgeDNSServerKey = "EDGE_DNS_SERVER"

	// Deprecated: Please use EDGE_DNS_SERVERS instead.
	EdgeDNSServerPortKey = "EDGE_DNS_SERVER_PORT"
)

const (
	localhost = "localhost"

	localhostIPv4 = "127.0.0.1"
)

// ResolveOperatorConfig executes once. It reads operator's configuration
// from environment variables into &Config and validates
func (dr *DependencyResolver) ResolveOperatorConfig() (*Config, error) {
	var recognizedDNSTypes []EdgeDNSType
	dr.onceConfig.Do(func() {

		dr.config = &Config{}

		// binding
		dr.errorConfig = env.Bind(dr.config)
		if dr.errorConfig != nil {
			return
		}

		// calculation
		fallbackDNS := fmt.Sprintf("%s:%v", dr.config.fallbackEdgeDNSServerName, dr.config.fallbackEdgeDNSServerPort)
		edgeDNSServerList := env.GetEnvAsArrayOfStringsOrFallback(EdgeDNSServersKey, []string{fallbackDNS})
		dr.config.EdgeDNSServers = parseEdgeDNSServers(edgeDNSServerList)
		dr.config.ExtClustersGeoTags = excludeGeoTag(dr.config.ExtClustersGeoTags, dr.config.ClusterGeoTag)
		dr.config.Log.Level, _ = zerolog.ParseLevel(strings.ToLower(dr.config.Log.level))
		dr.config.Log.Format = parseLogOutputFormat(strings.ToLower(dr.config.Log.format))
		dr.config.EdgeDNSType, recognizedDNSTypes = getEdgeDNSType(dr.config)

		dr.errorConfig = dr.validateConfig(dr.config, recognizedDNSTypes)
		// validation
		if dr.errorConfig == nil {
			dr.config.DelegationZones = parseDelegationZones(dr.config)
		}
	})
	return dr.config, dr.errorConfig
}

func (dr *DependencyResolver) validateConfig(config *Config, recognizedDNSTypes []EdgeDNSType) (err error) {
	const dnsNameMax = 253
	const dnsLabelMax = 63
	if config.Log.Level == zerolog.NoLevel {
		return fmt.Errorf("invalid '%s', allowed values ['','%s','%s','%s','%s','%s','%s','%s']", LogLevelKey,
			zerolog.TraceLevel, zerolog.DebugLevel, zerolog.InfoLevel, zerolog.WarnLevel, zerolog.FatalLevel,
			zerolog.DebugLevel, zerolog.PanicLevel)
	}
	if config.Log.Format == NoFormat {
		return fmt.Errorf("invalid '%s', allowed values ['','%s','%s']", LogFormatKey, JSONFormat, SimpleFormat)
	}
	if config.EdgeDNSType == DNSTypeMultipleProviders {
		return fmt.Errorf("several EdgeDNS recognized %s", recognizedDNSTypes)
	}
	err = field(K8gbNamespaceKey, config.K8gbNamespace).isNotEmpty().matchRegexp(k8sNamespaceRegex).err
	if err != nil {
		return err
	}
	err = field(ReconcileRequeueSecondsKey, config.ReconcileRequeueSeconds).isHigherThanZero().err
	if err != nil {
		return err
	}
	err = field(NSRecordTTLKey, config.NSRecordTTL).isHigherThanZero().err
	if err != nil {
		return err
	}
	err = field(ClusterGeoTagKey, config.ClusterGeoTag).isNotEmpty().matchRegexp(geoTagRegex).err
	if err != nil {
		return err
	}
	err = field(ExtClustersGeoTagsKey, config.ExtClustersGeoTags).hasItems().hasUniqueItems().err
	if err != nil {
		return err
	}
	for i, geoTag := range config.ExtClustersGeoTags {
		err = field(fmt.Sprintf("%s[%v]", ExtClustersGeoTagsKey, i), geoTag).
			isNotEmpty().matchRegexp(geoTagRegex).err
		if err != nil {
			return err
		}
	}
	err = field(EdgeDNSServersKey, os.Getenv(EdgeDNSServersKey)).isNotEmpty().matchRegexp(hostNamesWithPortsRegex1).err
	if err != nil {
		return err
	}
	err = field(EdgeDNSServersKey, os.Getenv(EdgeDNSServersKey)).isNotEmpty().matchRegexp(hostNamesWithPortsRegex2).err
	if err != nil {
		return err
	}
	err = validateLocalhostNotAmongDNSServers(config)
	if err != nil {
		return err
	}
	err = field(EdgeDNSServersKey, config.EdgeDNSServers).isNotEmpty().matchRegexp(hostNamesWithPortsRegex1).err
	if err != nil {
		return err
	}
	for _, s := range config.EdgeDNSServers {
		if s.Port < 1 || s.Port > 65535 {
			return fmt.Errorf("error for port of edge dns server(%v): it must be a positive integer between 1 and 65535", s)
		}
	}
	err = field(EdgeDNSZoneKey, config.EdgeDNSZone).isNotEmpty().matchRegexp(hostNameRegex).err
	if err != nil {
		return err
	}
	err = field(DNSZoneKey, config.DNSZone).isNotEmpty().matchRegexp(hostNameRegex).err
	if err != nil {
		return err
	}
	// do full Infoblox validation only in case that Host exists
	if isNotEmpty(config.Infoblox.Host) {
		err = validateConfigForInfoblox(config)
		if err != nil {
			return err
		}
	}
	validateLabels := func(label string) error {
		labels := strings.Split(label, ".")
		for _, l := range labels {
			if len(l) > dnsLabelMax {
				return fmt.Errorf("%s exceeds %v characters limit", l, dnsLabelMax)
			}
		}
		return nil
	}

	serverNames := config.GetExternalClusterNSNames()
	serverNames[config.ClusterGeoTag] = config.GetClusterNSName()
	for geoTag, nsName := range serverNames {
		if len(nsName) > dnsNameMax {
			return fmt.Errorf("ns name '%s' exceeds %v charactes limit for [GeoTag: '%s', %s: '%s', %s: '%s']",
				nsName, dnsLabelMax, geoTag, EdgeDNSZoneKey, config.EdgeDNSZone, DNSZoneKey, config.DNSZone)
		}
		if err := validateLabels(nsName); err != nil {
			return fmt.Errorf("error for geo tag: %s. %s in ns name %s", geoTag, err, nsName)
		}
	}

	mHost, mPort, err := parseMetricsAddr(config.MetricsAddress)
	if err != nil {
		return fmt.Errorf("invalid %s: expecting MetricsAddress in form {host}:port (%s)", MetricsAddressKey, err)
	}
	err = field(MetricsAddressKey, mHost).matchRegexps(hostNameRegex, ipAddressRegex).err
	if err != nil {
		return err
	}
	err = field(MetricsAddressKey, mPort).isLessOrEqualTo(65535).isHigherThan(1024).err
	if err != nil {
		return err
	}
	return nil
}

func validateLocalhostNotAmongDNSServers(config *Config) error {
	containsLocalhost := func(list utils.DNSList) bool {
		for i := 1; i < len(list); i++ { // skipping first because localhost or 127.0.0.1 can occur on the first position
			if list[i].Host == localhost || list[i].Host == localhostIPv4 {
				return true
			}
		}
		return false
	}
	if len(config.EdgeDNSServers) > 1 && containsLocalhost(config.EdgeDNSServers) {
		return fmt.Errorf("invalid %s: the list can't contain 'localhost' or '127.0.0.1' on other than the first position", EdgeDNSServersKey)
	}
	return nil
}

func validateConfigForInfoblox(config *Config) error {
	err := field(InfobloxGridHostKey, config.Infoblox.Host).matchRegexps(hostNameRegex, ipAddressRegex).err
	if err != nil {
		return err
	}
	err = field(InfobloxVersionKey, config.Infoblox.Version).isNotEmpty().matchRegexp(versionNumberRegex).err
	if err != nil {
		return err
	}
	err = field(InfobloxPortKey, config.Infoblox.Port).isHigherThanZero().isLessOrEqualTo(65535).err
	if err != nil {
		return err
	}
	err = field(InfobloxUsernameKey, config.Infoblox.Username).isNotEmpty().err
	if err != nil {
		return err
	}
	err = field(InfobloxPasswordKey, config.Infoblox.Password).isNotEmpty().err
	if err != nil {
		return err
	}
	err = field(InfobloxHTTPPoolConnectionsKey, config.Infoblox.HTTPPoolConnections).isHigherOrEqualToZero().err
	if err != nil {
		return err
	}
	err = field(InfobloxHTTPRequestTimeoutKey, config.Infoblox.HTTPRequestTimeout).isHigherThanZero().err
	if err != nil {
		return err
	}
	return nil
}

func (dr *DependencyResolver) GetDeprecations() (deprecations []string) {
	type oldVar = string
	type newVar struct {
		Name string
		Msg  string
	}

	var deprecated = map[oldVar]newVar{
		EdgeDNSServerKey: newVar{
			Name: EdgeDNSServersKey,
			Msg:  "Pass the hostname or IP address as comma-separated list",
		},
		EdgeDNSServerPortKey: newVar{
			Name: EdgeDNSServersKey,
			Msg: "Port is an optional item in the comma-separated list of dns edge servers, in following form: dns1:53,dns2 (if not provided after the " +
				"hostname and colon, it defaults to '53')",
		},
	}

	for k, v := range deprecated {
		if os.Getenv(k) != "" {
			deprecations = append(deprecations, fmt.Sprintf("'%s' has been deprecated, use %s instead. Details: %s", k, v.Name, v.Msg))
		}
	}
	return
}

func parseMetricsAddr(metricsAddr string) (host string, port int, err error) {
	ma := strings.Split(metricsAddr, ":")
	if len(ma) != 2 {
		err = fmt.Errorf("invalid format {host}:port (%s)", metricsAddr)
		return
	}
	host = ma[0]
	port, err = strconv.Atoi(ma[1])
	return
}

func parseEdgeDNSServers(serverList []string) (r []utils.DNSServer) {
	r = []utils.DNSServer{}
	var host, portStr string
	var err error
	for _, chunk := range serverList {
		chunk = strings.TrimSpace(chunk)
		switch strings.Count(chunk, ":") {
		case 0: // ipv4 or domain w/o port
			host = chunk
			portStr = "53"
		case 1: // ipv4 or domain w/ port
			host, portStr, err = net.SplitHostPort(chunk)
			if err != nil {
				continue
			}
		default: // ipv6 or http://foo:bar
			// not supported
			continue
		}
		var port int
		port, err = strconv.Atoi(portStr)
		if err != nil {
			port = 53
		}
		if host != "" {
			r = append(r, utils.DNSServer{
				Host: host,
				Port: port,
			})
		}
	}
	return r
}

// excludeGeoTag excludes the clusterGeoTag from external geo tags
func excludeGeoTag(tags []string, tag string) (r []string) {
	r = []string{}
	for _, t := range tags {
		if tag != t {
			r = append(r, t)
		}
	}
	return
}

// getEdgeDNSType contains logic retrieving EdgeDNSType.
func getEdgeDNSType(config *Config) (EdgeDNSType, []EdgeDNSType) {
	recognized := make([]EdgeDNSType, 0)
	if config.extDNSEnabled {
		recognized = append(recognized, DNSTypeExternal)
	}
	if isNotEmpty(config.Infoblox.Host) {
		recognized = append(recognized, DNSTypeInfoblox)
	}
	switch len(recognized) {
	case 0:
		return DNSTypeNoEdgeDNS, recognized
	case 1:
		return recognized[0], recognized
	}
	return DNSTypeMultipleProviders, recognized
}

func parseLogOutputFormat(value string) LogFormat {
	switch value {
	case json:
		return JSONFormat
	case simple:
		return SimpleFormat
	}
	return NoFormat
}

func (c *Config) GetExternalClusterNSNames() (m map[string]string) {
	m = make(map[string]string, len(c.ExtClustersGeoTags))
	for _, tag := range c.ExtClustersGeoTags {
		m[tag] = getNsName(tag, c.DNSZone, c.EdgeDNSZone, c.EdgeDNSServers[0].Host)
	}
	return
}

func (c *Config) GetClusterNSName() string {
	return getNsName(c.ClusterGeoTag, c.DNSZone, c.EdgeDNSZone, c.EdgeDNSServers[0].Host)
}

// getNsName returns NS for geo tag.
// The values is combination of DNSZone, EdgeDNSZone and (Ext)ClusterGeoTag, see:
// DNS_ZONE k8gb-test.gslb.cloud.example.com
// EDGE_DNS_ZONE: cloud.example.com
// CLUSTER_GEOTAG: us
// will generate "gslb-ns-us-k8gb-test-gslb.cloud.example.com"
// If edgeDNSServer == localhost or 127.0.0.1 than edgeDNSServer is returned.
// The function is private and expects only valid inputs.
func getNsName(tag, dnsZone, edgeDNSZone, edgeDNSServer string) string {
	if edgeDNSServer == "127.0.0.1" || edgeDNSServer == "localhost" {
		return edgeDNSServer
	}
	const prefix = "gslb-ns"
	d := strings.TrimSuffix(dnsZone, "."+edgeDNSZone)
	domainX := strings.ReplaceAll(d, ".", "-")
	return fmt.Sprintf("%s-%s-%s.%s", prefix, tag, domainX, edgeDNSZone)
}
