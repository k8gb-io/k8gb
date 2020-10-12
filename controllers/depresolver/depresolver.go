// Package depresolver abstracts and implements k8gb dependencies resolver.
// depresolver responsibilities
// - abstracts multiple configurations into single point of access
// - provides predefined values when configuration is missing
// - validates configuration
// - executes once
package depresolver

import (
	"context"
	"sync"

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
)

//Infoblox configuration
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
	// DNSTypeRoute53 switch
	// TODO: hide for depresolver subscriber as depresolver retrieves EdgeDNSType. Maybe we can change configuration and set EdgeDNSType directly instead of DNSTypeRoute53 boolean
	Route53Enabled bool
	// Infoblox configuration
	Infoblox Infoblox
}

// DependencyResolver resolves configuration for GSLB
type DependencyResolver struct {
	client      client.Client
	config      *Config
	context     context.Context
	onceConfig  sync.Once
	onceSpec    sync.Once
	errorConfig error
	errorSpec   error
}

// NewDependencyResolver returns a new depresolver.DependencyResolver
func NewDependencyResolver(context context.Context, client client.Client) *DependencyResolver {
	resolver := new(DependencyResolver)
	resolver.client = client
	resolver.context = context
	return resolver
}
