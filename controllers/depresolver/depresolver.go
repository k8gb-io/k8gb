// Package depresolver abstracts and implements k8gb dependencies resolver.
// depresolver responsibilities
// - abstracts multiple configurations into single point of access
// - provides predefined values when configuration is missing
// - validates configuration
// - executes once
// TODO: Add the rest of configuration to be resolved
package depresolver

import (
	"context"
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

//Config holds operator configuration
type Config struct {
	//Reschedule of Reconcile loop to pickup external Gslb targets
	ReconcileRequeueSeconds int
	// Cluster Geo Tag to determine specific location
	ClusterGeoTag string
	// Route53 switch
	Route53Enabled bool
}

//DependencyResolver resolves configuration for GSLB
type DependencyResolver struct {
	client      client.Client
	config      *Config
	context     context.Context
	onceConfig  sync.Once
	onceSpec    sync.Once
	errorConfig error
	errorSpec   error
}

const (
	lessOrEqualToZeroErrorMessage = "\"%s is less or equal to zero\""
	lessThanZeroErrorMessage      = "\"%s is less than zero\""
	doesNotMatchRegexMessage      = "\"%s does not match /%s/ regexp rule\""
)

//NewDependencyResolver returns a new depresolver.DependencyResolver
func NewDependencyResolver(context context.Context, client client.Client) *DependencyResolver {
	resolver := new(DependencyResolver)
	resolver.client = client
	resolver.context = context
	return resolver
}
