package depresolver

import (
	"context"
	"errors"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
)

var predefinedStrategy = k8gbv1beta1.Strategy{
	DNSTtlSeconds:              30,
	SplitBrainThresholdSeconds: 300,
}

// ResolveGslbSpec executes once during reconciliation. At first cycle it reads
// omitempty properties and attach predefined values in case they are not defined.
// ResolveGslbSpec returns error if any input is invalid
func (dr *DependencyResolver) ResolveGslbSpec(ctx context.Context, gslb *k8gbv1beta1.Gslb) error {
	dr.onceSpec.Do(func() {
		strategy := &gslb.Spec.Strategy
		// set predefined values if missing in the yaml
		if strategy.DNSTtlSeconds == 0 {
			strategy.DNSTtlSeconds = predefinedStrategy.DNSTtlSeconds
		}
		if strategy.SplitBrainThresholdSeconds == 0 {
			strategy.SplitBrainThresholdSeconds = predefinedStrategy.SplitBrainThresholdSeconds
		}
		dr.errorSpec = dr.validateStrategy(strategy)
		if dr.errorSpec == nil {
			dr.errorSpec = dr.client.Update(ctx, gslb)
		}
		for _, rule := range gslb.Spec.Ingress.Rules {
			if rule.HTTP == nil {
				dr.errorSpec = errors.New("missing .spec.ingress.rules.http")
			}
		}
	})
	return dr.errorSpec
}

func (dr *DependencyResolver) validateStrategy(strategy *k8gbv1beta1.Strategy) (err error) {
	err = field("DNSTtlSeconds", strategy.DNSTtlSeconds).isHigherOrEqualToZero().err
	if err != nil {
		return
	}
	err = field("SplitBrainThresholdSeconds", strategy.SplitBrainThresholdSeconds).isHigherOrEqualToZero().err
	if err != nil {
		return
	}
	return
}
