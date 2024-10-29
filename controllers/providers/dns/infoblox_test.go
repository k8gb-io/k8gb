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
	"os"
	"testing"

	"github.com/k8gb-io/k8gb/controllers/utils"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	"github.com/k8gb-io/k8gb/controllers/depresolver"
	"github.com/k8gb-io/k8gb/controllers/mocks"
	"github.com/k8gb-io/k8gb/controllers/providers/assistant"

	ibclient "github.com/infobloxopen/infoblox-go-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	ref       = "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:cloud.example.com/default"
	namespace = "test-gslb"
	gslbName  = "test-gslb"
)

var (
	defaultConfig = depresolver.Config{
		ReconcileRequeueSeconds: 30,
		ClusterGeoTag:           "us-west-1",
		ExtClustersGeoTags:      []string{"us-east-1"},
		EdgeDNSServers: []utils.DNSServer{
			{
				Host: "8.8.8.8",
				Port: 53,
			},
		},
		EdgeDNSZone:   "example.com",
		DNSZone:       "cloud.example.com",
		K8gbNamespace: "k8gb",
		Infoblox: depresolver.Infoblox{
			Host:     "fakeinfoblox.example.com",
			Username: "foo",
			Password: "blah",
			Port:     443,
			Version:  "0.0.0",
		},
	}

	defaultDelegatedZone = ibclient.ZoneDelegated{
		Fqdn:       defaultConfig.DNSZone,
		DelegateTo: []ibclient.NameServer{},
		Ref:        ref,
	}

	defaultGslb = new(k8gbv1beta1.Gslb)
)

func TestCanFilterOutDelegatedZoneEntryAccordingFQDNProvided(t *testing.T) {
	// arrange
	delegateTo := []ibclient.NameServer{
		{Address: "10.0.0.1", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.0.0.2", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.0.0.3", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.1.0.1", Name: "gslb-ns-za-cloud.example.com"},
		{Address: "10.1.0.2", Name: "gslb-ns-za-cloud.example.com"},
		{Address: "10.1.0.3", Name: "gslb-ns-za-cloud.example.com"},
	}
	want := []ibclient.NameServer{
		{Address: "10.0.0.1", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.0.0.2", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.0.0.3", Name: "gslb-ns-eu-cloud.example.com"},
	}
	customConfig := defaultConfig
	customConfig.EdgeDNSZone = "example.com"
	customConfig.ExtClustersGeoTags = []string{"za"}
	a := assistant.NewGslbAssistant(nil, customConfig.K8gbNamespace, customConfig.EdgeDNSServers)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockInfobloxClient(ctrl)
	provider := NewInfobloxDNS(customConfig, a, m)
	// act
	extClusters := customConfig.GetExternalClusterNSNames()
	got := provider.filterOutDelegateTo(delegateTo, extClusters["za"])
	// assert
	assert.Equal(t, want, got, "got:\n %q filtered out delegation records,\n\n want:\n %q", got, want)
}
func TestCanSanitizeDelegatedZone(t *testing.T) {
	// arrange
	local := []ibclient.NameServer{
		{Address: "10.0.0.3", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.0.0.1", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.0.0.2", Name: "gslb-ns-eu-cloud.example.com"},
	}
	upstream := []ibclient.NameServer{
		{Address: "10.0.0.3", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.1.0.3", Name: "gslb-ns-za-cloud.example.com"},
		{Address: "10.0.0.1", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.1.0.2", Name: "gslb-ns-za-cloud.example.com"},
		{Address: "10.0.0.2", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.1.0.1", Name: "gslb-ns-za-cloud.example.com"},
	}
	want := []ibclient.NameServer{
		{Address: "10.0.0.1", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.0.0.2", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.0.0.3", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.1.0.1", Name: "gslb-ns-za-cloud.example.com"},
		{Address: "10.1.0.2", Name: "gslb-ns-za-cloud.example.com"},
		{Address: "10.1.0.3", Name: "gslb-ns-za-cloud.example.com"},
	}
	customConfig := defaultConfig
	customConfig.EdgeDNSZone = "example.com"
	customConfig.ExtClustersGeoTags = []string{"za"}
	customConfig.ClusterGeoTag = "eu"
	a := assistant.NewGslbAssistant(nil, customConfig.K8gbNamespace, customConfig.EdgeDNSServers)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockInfobloxClient(ctrl)
	provider := NewInfobloxDNS(customConfig, a, m)
	// act
	got := provider.sanitizeDelegateZone(local, upstream)
	// assert
	assert.Equal(t, want, got, "got:\n %q filtered out delegation records,\n\n want:\n %q", got, want)
}

func TestSortNameServer(t *testing.T) {
	delegateTo := []ibclient.NameServer{
		{Address: "10.0.0.3", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.1.0.3", Name: "gslb-ns-za-cloud.example.com"},
		{Address: "10.0.0.1", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.1.0.2", Name: "gslb-ns-za-cloud.example.com"},
		{Address: "10.0.0.2", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.1.0.1", Name: "gslb-ns-za-cloud.example.com"},
	}
	want := []ibclient.NameServer{
		{Address: "10.0.0.1", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.0.0.2", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.0.0.3", Name: "gslb-ns-eu-cloud.example.com"},
		{Address: "10.1.0.1", Name: "gslb-ns-za-cloud.example.com"},
		{Address: "10.1.0.2", Name: "gslb-ns-za-cloud.example.com"},
		{Address: "10.1.0.3", Name: "gslb-ns-za-cloud.example.com"},
	}
	sortZones(delegateTo)
	assert.Equal(t, want, delegateTo, "got:\n %q \n\n want:\n %q", delegateTo, want)
}

func TestInfobloxCreateZoneDelegationForExternalDNS(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := mocks.NewMockAssistant(ctrl)
	cl := mocks.NewMockInfobloxClient(ctrl)
	con := mocks.NewMockIBConnector(ctrl)
	con.EXPECT().CreateObject(gomock.Any()).Return(ref, nil).AnyTimes()
	con.EXPECT().UpdateObject(gomock.Any(), gomock.Any()).Return(ref, nil).Times(1)
	con.EXPECT().GetObject(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(2, []ibclient.ZoneDelegated{defaultDelegatedZone}).Return(nil)
	cl.EXPECT().GetObjectManager().Return(ibclient.NewObjectManager(con, "k8gbclient", ""), nil).Times(1)
	config := defaultConfig
	provider := NewInfobloxDNS(config, a, cl)

	// act
	err := provider.CreateZoneDelegationForExternalDNS(defaultGslb)
	// assert
	assert.NoError(t, err)
}

func TestInfobloxCreateZoneDelegationForExternalDNSWithSplitBrainEnabled(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := mocks.NewMockAssistant(ctrl)
	cl := mocks.NewMockInfobloxClient(ctrl)
	con := mocks.NewMockIBConnector(ctrl)
	a.EXPECT().InspectTXTThreshold(gomock.Any(), gomock.Any()).Do(func(fqdn string, _ interface{}) {
		require.Equal(t, "test-gslb-heartbeat-us-east-1.example.com", fqdn)
	}).Return(nil).Times(1)
	con.EXPECT().CreateObject(gomock.Any()).Return(ref, nil).AnyTimes()
	con.EXPECT().UpdateObject(gomock.Any(), gomock.Any()).Return(ref, nil).Times(2)
	con.EXPECT().GetObject(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(2, []ibclient.ZoneDelegated{defaultDelegatedZone}).Return(nil)
	cl.EXPECT().GetObjectManager().Return(ibclient.NewObjectManager(con, "k8gbclient", ""), nil).Times(1)
	con.EXPECT().GetObject(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(2, []ibclient.RecordTXT{{Ref: ref}}).
		Return(nil).Do(func(arg0 *ibclient.RecordTXT, _, _ interface{}) {
		require.Equal(t, "test-gslb-heartbeat-us-west-1.example.com", arg0.Name)
	}).AnyTimes()
	config := defaultConfig
	config.SplitBrainCheck = true
	provider := NewInfobloxDNS(config, a, cl)

	// act
	err := provider.CreateZoneDelegationForExternalDNS(defaultGslb)
	// assert
	assert.NoError(t, err)
}

func TestInfobloxCreateZoneDelegationForExternalDNSWithSplitBrainEnabledCreatingNewHeartBeatRecord(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := mocks.NewMockAssistant(ctrl)
	cl := mocks.NewMockInfobloxClient(ctrl)
	con := mocks.NewMockIBConnector(ctrl)
	a.EXPECT().InspectTXTThreshold(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	con.EXPECT().CreateObject(gomock.Any()).Return(ref, nil).AnyTimes()
	con.EXPECT().UpdateObject(gomock.Any(), gomock.Any()).Return(ref, nil).Times(1)
	con.EXPECT().GetObject(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(2, []ibclient.ZoneDelegated{defaultDelegatedZone}).Return(nil)
	cl.EXPECT().GetObjectManager().Return(ibclient.NewObjectManager(con, "k8gbclient", ""), nil).Times(1)
	con.EXPECT().GetObject(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(2, []ibclient.RecordTXT{}).Return(nil).AnyTimes()
	config := defaultConfig
	config.SplitBrainCheck = true
	provider := NewInfobloxDNS(config, a, cl)

	// act
	err := provider.CreateZoneDelegationForExternalDNS(defaultGslb)
	// assert
	assert.NoError(t, err)
}

func TestInfobloxFinalize(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := mocks.NewMockAssistant(ctrl)
	cl := mocks.NewMockInfobloxClient(ctrl)
	con := mocks.NewMockIBConnector(ctrl)
	con.EXPECT().DeleteObject(gomock.Any()).Return(ref, nil).Do(func(arg0 string) {
		require.Equal(t, arg0, ref)
	}).AnyTimes()
	con.EXPECT().GetObject(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(2, []ibclient.ZoneDelegated{defaultDelegatedZone}).
		Return(nil).Times(1)
	con.EXPECT().GetObject(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(2, []ibclient.RecordTXT{{Ref: ref}}).
		Return(nil).Do(func(arg0 *ibclient.RecordTXT, _, _ interface{}) {
		require.Equal(t, "test-gslb-heartbeat-us-west-1.example.com", arg0.Name)
	}).Times(1)
	cl.EXPECT().GetObjectManager().Return(ibclient.NewObjectManager(con, "k8gbclient", ""), nil).Times(1)
	config := defaultConfig
	provider := NewInfobloxDNS(config, a, cl)

	// act
	err := provider.Finalize(defaultGslb, nil)

	// assert
	assert.NoError(t, err)
}

func TestEmptySort(t *testing.T) {
	// arrange
	delegateTo := make([]ibclient.NameServer, 0)
	// act
	sortZones(delegateTo)
	// assert
	assert.Equal(t, 0, len(delegateTo))
}

func TestNilSort(t *testing.T) {
	// arrange
	delegateTo := []ibclient.NameServer(nil)
	// act
	sortZones(delegateTo)
	// assert
	assert.Nil(t, delegateTo)
}

func TestMain(m *testing.M) {
	defaultGslb.Name = gslbName
	defaultGslb.Namespace = namespace
	os.Exit(m.Run())
}
