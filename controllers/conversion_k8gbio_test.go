package controllers

import (
	"testing"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	k8gbv1beta1io "github.com/k8gb-io/k8gb/api/v1beta1io"
	"github.com/stretchr/testify/require"
)

func TestConvertLegacyToIOCopiesSpec(t *testing.T) {
	legacy := &k8gbv1beta1.Gslb{
		Spec: k8gbv1beta1.GslbSpec{
			Strategy: k8gbv1beta1.Strategy{Type: "roundRobin", DNSTtlSeconds: 30},
		},
	}

	io := convertGslbLegacyToIO(legacy)

	require.Equal(t, k8gbv1beta1io.GroupVersion.String(), io.APIVersion)
	require.Equal(t, "Gslb", io.Kind)
	require.Equal(t, legacy.Spec.Strategy.Type, io.Spec.Strategy.Type)
	require.Equal(t, legacy.Spec.Strategy.DNSTtlSeconds, io.Spec.Strategy.DNSTtlSeconds)
}
