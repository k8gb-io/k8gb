package test

import "github.com/AbsaOSS/gopkg/env"

// you can split into more structs if you like
type testSettings struct {
	DNSZone         string
	PrimaryGeoTag   string
	SecondaryGeoTag string
	DNSServer1      string
	Port1           int
	DNSServer2      string
	Port2           int
	Cluster1        string
	Cluster2        string
	PodinfoImage    string
}

var (
	settings testSettings
)

func init() {
	p1, _ := env.GetEnvAsIntOrFallback("DNS_SERVER1_PORT", 5053)
	p2, _ := env.GetEnvAsIntOrFallback("DNS_SERVER2_PORT", 5054)
	settings = testSettings{
		DNSZone:         env.GetEnvAsStringOrFallback("GSLB_DOMAIN", "cloud.example.com"),
		PrimaryGeoTag:   env.GetEnvAsStringOrFallback("PRIMARY_GEO_TAG", "eu"),
		SecondaryGeoTag: env.GetEnvAsStringOrFallback("SECONDARY_GEO_TAG", "us"),
		DNSServer1:      env.GetEnvAsStringOrFallback("DNS_SERVER1", "localhost"),
		Port1:           p1,
		DNSServer2:      env.GetEnvAsStringOrFallback("DNS_SERVER2", "localhost"),
		Port2:           p2,
		Cluster1:        env.GetEnvAsStringOrFallback("K8GB_CLUSTER1", "k3d-test-gslb1"),
		Cluster2:        env.GetEnvAsStringOrFallback("K8GB_CLUSTER2", "k3d-test-gslb2"),
		PodinfoImage:    env.GetEnvAsStringOrFallback("PODINFO_IMAGE_REPO", "ghcr.io/stefanprodan/podinfo"),
	}
}
