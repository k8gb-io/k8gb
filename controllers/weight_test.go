package controllers

import (
	"context"
	"fmt"
	"testing"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	"github.com/k8gb-io/k8gb/controllers/depresolver"
	"github.com/k8gb-io/k8gb/controllers/mocks"
	"github.com/k8gb-io/k8gb/controllers/providers/assistant"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

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

func TestWeight(t *testing.T) {
	// arrange
	type wrr struct {
		weight  string
		targets []string
	}
	var tests = []struct {
		name           string
		data           map[string]wrr
		injectWeights  bool
		expectedLabels map[string]string
	}{
		{
			name:          "eu35-us50-za15",
			injectWeights: true,
			data: map[string]wrr{
				"eu": {weight: "35%", targets: []string{"10.10.0.1", "10.10.0.2"}},
				"us": {weight: "50%", targets: []string{"10.0.0.1", "10.0.0.2"}},
				"za": {weight: "15%", targets: []string{"10.22.0.1", "10.22.0.2", "10.22.1.1"}},
			},
			expectedLabels: map[string]string{
				"strategy":       depresolver.RoundRobinStrategy,
				"weight-us-0-50": "10.0.0.1",
				"weight-us-1-50": "10.0.0.2",
				"weight-eu-0-35": "10.10.0.1",
				"weight-eu-1-35": "10.10.0.2",
				"weight-za-0-15": "10.22.0.1",
				"weight-za-1-15": "10.22.0.2",
				"weight-za-2-15": "10.22.1.1",
			},
		},
		{
			name:          "eu100-us0-za0",
			injectWeights: true,
			data: map[string]wrr{
				"eu": {weight: "100%", targets: []string{"10.10.0.1", "10.10.0.2"}},
				"us": {weight: "0%", targets: []string{"10.0.0.1", "10.0.0.2"}},
				"za": {weight: "0%", targets: []string{"10.22.0.1", "10.22.0.2", "10.22.1.1"}},
			},
			expectedLabels: map[string]string{
				"strategy":        depresolver.RoundRobinStrategy,
				"weight-us-0-0":   "10.0.0.1",
				"weight-us-1-0":   "10.0.0.2",
				"weight-eu-0-100": "10.10.0.1",
				"weight-eu-1-100": "10.10.0.2",
				"weight-za-0-0":   "10.22.0.1",
				"weight-za-1-0":   "10.22.0.2",
				"weight-za-2-0":   "10.22.1.1",
			},
		},
		{
			name:          "weights-without-external-targets",
			injectWeights: true,
			data: map[string]wrr{
				"eu": {weight: "25%", targets: []string{}},
				"us": {weight: "75%", targets: []string{}},
				"za": {weight: "0%", targets: []string{}},
			},
			expectedLabels: map[string]string{
				"strategy": depresolver.RoundRobinStrategy,
			},
		},
		{
			name:          "weights-without-external-targets",
			injectWeights: true,
			data: map[string]wrr{
				"eu": {weight: "25%", targets: []string{}},
				"us": {weight: "75%", targets: []string{}},
				"za": {weight: "0%", targets: []string{}},
			},
			expectedLabels: map[string]string{
				"strategy": depresolver.RoundRobinStrategy,
			},
		},
		{
			name:          "no weights without external targets",
			injectWeights: false,
			data:          map[string]wrr{},
			expectedLabels: map[string]string{
				"strategy": depresolver.RoundRobinStrategy,
			},
		},
		{
			name:          "no weights with external targets",
			injectWeights: false,
			data: map[string]wrr{
				"eu": {weight: "100%", targets: []string{"10.10.0.1", "10.10.0.2"}},
				"us": {weight: "0%", targets: []string{"10.0.0.1", "10.0.0.2"}},
				"za": {weight: "0%", targets: []string{"10.22.0.1", "10.22.0.2", "10.22.1.1"}},
			},
			expectedLabels: map[string]string{
				"strategy": depresolver.RoundRobinStrategy,
			},
		},
		{
			name:          "empty weights",
			injectWeights: true,
			data:          map[string]wrr{},
			expectedLabels: map[string]string{
				"strategy": depresolver.RoundRobinStrategy,
			},
		},
	}

	// act
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			injectWeight := func(ctx context.Context, gslb *k8gbv1beta1.Gslb, client client.Client) error {
				if !test.injectWeights {
					return nil
				}
				gslb.Spec.Strategy.Weight = make(map[string]k8gbv1beta1.Percentage, 0)
				for k, w := range test.data {
					gslb.Spec.Strategy.Weight[k] = k8gbv1beta1.Percentage(w.weight)
				}
				return nil
			}

			assertAnnotation := func(gslb *k8gbv1beta1.Gslb, ep *externaldns.DNSEndpoint) error {
				require.NotNil(t, ep)
				require.NotNil(t, gslb)
				// annotation is equal to tested value
				for _, e := range ep.Spec.Endpoints {
					for k, v := range e.Labels {
						assert.Equal(t, test.expectedLabels[k], v)
					}
					assert.Equal(t, len(test.expectedLabels), len(e.Labels))
				}
				return nil
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// settings := provideSettings(t, predefinedConfig)
			m := mocks.NewMockProvider(ctrl)
			r := mocks.NewMockGslbResolver(ctrl)
			m.EXPECT().GslbIngressExposedIPs(gomock.Any()).Return([]string{}, nil).Times(1)
			m.EXPECT().SaveDNSEndpoint(gomock.Any(), gomock.Any()).Do(assertAnnotation).Return(fmt.Errorf("save DNS error")).Times(1)
			m.EXPECT().CreateZoneDelegationForExternalDNS(gomock.Any()).Return(nil).AnyTimes()
			r.EXPECT().ResolveGslbSpec(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(injectWeight).AnyTimes()

			ts := assistant.Targets{}
			for k, w := range test.data {
				ts[k] = &assistant.Target{IPs: w.targets}
			}
			m.EXPECT().GetExternalTargets("roundrobin.cloud.example.com").Return(ts).Times(1)
			m.EXPECT().GetExternalTargets("notfound.cloud.example.com").Return(assistant.Targets{}).Times(1)
			m.EXPECT().GetExternalTargets("unhealthy.cloud.example.com").Return(assistant.Targets{}).Times(1)

			settings := provideSettings(t, predefinedConfig)
			settings.reconciler.DNSProvider = m
			settings.reconciler.DepResolver = r

			// act, assert
			_, _ = settings.reconciler.Reconcile(context.TODO(), settings.request)
		})
	}
}
