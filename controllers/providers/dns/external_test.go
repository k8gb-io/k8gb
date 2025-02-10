package dns

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	utils2 "github.com/k8gb-io/k8gb/controllers/utils"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	"github.com/k8gb-io/k8gb/controllers/depresolver"
	"github.com/k8gb-io/k8gb/controllers/mocks"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	externaldns "sigs.k8s.io/external-dns/endpoint"

	"github.com/k8gb-io/k8gb/controllers/providers/assistant"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// test data
var targetIPs = []string{
	"10.0.1.38",
	"10.0.1.40",
	"10.0.1.39",
}
var a = struct {
	Config              depresolver.Config
	Gslb                *k8gbv1beta1.Gslb
	TargetIPs           []string
	TargetNSNamesSorted []string
}{
	Config: depresolver.Config{
		ReconcileRequeueSeconds: 30,
		NSRecordTTL:             30,
		ClusterGeoTag:           "us",
		EdgeDNSServers: []utils2.DNSServer{
			{
				Host: "dns.cloud.example.com",
				Port: 53,
			},
		},
		DelegationZones: depresolver.DelegationZones{
			{
				Domain:            "cloud.example.com",
				Zone:              "example.com",
				ClusterNSName:     "gslb-ns-us-cloud.example.com",
				ExtClusterNSNames: map[string]string{"eu": "gslb-ns-eu-cloud.example.com", "za": "gslb-ns-za-cloud.example.com"},
			},
		},
		K8gbNamespace: "k8gb",
	},
	Gslb: func() *k8gbv1beta1.Gslb {
		var crSampleYaml = "../../../deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr_roundrobin_ingress.yaml"
		gslbYaml, _ := os.ReadFile(crSampleYaml)
		gslb, _ := utils2.YamlToGslb(gslbYaml)
		gslb.Status.LoadBalancer.ExposedIPs = targetIPs
		return gslb
	}(),
	TargetIPs: targetIPs,
	TargetNSNamesSorted: []string{
		"gslb-ns-eu-cloud.example.com",
		"gslb-ns-us-cloud.example.com",
		"gslb-ns-za-cloud.example.com",
	},
}

var expectedDNSEndpoint = &externaldns.DNSEndpoint{
	ObjectMeta: metav1.ObjectMeta{
		Name:        fmt.Sprintf("k8gb-ns-%s", externalDNSTypeCommon),
		Namespace:   a.Config.K8gbNamespace,
		Annotations: map[string]string{"k8gb.absa.oss/dnstype": string(externalDNSTypeCommon)},
	},
	Spec: externaldns.DNSEndpointSpec{
		Endpoints: []*externaldns.Endpoint{
			{
				DNSName:    a.Config.DelegationZones[0].Domain,
				RecordTTL:  externaldns.TTL(a.Config.NSRecordTTL),
				RecordType: "NS",
				Targets:    a.TargetNSNamesSorted,
			},
			{
				DNSName:    "gslb-ns-us-cloud.example.com",
				RecordTTL:  externaldns.TTL(a.Config.NSRecordTTL),
				RecordType: "A",
				Targets:    a.TargetIPs,
			},
		},
	},
}

func TestCreateZoneDelegationOnExternalDNS(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockAssistant(ctrl)
	ep1 := expectedDNSEndpoint.DeepCopy()
	ep1.Name = "k8gb-ns-extdns-cloud-example-com"
	p := NewExternalDNS(a.Config, m)
	m.EXPECT().SaveDNSEndpoint(a.Config.K8gbNamespace, gomock.Eq(ep1)).Return(nil).Times(1).
		Do(func(ns string, ep *externaldns.DNSEndpoint) {
			require.True(t, reflect.DeepEqual(ep, ep1))
			require.Equal(t, ns, a.Config.K8gbNamespace)
		})
	gslb := a.Gslb.DeepCopy()
	gslb.Status.Servers = []*k8gbv1beta1.Server{
		{
			Host: "cloud.example.com",
		},
	}
	// act
	err := p.CreateZoneDelegationForExternalDNS(gslb)
	// assert
	assert.NoError(t, err)
}

func TestCreateZoneDelegationOnExternalDNSWithMultipleEndpoints(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockAssistant(ctrl)
	di := depresolver.DelegationZoneInfo{
		Domain:            "common.sampledomain.com",
		Zone:              "sampledomain.com",
		ClusterNSName:     "gslb-ns-us-common.sampledomain.com",
		ExtClusterNSNames: map[string]string{"za": "gslb-ns-za-common.sampledomain.com", "eu": "gslb-ns-eu-common.sampledomain.com"},
	}
	a.Config.DelegationZones = append(a.Config.DelegationZones, di)
	ep1 := expectedDNSEndpoint.DeepCopy()
	ep1.Name = "k8gb-ns-extdns-cloud-example-com"
	ep2 := expectedDNSEndpoint.DeepCopy()
	ep2.Name = "k8gb-ns-extdns-common-sampledomain-com"
	ep2.Spec.Endpoints[0].DNSName = "common.sampledomain.com"
	ep2.Spec.Endpoints[1].DNSName = "gslb-ns-us-common.sampledomain.com"
	ep2.Spec.Endpoints[0].Targets = []string{
		"gslb-ns-eu-common.sampledomain.com",
		"gslb-ns-us-common.sampledomain.com",
		"gslb-ns-za-common.sampledomain.com",
	}
	p := NewExternalDNS(a.Config, m)
	m.EXPECT().SaveDNSEndpoint(a.Config.K8gbNamespace, gomock.Eq(ep1)).Return(nil).Times(1).
		Do(func(ns string, ep *externaldns.DNSEndpoint) {
			require.True(t, reflect.DeepEqual(ep, ep1))
			require.Equal(t, ns, a.Config.K8gbNamespace)
		})

	m.EXPECT().SaveDNSEndpoint(a.Config.K8gbNamespace, gomock.Eq(ep2)).Return(nil).Times(1).
		Do(func(ns string, ep *externaldns.DNSEndpoint) {
			require.True(t, reflect.DeepEqual(ep, ep2))
			require.Equal(t, ns, a.Config.K8gbNamespace)
		})
	gslb1 := a.Gslb.DeepCopy()
	gslb1.Status.Servers = []*k8gbv1beta1.Server{
		{
			Host: "cloud.example.com",
		},
	}
	gslb2 := a.Gslb.DeepCopy()
	gslb2.Status.Servers = []*k8gbv1beta1.Server{
		{
			Host: "common.sampledomain.com",
		},
	}
	gslb3 := a.Gslb.DeepCopy()
	gslb3.Status.Servers = []*k8gbv1beta1.Server{
		{
			Host: "common.dummy.com",
		},
	}
	// act
	err := p.CreateZoneDelegationForExternalDNS(gslb1)
	assert.NoError(t, err)
	err = p.CreateZoneDelegationForExternalDNS(gslb2)
	assert.NoError(t, err)
	err = p.CreateZoneDelegationForExternalDNS(gslb3)
	assert.Error(t, err)
}

func TestSaveNewDNSEndpointOnExternalDNS(t *testing.T) {
	// arrange
	var ep = &corev1.Endpoints{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Endpoints",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "k8gb-ns-extdns",
			Namespace: "test-gslb",
		},
	}
	endpointToSave := expectedDNSEndpoint
	endpointToSave.Namespace = a.Gslb.Namespace

	runtimeScheme := runtime.NewScheme()
	schemeBuilder := &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}}
	schemeBuilder.Register(&externaldns.DNSEndpoint{}, &externaldns.DNSEndpointList{})
	require.NoError(t, corev1.AddToScheme(runtimeScheme))
	require.NoError(t, k8gbv1beta1.AddToScheme(runtimeScheme))
	require.NoError(t, schemeBuilder.AddToScheme(runtimeScheme))

	var cl = fake.NewClientBuilder().WithScheme(runtimeScheme).WithObjects(ep).Build()

	assistant := assistant.NewGslbAssistant(cl, a.Config.K8gbNamespace, a.Config)
	p := NewExternalDNS(a.Config, assistant)
	// act, assert
	err := p.SaveDNSEndpoint(a.Gslb, expectedDNSEndpoint)
	assert.NoError(t, err)
}

func TestSaveExistingDNSEndpointOnExternalDNS(t *testing.T) {
	// arrange
	endpointToSave := expectedDNSEndpoint
	endpointToSave.Namespace = a.Gslb.Namespace

	runtimeScheme := runtime.NewScheme()
	schemeBuilder := &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}}
	schemeBuilder.Register(&externaldns.DNSEndpoint{}, &externaldns.DNSEndpointList{})
	require.NoError(t, corev1.AddToScheme(runtimeScheme))
	require.NoError(t, k8gbv1beta1.AddToScheme(runtimeScheme))
	require.NoError(t, schemeBuilder.AddToScheme(runtimeScheme))

	var cl = fake.NewClientBuilder().WithScheme(runtimeScheme).WithObjects(endpointToSave).Build()
	assistant := assistant.NewGslbAssistant(cl, a.Config.K8gbNamespace, a.Config)
	p := NewExternalDNS(a.Config, assistant)
	// act, assert
	err := p.SaveDNSEndpoint(a.Gslb, endpointToSave)
	assert.NoError(t, err)
}
