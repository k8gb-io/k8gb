package dns

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/AbsaOSS/k8gb/controllers/depresolver"
	"github.com/AbsaOSS/k8gb/controllers/providers/assistant"

	ibclient "github.com/infobloxopen/infoblox-go-client"
	"github.com/stretchr/testify/assert"
)

var predefinedConfig = depresolver.Config{
	ReconcileRequeueSeconds: 30,
	ClusterGeoTag:           "us-west-1",
	ExtClustersGeoTags:      []string{"us-east-1"},
	EdgeDNSServer:           "8.8.8.8",
	EdgeDNSZone:             "example.com",
	DNSZone:                 "cloud.example.com",
	K8gbNamespace:           "k8gb",
	Infoblox: depresolver.Infoblox{
		Host:     "fakeinfoblox.example.com",
		Username: "foo",
		Password: "blah",
		Port:     443,
		Version:  "0.0.0",
	},
	Override: depresolver.Override{
		FakeInfobloxEnabled: true,
	},
}

func TestCanFilterOutDelegatedZoneEntryAccordingFQDNProvided(t *testing.T) {
	// arrange
	delegateTo := []ibclient.NameServer{
		{Address: "10.0.0.1", Name: "gslb-ns-cloud-example-com-eu.example.com"},
		{Address: "10.0.0.2", Name: "gslb-ns-cloud-example-com-eu.example.com"},
		{Address: "10.0.0.3", Name: "gslb-ns-cloud-example-com-eu.example.com"},
		{Address: "10.1.0.1", Name: "gslb-ns-cloud-example-com-za.example.com"},
		{Address: "10.1.0.2", Name: "gslb-ns-cloud-example-com-za.example.com"},
		{Address: "10.1.0.3", Name: "gslb-ns-cloud-example-com-za.example.com"},
	}
	want := []ibclient.NameServer{
		{Address: "10.0.0.1", Name: "gslb-ns-cloud-example-com-eu.example.com"},
		{Address: "10.0.0.2", Name: "gslb-ns-cloud-example-com-eu.example.com"},
		{Address: "10.0.0.3", Name: "gslb-ns-cloud-example-com-eu.example.com"},
	}
	customConfig := predefinedConfig
	customConfig.EdgeDNSZone = "example.com"
	customConfig.ExtClustersGeoTags = []string{"za"}
	a := assistant.NewGslbAssistant(nil, nil, customConfig.K8gbNamespace, customConfig.EdgeDNSServer)
	provider, err := NewInfobloxDNS(customConfig, a)
	require.NoError(t, err)
	// act
	extClusters := nsServerNameExt(customConfig)
	got := provider.filterOutDelegateTo(delegateTo, extClusters[0])
	// assert
	assert.Equal(t, want, got, "got:\n %q filtered out delegation records,\n\n want:\n %q", got, want)
}
