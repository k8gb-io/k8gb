package dns

import (
	"io/ioutil"
	"testing"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	"github.com/AbsaOSS/k8gb/controllers/depresolver"
	"github.com/AbsaOSS/k8gb/controllers/internal/utils"

	"github.com/stretchr/testify/assert"
)

var commonConfig = depresolver.Config{
	ClusterGeoTag:      "us",
	DNSZone:            "example.com",
	ExtClustersGeoTags: []string{"uk", "eu"},
	EdgeDNSZone:        "8.8.8.8",
}

func TestNsServerName(t *testing.T) {
	// arrange
	// act
	result := nsServerName(commonConfig)
	// assert
	assert.Equal(t, "gslb-ns-example-com-us.8.8.8.8", result)
}

func TestEmptyClusterGeoTagNSServerName(t *testing.T) {
	// arrange
	config := commonConfig
	config.ClusterGeoTag = ""
	// act
	result := nsServerName(config)
	// assert
	assert.Equal(t, "gslb-ns-example-com-.8.8.8.8", result)
}

func TestNsServerNameExt(t *testing.T) {
	// arrange
	expected := []string{"gslb-ns-example-com-uk.8.8.8.8", "gslb-ns-example-com-eu.8.8.8.8"}
	// act
	result := nsServerNameExt(commonConfig)
	// assert
	assert.Equal(t, expected, result)
}

func TestNsServerNameExtWithEmptyGeoTag(t *testing.T) {
	// arrange
	config := commonConfig
	config.ExtClustersGeoTags = []string{}
	// act
	result := nsServerNameExt(config)
	// assert
	assert.Equal(t, []string{}, result)
}

func TestGeneratesProperExternalNSTargetFQDNsAccordingToTheGeoTags(t *testing.T) {
	// arrange
	want := []string{"gslb-ns-cloud-example-com-za.example.com"}
	customConfig := predefinedConfig
	customConfig.EdgeDNSZone = "example.com"
	customConfig.ExtClustersGeoTags = []string{"za"}
	// act
	got := nsServerNameExt(customConfig)
	// assert
	assert.Equal(t, want, got, "got:\n %q externalGslb NS records,\n\n want:\n %q", got, want)
}

func TestCanGenerateExternalHeartbeatFQDNs(t *testing.T) {
	// arrange
	want := []string{"test-gslb-heartbeat-za.example.com"}
	customConfig := predefinedConfig
	customConfig.EdgeDNSZone = "example.com"
	customConfig.ExtClustersGeoTags = []string{"za"}
	gslb := getGSLB(t)
	// act
	got := getExternalClusterHeartbeatFQDNs(gslb, customConfig)
	// assert
	assert.Equal(t, want, got, "got:\n %s unexpected heartbeat records,\n\n want:\n %s", got, want)
}

func getGSLB(t *testing.T) *k8gbv1beta1.Gslb {
	var crSampleYaml = "../../../deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr.yaml"
	gslbYaml, err := ioutil.ReadFile(crSampleYaml)
	if err != nil {
		t.Fatalf("Can't open example CR file: %s", crSampleYaml)
	}
	gslb, _ := utils.YamlToGslb(gslbYaml)
	return gslb
}
