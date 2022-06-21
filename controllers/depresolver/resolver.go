package depresolver

import (
	"context"
	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Resolver interface {
	// ResolveOperatorConfig executes once. It reads operator's configuration
	// from environment variables into &Config and validates
	ResolveOperatorConfig() (*Config, error)

	// ResolveGslbSpec fills Gslb by spec values. It executes always, when gslb is initialised.
	// If spec value is not defined, it will use the default value. Function returns error if input is invalid.
	ResolveGslbSpec(ctx context.Context, gslb *k8gbv1beta1.Gslb, client client.Client) error
}
