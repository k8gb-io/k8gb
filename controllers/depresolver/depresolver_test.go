package depresolver

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	"github.com/AbsaOSS/k8gb/controllers/internal/utils"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

var predefinedConfig = Config{
	30,
	"us",
	[]string{"uk", "eu"},
	DNSTypeInfoblox,
	"cloud.example.com",
	"8.8.8.8",
	"example.com",
	false,
	Infoblox{
		"Infoblox.host.com",
		"0.0.3",
		443,
		"Infoblox",
		"secret",
	},
}

func TestResolveSpecWithFilledFields(t *testing.T) {
	// arrange
	cl, gslb := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	err := resolver.ResolveGslbSpec(gslb)
	// assert
	assert.NoError(t, err)
	assert.Equal(t, 35, gslb.Spec.Strategy.DNSTtlSeconds)
	assert.Equal(t, 305, gslb.Spec.Strategy.SplitBrainThresholdSeconds)
}

func TestResolveSpecWithoutFields(t *testing.T) {
	// arrange
	cl, gslb := getTestContext("./testdata/free_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	err := resolver.ResolveGslbSpec(gslb)
	// assert
	assert.NoError(t, err)
	assert.Equal(t, predefinedStrategy.DNSTtlSeconds, gslb.Spec.Strategy.DNSTtlSeconds)
	assert.Equal(t, predefinedStrategy.SplitBrainThresholdSeconds, gslb.Spec.Strategy.SplitBrainThresholdSeconds)
}

func TestResolveSpecWithZeroSplitBrain(t *testing.T) {
	// arrange
	cl, gslb := getTestContext("./testdata/filled_omitempty_with_zero_splitbrain.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	err := resolver.ResolveGslbSpec(gslb)
	// assert
	assert.NoError(t, err)
	assert.Equal(t, 35, gslb.Spec.Strategy.DNSTtlSeconds)
	assert.Equal(t, predefinedStrategy.SplitBrainThresholdSeconds, gslb.Spec.Strategy.SplitBrainThresholdSeconds)
}

func TestResolveSpecWithEmptyFields(t *testing.T) {
	// arrange
	cl, gslb := getTestContext("./testdata/invalid_omitempty_empty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	err := resolver.ResolveGslbSpec(gslb)
	// assert
	assert.NoError(t, err)
	assert.Equal(t, predefinedStrategy.DNSTtlSeconds, gslb.Spec.Strategy.DNSTtlSeconds)
	assert.Equal(t, predefinedStrategy.SplitBrainThresholdSeconds, gslb.Spec.Strategy.SplitBrainThresholdSeconds)
}

func TestResolveSpecWithNegativeFields(t *testing.T) {
	// arrange
	cl, gslb := getTestContext("./testdata/invalid_omitempty_negative.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	err := resolver.ResolveGslbSpec(gslb)
	// assert
	assert.Error(t, err)
}

func TestSpecRunOnce(t *testing.T) {
	// arrange
	cl, gslb := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	err1 := resolver.ResolveGslbSpec(gslb)
	gslb.Spec.Strategy.DNSTtlSeconds = -100
	err2 := resolver.ResolveGslbSpec(gslb)
	// assert
	assert.NoError(t, err1)
	// err2 would not be empty
	assert.NoError(t, err2)
}

func TestResolveConfigWithMultipleInvalidEnv(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSZone = ""
	expected.EdgeDNSServer = ""
	expected.ExtClustersGeoTags = []string{}
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestResolveConfigWithoutEnvVarsSet(t *testing.T) {
	// arrange
	defer cleanup()
	defaultConfig := Config{}
	defaultConfig.ReconcileRequeueSeconds = 30
	defaultConfig.EdgeDNSType = DNSTypeNoEdgeDNS
	defaultConfig.ExtClustersGeoTags = []string{}
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	config, err := resolver.ResolveOperatorConfig()
	// assert
	assert.Error(t, err)
	assert.Equal(t, defaultConfig, *config)
}

func TestResolveConfigWithReconcileRequeueSecondsSync(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.ReconcileRequeueSeconds = 3
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestResolveConfigWithTextReconcileRequeueSecondsSync(t *testing.T) {
	// arrange
	defer cleanup()
	configureEnvVar(predefinedConfig)
	_ = os.Setenv(ReconcileRequeueSecondsKey, "invalid")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	config, err := resolver.ResolveOperatorConfig()
	// assert
	assert.NoError(t, err)
	assert.Equal(t, predefinedConfig, *config)
}

func TestResolveConfigWithEmptyReconcileRequeueSecondsSync(t *testing.T) {
	// arrange
	defer cleanup()
	configureEnvVar(predefinedConfig)
	_ = os.Setenv(ReconcileRequeueSecondsKey, "")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	config, err := resolver.ResolveOperatorConfig()
	// assert
	assert.NoError(t, err)
	assert.Equal(t, predefinedConfig, *config)
}

func TestResolveConfigWithNegativeReconcileRequeueSecondsKey(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.ReconcileRequeueSeconds = -1
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestResolveConfigWithZeroReconcileRequeueSecondsKey(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.ReconcileRequeueSeconds = 0
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestResolveConfigWithEmptyReconcileRequeueSecondsKey(t *testing.T) {
	// arrange
	defer cleanup()
	configureEnvVar(predefinedConfig)
	_ = os.Setenv(ReconcileRequeueSecondsKey, "")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	config, err := resolver.ResolveOperatorConfig()
	// assert
	assert.NoError(t, err)
	assert.Equal(t, predefinedConfig, *config)
}

func TestResolveConfigWithoutReconcileRequeueSecondsKey(t *testing.T) {
	// arrange
	defer cleanup()
	// act,assert
	arrangeVariablesAndAssert(t, predefinedConfig, assert.NoError, ReconcileRequeueSecondsKey)
}

func TestResolveConfigWithMalformedGeoTag(t *testing.T) {
	// arrange
	for _, tag := range []string{"eu-west.1", "eu?", " ", "eu west1", "?/"} {
		defer cleanup()
		expected := predefinedConfig
		expected.ClusterGeoTag = tag
		// act,assert
		arrangeVariablesAndAssert(t, expected, assert.Error)
	}
}

func TestResolveConfigWithProperGeoTag(t *testing.T) {
	// arrange
	for _, tag := range []string{"eu-west-1", "eu-west1", "us", "1", "US"} {
		defer cleanup()
		expected := predefinedConfig
		expected.ClusterGeoTag = tag
		// act,assert
		arrangeVariablesAndAssert(t, expected, assert.NoError)
	}
}

func TestConfigRunOnce(t *testing.T) {
	// arrange
	defer cleanup()
	configureEnvVar(predefinedConfig)
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	config1, err1 := resolver.ResolveOperatorConfig()
	_ = os.Setenv(ReconcileRequeueSecondsKey, "100")
	// resolve again with new values
	config2, err2 := resolver.ResolveOperatorConfig()
	// assert
	assert.NoError(t, err1)
	assert.Equal(t, predefinedConfig, *config1)
	// config2, err2 would be equal
	assert.NoError(t, err2)
	assert.Equal(t, *config1, *config2)
}

func TestResolveConfigWithMalformedRoute53Enabled(t *testing.T) {
	// arrange
	defer cleanup()
	configureEnvVar(predefinedConfig)
	_ = os.Setenv(Route53EnabledKey, "i.am.wrong??.")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	config, err := resolver.ResolveOperatorConfig()
	// assert
	assert.NoError(t, err)
	assert.Equal(t, false, config.Route53Enabled)
}

func TestResolveConfigWithProperRoute53Enabled(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.Route53Enabled = true
	expected.Infoblox.Host = ""
	expected.EdgeDNSType = DNSTypeRoute53
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestResolveConfigWithoutRoute53(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.Route53Enabled = false
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError, Route53EnabledKey)
}

func TestResolveConfigWithEmptyRoute53(t *testing.T) {
	// arrange
	defer cleanup()
	configureEnvVar(predefinedConfig)
	_ = os.Setenv(Route53EnabledKey, "")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	config, err := resolver.ResolveOperatorConfig()
	// assert
	assert.NoError(t, err)
	assert.Equal(t, false, config.Route53Enabled)
}

func TestResolveConfigWithEmptyEdgeDnsServer(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSServer = ""
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestResolveConfigWithNoEdgeDnsServer(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSServer = ""
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error, EdgeDNSServerKey)
}

func TestResolveConfigWithEmptyIpAddressInEdgeDnsServer(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSServer = "22.147.90.2"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestResolveConfigWithHostnameEdgeDnsServer(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSServer = "server-nonprod.on.domain.l3.2l.com"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)

}

func TestResolveConfigWithInvalidHostnameEdgeDnsServer(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSServer = "https://server-nonprod.on.domain.l3.2l.com"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestResolveConfigWithInvalidIpAddressEdgeDnsServer(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSServer = "22.147.90.2."
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestResolveConfigWithEmptyEdgeDnsZone(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSZone = ""
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestResolveConfigWithoutEdgeDnsZone(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSZone = ""
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error, EdgeDNSZoneKey)
}

func TestResolveConfigWithHostnameEdgeDnsZone(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSZone = "company.2l.com"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestResolveConfigWithInvalidHostnameEdgeDnsZone(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSZone = "https://zone.com"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestResolveConfigWithValidHostnameDnsZone(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.DNSZone = "3l2.zo-ne.com"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestResolveConfigWithEmptyHostnameDnsZone(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.DNSZone = ""
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestResolveConfigWithInvalidHostnameDnsZone(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.DNSZone = "dns-zo?ne"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestResolveConfigWithoutDnsZone(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.DNSZone = ""
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error, DNSZoneKey)
}

func TestResolveEmptyExtGeoTags(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.DNSZone = ""
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error, DNSZoneKey)
}

func TestResolveOneExtGeoTags(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.ExtClustersGeoTags = []string{"foo"}
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestResolveMultipleExtGeoTags(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.ExtClustersGeoTags = []string{"foo", "blah", "boom"}
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestResolveUnsetExtGeoTags(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.ExtClustersGeoTags = []string{}
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error, ExtClustersGeoTagsKey)
}

func TestResolveInvalidExtGeoTags(t *testing.T) {
	// arrange
	for _, arr := range [][]string{{"good-tag", ".wrong.tag?"}, {"", ""}} {
		defer cleanup()
		expected := predefinedConfig
		expected.ExtClustersGeoTags = arr
		// act,assert
		arrangeVariablesAndAssert(t, expected, assert.Error)
	}
}

func TestResolveGeoTagExistsWithinExtGeoTags(t *testing.T) {
	// arrange
	defer cleanup()
	tag := "us-west1"
	expected := predefinedConfig
	expected.ClusterGeoTag = tag
	expected.ExtClustersGeoTags = []string{"us-east1", tag}
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestResolveGeoTagWithRepeatingExtGeoTags(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.ExtClustersGeoTags = []string{"foo", "blah", "foo"}
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestRoute53IsEnabledAndInfobloxIsConfigured(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.Route53Enabled = true
	expected.EdgeDNSType = DNSTypeRoute53 | DNSTypeInfoblox
	expected.Infoblox.Host = "Infoblox.domain"
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = "foo"
	expected.Infoblox.Password = "blah"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestRoute53IsDisabledAndInfobloxIsNotConfigured(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeNoEdgeDNS
	expected.Route53Enabled = false
	expected.Infoblox.Host = ""
	// act,assert
	// that's how our integration tests are running
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestRoute53IsDisabledButInfobloxIsConfigured(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.Route53Enabled = false
	expected.EdgeDNSType = DNSTypeInfoblox
	expected.Infoblox.Host = "Infoblox.domain"
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = "foo"
	expected.Infoblox.Password = "blah"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestRoute53IsEnabledButInfobloxIsNotConfigured(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.Route53Enabled = true
	expected.EdgeDNSType = DNSTypeRoute53
	expected.Infoblox.Host = ""
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = "foo"
	expected.Infoblox.Password = "blah"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestInfobloxGridHostIsEmpty(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeRoute53
	expected.Route53Enabled = true
	expected.Infoblox.Host = ""
	expected.Infoblox.Version = ""
	expected.Infoblox.Port = 0
	expected.Infoblox.Username = ""
	expected.Infoblox.Password = ""
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestInfobloxGridHostIsNotEmpty(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeInfoblox
	expected.Infoblox.Host = "test.domain"
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = "foo"
	expected.Infoblox.Password = "blah"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestInfobloxGridHostIsNotEmptyButInfobloxPropsAreEmpty(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeInfoblox
	expected.Infoblox.Host = "test.domain"
	expected.Infoblox.Version = ""
	expected.Infoblox.Port = 0
	expected.Infoblox.Username = ""
	expected.Infoblox.Password = ""
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestInfobloxGridHostIsEmptyButInfobloxPropsAreFilled(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeRoute53
	expected.Route53Enabled = true
	expected.Infoblox.Host = ""
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = "foo"
	expected.Infoblox.Password = "blah"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestInfobloxGridHostIsUnset(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeNoEdgeDNS
	expected.Route53Enabled = false
	expected.Infoblox.Host = ""
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = "foo"
	expected.Infoblox.Password = "blah"
	// act,assert
	// values are ignored and not validated
	arrangeVariablesAndAssert(t, expected, assert.NoError, InfobloxGridHostKey)
}

func TestInfobloxGridHostIsInvalid(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeInfoblox
	expected.Route53Enabled = false
	expected.Infoblox.Host = "dnfkjdnf kj"
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = "foo"
	expected.Infoblox.Password = "blah"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestInfobloxVersionIsValid(t *testing.T) {
	// arrange
	defer cleanup()
	// version can be empty!
	for _, v := range []string{"0.0.1", "v0.0.1", "v0.0.0-patch1", "2.3.5-patch1"} {
		expected := predefinedConfig
		expected.EdgeDNSType = DNSTypeInfoblox
		expected.Infoblox.Host = "test.domain"
		expected.Infoblox.Version = v
		expected.Infoblox.Port = 443
		expected.Infoblox.Username = "foo"
		expected.Infoblox.Password = "blah"
		// act,assert
		arrangeVariablesAndAssert(t, expected, assert.NoError)
	}
}

func TestInfobloxVersionIsInvalid(t *testing.T) {
	// arrange
	defer cleanup()
	for _, v := range []string{"0.1.*", "kkojo", "k12k", ""} {
		expected := predefinedConfig
		expected.EdgeDNSType = DNSTypeInfoblox
		expected.Infoblox.Host = "test.domain"
		expected.Infoblox.Version = v
		expected.Infoblox.Port = 443
		expected.Infoblox.Username = "foo"
		expected.Infoblox.Password = "blah"
		// act,assert
		arrangeVariablesAndAssert(t, expected, assert.Error)
	}
}

func TestInfobloxVersionIsUnset(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeInfoblox
	expected.Infoblox.Host = "test.domain"
	expected.Infoblox.Version = ""
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = "foo"
	expected.Infoblox.Password = "blah"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error, InfobloxVersionKey)
}

func TestInvalidInfobloxPort(t *testing.T) {
	// arrange
	defer cleanup()
	for _, p := range []int{-1, 0, 65536} {
		expected := predefinedConfig
		expected.EdgeDNSType = DNSTypeInfoblox
		expected.Infoblox.Host = "test.domain"
		expected.Infoblox.Version = "0.0.1"
		expected.Infoblox.Port = p
		expected.Infoblox.Username = "foo"
		expected.Infoblox.Password = "blah"
		// act,assert
		arrangeVariablesAndAssert(t, expected, assert.Error)
	}
}

func TestUnsetInfobloxPort(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeInfoblox
	expected.Infoblox.Host = "test.domain"
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 0
	expected.Infoblox.Username = "foo"
	expected.Infoblox.Password = "blah"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error, InfobloxPortKey)
}

func TestValidInfobloxUserPasswordAndPort(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeInfoblox
	expected.Infoblox.Host = "test.domain"
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = "infobloxUser"
	expected.Infoblox.Password = "blah"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.NoError)
}

func TestEmptyInfobloxUser(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeInfoblox
	expected.Infoblox.Host = "test.domain"
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = ""
	expected.Infoblox.Password = "blah"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestUnsetInfobloxUser(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeInfoblox
	expected.Infoblox.Host = "test.domain"
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = ""
	expected.Infoblox.Password = "blah"
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error, InfobloxUsernameKey)
}

func TestEmptyInfobloxPassword(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeInfoblox
	expected.Infoblox.Host = "test.domain"
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = "infobloxUser"
	expected.Infoblox.Password = ""
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error)
}

func TestUnsetInfobloxPassword(t *testing.T) {
	// arrange
	defer cleanup()
	expected := predefinedConfig
	expected.EdgeDNSType = DNSTypeInfoblox
	expected.Infoblox.Host = "test.domain"
	expected.Infoblox.Version = "0.0.1"
	expected.Infoblox.Port = 443
	expected.Infoblox.Username = "foo"
	expected.Infoblox.Password = ""
	// act,assert
	arrangeVariablesAndAssert(t, expected, assert.Error, InfobloxPasswordKey)
}

// arrangeVariablesAndAssert sets string environment variables and asserts `expected` argument with ResolveOperatorConfig() output. The last parameter unsets the values
func arrangeVariablesAndAssert(t *testing.T, expected Config, errf func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool, unset ...string) {
	configureEnvVar(expected)
	for _, v := range unset {
		_ = os.Unsetenv(v)
	}
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	// act
	config, err := resolver.ResolveOperatorConfig()
	// assert
	assert.Equal(t, expected, *config)
	errf(t, err)
}

func cleanup() {
	for _, s := range []string{ReconcileRequeueSecondsKey, ClusterGeoTagKey, ExtClustersGeoTagsKey, EdgeDNSZoneKey, DNSZoneKey, EdgeDNSServerKey,
		Route53EnabledKey, InfobloxGridHostKey, InfobloxVersionKey, InfobloxPortKey, InfobloxUsernameKey, InfobloxPasswordKey} {
		if os.Unsetenv(s) != nil {
			panic(fmt.Errorf("cleanup %s", s))
		}
	}
}

func configureEnvVar(config Config) {
	_ = os.Setenv(ReconcileRequeueSecondsKey, strconv.Itoa(config.ReconcileRequeueSeconds))
	_ = os.Setenv(ClusterGeoTagKey, config.ClusterGeoTag)
	_ = os.Setenv(ExtClustersGeoTagsKey, strings.Join(config.ExtClustersGeoTags, ","))
	_ = os.Setenv(EdgeDNSServerKey, config.EdgeDNSServer)
	_ = os.Setenv(EdgeDNSZoneKey, config.EdgeDNSZone)
	_ = os.Setenv(DNSZoneKey, config.DNSZone)
	_ = os.Setenv(Route53EnabledKey, strconv.FormatBool(config.Route53Enabled))
	_ = os.Setenv(InfobloxGridHostKey, config.Infoblox.Host)
	_ = os.Setenv(InfobloxVersionKey, config.Infoblox.Version)
	_ = os.Setenv(InfobloxPortKey, strconv.Itoa(config.Infoblox.Port))
	_ = os.Setenv(InfobloxUsernameKey, config.Infoblox.Username)
	_ = os.Setenv(InfobloxPasswordKey, config.Infoblox.Password)
}

func getTestContext(testData string) (client.Client, *k8gbv1beta1.Gslb) {
	// Create a fake client to mock API calls.
	var gslbYaml, err = ioutil.ReadFile(testData)
	if err != nil {
		panic(fmt.Errorf("can't open example CR file: %s", testData))
	}
	gslb, err := utils.YamlToGslb(gslbYaml)
	if err != nil {
		panic(err)
	}
	objs := []runtime.Object{
		gslb,
	}
	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(k8gbv1beta1.GroupVersion, gslb)
	// Register external-dns DNSEndpoint CRD
	s.AddKnownTypes(schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}, &externaldns.DNSEndpoint{})
	cl := fake.NewFakeClientWithScheme(s, objs...)
	return cl, gslb
}
