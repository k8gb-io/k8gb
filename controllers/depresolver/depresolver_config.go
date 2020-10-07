package depresolver

import (
	"fmt"
	"regexp"

	"github.com/AbsaOSS/k8gb/controllers/internal/env"
)

const (
	reconcileRequeueSecondsKey = "RECONCILE_REQUEUE_SECONDS"
	clusterGeoTagKey           = "CLUSTER_GEO_TAG"
	route53EnabledKey          = "ROUTE53_ENABLED"
)

// ResolveOperatorConfig executes once. It reads operator's configuration
// from environment variables into &Config and validates
func (dr *DependencyResolver) ResolveOperatorConfig() (*Config, error) {
	dr.onceConfig.Do(func() {
		dr.config = &Config{}
		// set predefined values when not read from the environment variables
		dr.config.ReconcileRequeueSeconds, _ = env.GetEnvAsIntOrFallback(reconcileRequeueSecondsKey, 30)
		dr.config.ClusterGeoTag = env.GetEnvAsStringOrFallback(clusterGeoTagKey, "unset")
		dr.config.Route53Enabled = env.GetEnvAsBoolOrFallback(route53EnabledKey, false)
		dr.errorConfig = dr.validateConfig(dr.config)
	})
	return dr.config, dr.errorConfig
}

func (dr *DependencyResolver) validateConfig(config *Config) error {
	if config.ReconcileRequeueSeconds <= 0 {
		return fmt.Errorf(lessOrEqualToZeroErrorMessage, "ReconcileRequeueSeconds")
	}
	geoTagRegexString := "^[a-zA-Z\\-\\d]*$"
	geoTagRegex, _ := regexp.Compile(geoTagRegexString)
	if !geoTagRegex.Match([]byte(config.ClusterGeoTag)) {
		return fmt.Errorf(doesNotMatchRegexMessage, "ClusterGeoTag", geoTagRegexString)
	}
	return nil
}
