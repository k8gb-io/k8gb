// Package depresolver abstracts and implements k8gb dependencies resolver.
// depresolver responsibilities
// - abstracts multiple configurations into single point of access
// - provides predefined values when configuration is missing
// - validates configuration
// - executes once
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
	"sync"

	"github.com/k8gb-io/k8gb/controllers/utils"

	"github.com/k8gb-io/k8gb/api/v1beta1"
	"github.com/rs/zerolog"
)

// LogFormat specifies how the logger prints values
type LogFormat int8

const (
	// JSONFormat prints messages as single json record
	JSONFormat LogFormat = 1 << iota
	// SimpleFormat prints messages in human readable way
	SimpleFormat
	// NoFormat, returned in situation when format is not recognised
	NoFormat
)

const (
	json         = "json"
	simple       = "simple"
	unrecognised = "noformat"
)

func (f LogFormat) String() string {
	switch f {
	case JSONFormat:
		return json
	case SimpleFormat:
		return simple
	}
	return unrecognised
}

// EdgeDNSType specifies to which edge DNS is k8gb connecting
type EdgeDNSType string

const (
	// DNSTypeNoEdgeDNS is default DNSType. Is used during integration testing when no edgeDNS provider exists
	DNSTypeNoEdgeDNS EdgeDNSType = "NoEdgeDNS"
	// DNSTypeInfoblox type
	DNSTypeInfoblox EdgeDNSType = "Infoblox"
	// DNSTypeRoute53 type
	DNSTypeExternal EdgeDNSType = "ExtDNS"
	// DNSTypeMultipleProviders type
	DNSTypeMultipleProviders EdgeDNSType = "MultipleProviders"
)

const (
	// GeoIP strategy
	GeoStrategy = "geoip"
	// RoundRobin strategy
	RoundRobinStrategy = "roundRobin"
	// Failover strategy
	FailoverStrategy = "failover"
)

// Log configuration
type Log struct {
	// Level [panic, fatal, error,warn,info,debug,trace], defines level of logger, default: info
	Level zerolog.Level
	// Format [simple,json] specifies how the logger prints values
	Format LogFormat
	// NoColor prints colored output if Format == simple
	NoColor bool `env:"NO_COLOR, default=false"`
	// format is binding source for Format
	format string `env:"LOG_FORMAT, default=simple"`
	// level is binding source for Level
	level string `env:"LOG_LEVEL, default=info"`
}

// Infoblox configuration
type Infoblox struct {
	// Host
	Host string `env:"INFOBLOX_GRID_HOST"`
	// Version
	Version string `env:"INFOBLOX_WAPI_VERSION"`
	// Port
	Port int `env:"INFOBLOX_WAPI_PORT, default=0"`
	// Username
	Username string `env:"INFOBLOX_WAPI_USERNAME"`
	// Password
	Password string `env:"INFOBLOX_WAPI_PASSWORD"`
	// HTTPRequestTimeout seconds
	HTTPRequestTimeout int `env:"INFOBLOX_HTTP_REQUEST_TIMEOUT, default=20"`
	// HTTPPoolConnections seconds
	HTTPPoolConnections int `env:"INFOBLOX_HTTP_POOL_CONNECTIONS, default=10"`
}

// Config is operator configuration returned by depResolver
type Config struct {
	// Reschedule of Reconcile loop to pickup external Gslb targets
	ReconcileRequeueSeconds int `env:"RECONCILE_REQUEUE_SECONDS, default=30"`
	// TTL of the NS and respective glue record used by external DNS
	NSRecordTTL int `env:"NS_RECORD_TTL, default=30"`
	// ClusterGeoTag to determine specific location
	ClusterGeoTag string `env:"CLUSTER_GEO_TAG"`
	// extClustersGeoTags to identify clusters in other locations in format separated by comma. i.e.: "eu,uk,us"
	extClustersGeoTags []string `env:"EXT_GSLB_CLUSTERS_GEO_TAGS, default=[]"`
	// EdgeDNSType is READONLY and is set automatically by configuration
	EdgeDNSType EdgeDNSType
	// EdgeDNSServers
	EdgeDNSServers utils.DNSList
	// to avoid breaking changes is used as fallback server for EdgeDNSServers
	fallbackEdgeDNSServerName string `env:"EDGE_DNS_SERVER"`
	// to avoid breaking changes is used as fallback server port for EdgeDNSServers
	fallbackEdgeDNSServerPort int `env:"EDGE_DNS_SERVER_PORT, default=53"`
	// edgeDNSZone main zone which would contain gslb zone to delegate; e.g. example.com
	edgeDNSZone string `env:"EDGE_DNS_ZONE"`
	// dnsZone controlled by gslb; e.g. cloud.example.com
	dnsZone string `env:"DNS_ZONE"`
	// DelegationZones
	DelegationZones DelegationZones
	// DelegationZones pairs of dnsZone ad edgeDNSZone, eg: DNS_ZONES=example.com:cloud.example.com;example.io:cloud.example.io
	dnsZones string `env:"DNS_ZONES"`
	// K8gbNamespace k8gb namespace
	K8gbNamespace string `env:"POD_NAMESPACE"`
	// Infoblox configuration
	Infoblox Infoblox
	// CoreDNSExposed flag
	CoreDNSExposed bool `env:"COREDNS_EXPOSED, default=false"`
	// IngressPath if not CoreDNSExposed the IngressPath must be set to get exposed IP's. Any ingress containing exposed IP's can be used.
	IngressPath string `env:"INGRESS_PATH"`
	// Log configuration
	Log Log
	// MetricsAddress in format address:port where address can be empty, IP address, or hostname, default: 0.0.0.0:8080
	MetricsAddress string `env:"METRICS_ADDRESS, default=0.0.0.0:8080"`
	// extDNSEnabled hidden. EdgeDNSType defines all enabled Enabled types
	extDNSEnabled bool `env:"EXTDNS_ENABLED, default=false"`
	// TracingEnabled flag decides whether to use a real otlp tracer or a noop one
	TracingEnabled bool `env:"TRACING_ENABLED, default=false"`
	// TracingSamplingRatio how many traces should be kept and sent (1.0 - all, 0.0 - none)
	TracingSamplingRatio float64 `env:"TRACING_SAMPLING_RATIO, default=1.0"`
	// OtelExporterOtlpEndpoint where the traces should be sent to (in case of otel collector deployed on the same pod as sidecar -> localhost:4318)
	// otel collector itself can be configured via a configmap to send it somewhere else
	OtelExporterOtlpEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT, default=localhost:4318"`
}

// DependencyResolver resolves configuration for GSLB
type DependencyResolver struct {
	config      *Config
	onceConfig  sync.Once
	errorConfig error
	errorSpec   error
	spec        v1beta1.GslbSpec
}

// NewDependencyResolver returns a new depresolver.DependencyResolver
func NewDependencyResolver() *DependencyResolver {
	resolver := new(DependencyResolver)
	return resolver
}
