package depresolver

import (
	"fmt"

	"k8s.io/kubernetes/pkg/util/env"
)

const (
	reconcileRequeueSecondsKey = "RECONCILE_REQUEUE_SECONDS"
)

// ResolveOperatorConfig executes once. It reads operator's configuration
// from environment variables into &Config and validates
func (dr *DependencyResolver) ResolveOperatorConfig() (*Config, error) {
	dr.onceConfig.Do(func() {
		dr.config = &Config{}
		// set predefined values when not read from the environment variables
		dr.config.ReconcileRequeueSeconds, _ = env.GetEnvAsIntOrFallback(reconcileRequeueSecondsKey, 30)
		dr.errorConfig = dr.validateConfig(dr.config)
	})
	return dr.config, dr.errorConfig
}

func (dr *DependencyResolver) validateConfig(config *Config) error {
	if config.ReconcileRequeueSeconds <= 0 {
		return fmt.Errorf(lessOrEqualToZeroErrorMessage, "ReconcileRequeueSeconds")
	}
	return nil
}
