package depresolver

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
	"github.com/AbsaOSS/ohmyglb/pkg/controller/gslb/internal/utils"
	externaldns "github.com/kubernetes-incubator/external-dns/endpoint"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var defaultConfig = Config{30, "unset"}

func TestResolveSpecWithFilledFields(t *testing.T) {
	//arrange
	cl, gslb := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	err := resolver.ResolveGslbSpec(gslb)
	//assert
	assert.NoError(t, err)
	assert.Equal(t, 35, gslb.Spec.Strategy.DNSTtlSeconds)
	assert.Equal(t, 305, gslb.Spec.Strategy.SplitBrainThresholdSeconds)
}

func TestResolveSpecWithoutFields(t *testing.T) {
	//arrange
	cl, gslb := getTestContext("./testdata/free_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	err := resolver.ResolveGslbSpec(gslb)
	//assert
	assert.NoError(t, err)
	assert.Equal(t, predefinedStrategy.DNSTtlSeconds, gslb.Spec.Strategy.DNSTtlSeconds)
	assert.Equal(t, predefinedStrategy.SplitBrainThresholdSeconds, gslb.Spec.Strategy.SplitBrainThresholdSeconds)
}

func TestResolveSpecWithZeroSplitBrain(t *testing.T) {
	//arrange
	cl, gslb := getTestContext("./testdata/filled_omitempty_with_zero_splitbrain.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	err := resolver.ResolveGslbSpec(gslb)
	//assert
	assert.NoError(t, err)
	assert.Equal(t, 35, gslb.Spec.Strategy.DNSTtlSeconds)
	assert.Equal(t, predefinedStrategy.SplitBrainThresholdSeconds, gslb.Spec.Strategy.SplitBrainThresholdSeconds)
}

func TestResolveSpecWithEmptyFields(t *testing.T) {
	//arrange
	cl, gslb := getTestContext("./testdata/invalid_omitempty_empty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	err := resolver.ResolveGslbSpec(gslb)
	//assert
	assert.NoError(t, err)
	assert.Equal(t, predefinedStrategy.DNSTtlSeconds, gslb.Spec.Strategy.DNSTtlSeconds)
	assert.Equal(t, predefinedStrategy.SplitBrainThresholdSeconds, gslb.Spec.Strategy.SplitBrainThresholdSeconds)
}

func TestResolveSpecWithNegativeFields(t *testing.T) {
	//arrange
	cl, gslb := getTestContext("./testdata/invalid_omitempty_negative.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	err := resolver.ResolveGslbSpec(gslb)
	//assert
	assert.Error(t, err)
}

func TestSpecRunOnce(t *testing.T) {
	//arrange
	cl, gslb := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	err1 := resolver.ResolveGslbSpec(gslb)
	gslb.Spec.Strategy.DNSTtlSeconds = -100
	err2 := resolver.ResolveGslbSpec(gslb)
	//assert
	assert.NoError(t, err1)
	// err2 would not be empty
	assert.NoError(t, err2)
}

func TestResolveConfigWithOneValidEnv(t *testing.T) {
	//arrange
	defer cleanup()
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	expected := Config{50, "unset"}
	_ = os.Setenv(reconcileRequeueSecondsKey, strconv.Itoa(expected.ReconcileRequeueSeconds))
	//act
	config, err := resolver.ResolveOperatorConfig()
	//assert
	assert.NoError(t, err)
	assert.Equal(t, expected, *config)
}

func TestResolveConfigWithoutEnv(t *testing.T) {
	//arrange
	defer cleanup()
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	config, err := resolver.ResolveOperatorConfig()
	//assert
	assert.NoError(t, err)
	assert.Equal(t, defaultConfig, *config)
}

func TestResolveConfigWithZeroReconcileRequeueSecondsSync(t *testing.T) {
	//arrange
	defer cleanup()
	_ = os.Setenv(reconcileRequeueSecondsKey, "0")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	_, err := resolver.ResolveOperatorConfig()
	//assert
	assert.Error(t, err)
}

func TestResolveConfigWithTextReconcileRequeueSecondsSync(t *testing.T) {
	//arrange
	defer cleanup()
	_ = os.Setenv(reconcileRequeueSecondsKey, "invalid")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	config, err := resolver.ResolveOperatorConfig()
	//assert
	assert.NoError(t, err)
	assert.Equal(t, defaultConfig, *config)
}

func TestResolveConfigWithEmptyReconcileRequeueSecondsSync(t *testing.T) {
	//arrange
	defer cleanup()
	_ = os.Setenv(reconcileRequeueSecondsKey, "")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	config, err := resolver.ResolveOperatorConfig()
	//assert
	assert.NoError(t, err)
	assert.Equal(t, defaultConfig, *config)
}

func TestResolveConfigWithNegativeReconcileRequeueSecondsKey(t *testing.T) {
	//arrange
	defer cleanup()
	_ = os.Setenv(reconcileRequeueSecondsKey, "-1")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	_, err := resolver.ResolveOperatorConfig()
	//assert
	assert.Error(t, err)
}

func TestResolveConfigWithZeroReconcileRequeueSecondsKey(t *testing.T) {
	//arrange
	defer cleanup()
	_ = os.Setenv(reconcileRequeueSecondsKey, "0")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	_, err := resolver.ResolveOperatorConfig()
	//assert
	assert.Error(t, err)
}

func TestResolveConfigWithEmptyReconcileRequeueSecondsKey(t *testing.T) {
	//arrange
	defer cleanup()
	_ = os.Setenv(reconcileRequeueSecondsKey, "")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	config, err := resolver.ResolveOperatorConfig()
	//assert
	assert.NoError(t, err)
	assert.Equal(t, defaultConfig, *config)
}

func TestResolveConfigWithMalformedGeoTag(t *testing.T) {
	//arrange
	defer cleanup()
	_ = os.Setenv(clusterGeoTagKey, "i.am.wrong??.")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	_, err := resolver.ResolveOperatorConfig()
	//assert
	assert.Error(t, err)
}

func TestResolveConfigWithProperGeoTag(t *testing.T) {
	//arrange
	defer cleanup()
	_ = os.Setenv(clusterGeoTagKey, "eu-west-1")
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	_, err := resolver.ResolveOperatorConfig()
	//assert
	assert.NoError(t, err)
}

func TestConfigRunOnce(t *testing.T) {
	//arrange
	defer cleanup()
	cl, _ := getTestContext("./testdata/filled_omitempty.yaml")
	resolver := NewDependencyResolver(context.TODO(), cl)
	//act
	config1, err1 := resolver.ResolveOperatorConfig()
	_ = os.Setenv(reconcileRequeueSecondsKey, "100")
	//resolve again with new values
	config2, err2 := resolver.ResolveOperatorConfig()
	//assert
	assert.NoError(t, err1)
	assert.Equal(t, defaultConfig, *config1)
	//config2, err2 would be equal
	assert.NoError(t, err2)
	assert.Equal(t, *config1, *config2)
}

func cleanup() {
	if os.Unsetenv(reconcileRequeueSecondsKey) != nil {
		panic(fmt.Errorf("cleanup %s", reconcileRequeueSecondsKey))
	}
	if os.Unsetenv(clusterGeoTagKey) != nil {
		panic(fmt.Errorf("cleanup %s", clusterGeoTagKey))
	}
}

func getTestContext(testData string) (client.Client, *ohmyglbv1beta1.Gslb) {
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
	s.AddKnownTypes(ohmyglbv1beta1.SchemeGroupVersion, gslb)
	// Register external-dns DNSEndpoint CRD
	s.AddKnownTypes(schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}, &externaldns.DNSEndpoint{})
	cl := fake.NewFakeClientWithScheme(s, objs...)
	return cl, gslb
}
