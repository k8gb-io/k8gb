package controllers

import (
	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	k8gbv1beta1io "github.com/k8gb-io/k8gb/api/v1beta1io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func convertGslbLegacyToIO(gslb *k8gbv1beta1.Gslb) *k8gbv1beta1io.Gslb {
	if gslb == nil {
		return &k8gbv1beta1io.Gslb{}
	}

	return &k8gbv1beta1io.Gslb{
		TypeMeta: metav1.TypeMeta{
			APIVersion: k8gbv1beta1io.GroupVersion.String(),
			Kind:       "Gslb",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      gslb.Name,
			Namespace: gslb.Namespace,
		},
		Spec: convertSpecLegacyToIO(gslb.Spec),
	}
}

func convertSpecLegacyToIO(spec k8gbv1beta1.GslbSpec) k8gbv1beta1io.GslbSpec {
	return k8gbv1beta1io.GslbSpec{
		Ingress:     k8gbv1beta1io.FromV1IngressSpec(k8gbv1beta1.ToV1IngressSpec(spec.Ingress)),
		Strategy:    convertStrategyLegacyToIO(spec.Strategy),
		ResourceRef: convertResourceRefLegacyToIO(spec.ResourceRef),
	}
}

func convertStrategyLegacyToIO(strategy k8gbv1beta1.Strategy) k8gbv1beta1io.Strategy {
	return k8gbv1beta1io.Strategy{
		Type:                       strategy.Type,
		Weight:                     strategy.Weight,
		PrimaryGeoTag:              strategy.PrimaryGeoTag,
		DNSTtlSeconds:              strategy.DNSTtlSeconds,
		SplitBrainThresholdSeconds: strategy.SplitBrainThresholdSeconds,
	}
}

func convertResourceRefLegacyToIO(ref k8gbv1beta1.ResourceRef) k8gbv1beta1io.ResourceRef {
	return k8gbv1beta1io.ResourceRef{
		ObjectReference: ref.ObjectReference,
		LabelSelector:   ref.LabelSelector,
	}
}
