// Package depresolver abstracts and implements k8gb dependencies resolver.
// depresolver responsibilities
// - abstracts multiple configurations into single point of access
// - provides predefined values when configuration is missing
// - validates configuration
// - executes once
package depresolver

import (
	"sync"

	"github.com/AbsaOSS/k8gb/api/v1beta1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EdgeDNSType specifies to which edge DNS is k8gb connecting
type EdgeDNSType int

const (
	// DNSTypeNoEdgeDNS is default DNSType. Is used during integration testing when no edgeDNS provider exists
	DNSTypeNoEdgeDNS EdgeDNSType = 1 << iota
	// DNSTypeInfoblox type
	DNSTypeInfoblox
	// DNSTypeRoute53 type
	DNSTypeRoute53
	// DNSTypeNS1 type
	DNSTypeNS1
)

// Infoblox configuration
// TODO: consider to make this private after refactor
type Infoblox struct {
	// Host
	Host string
	// Version
	Version string
	// Port
	Port int
	// Username
	Username string
	// Password
	Password string
	// HTTPRequestTimeout seconds; default = 20
	HTTPRequestTimeout int
	// HTTPPoolConnections seconds; default = 10
	HTTPPoolConnections int
}

// Override configuration
type Override struct {
	// FakeDNSEnabled; default=false
	FakeDNSEnabled bool
	// FakeInfobloxEnabled if true than Infoblox connection FQDN=`fakezone.example.com`; default = false
	FakeInfobloxEnabled bool
}

// Config is operator configuration returned by depResolver
type Config struct {
	// Reschedule of Reconcile loop to pickup external Gslb targets
	ReconcileRequeueSeconds int
	// ClusterGeoTag to determine specific location
	ClusterGeoTag string
	// ExtClustersGeoTags to identify clusters in other locations in format separated by comma. i.e.: "eu,uk,us"
	ExtClustersGeoTags []string
	// EdgeDNSType is READONLY and is set automatically by configuration
	EdgeDNSType EdgeDNSType
	// EdgeDNSServer
	EdgeDNSServer string
	// EdgeDNSZone main zone which would contain gslb zone to delegate; e.g. example.com
	EdgeDNSZone string
	// DNSZone controlled by gslb; e.g. cloud.example.com
	DNSZone string
	// K8gbNamespace k8gb namespace
	K8gbNamespace string
	// Infoblox configuration
	Infoblox Infoblox
	// Override the behavior of GSLB in the test environments
	Override Override
	// route53Enabled hidden. EdgeDNSType defines all enabled Enabled types
	route53Enabled bool
	// ns1Enabled flag
	ns1Enabled bool
	// CoreDNSExposed flag
	CoreDNSExposed bool
}

// DependencyResolver resolves configuration for GSLB
type DependencyResolver struct {
	client      client.Client
	config      *Config
	onceConfig  sync.Once
	errorConfig error
	errorSpec   error
	spec        v1beta1.GslbSpec
}

// NewDependencyResolver returns a new depresolver.DependencyResolver
func NewDependencyResolver(client client.Client) *DependencyResolver {
	resolver := new(DependencyResolver)
	resolver.client = client
	return resolver
}
