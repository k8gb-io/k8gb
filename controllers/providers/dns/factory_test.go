package dns

/*
Copyright 2021-2025 The k8gb Contributors.

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
	"context"
	"testing"

	"github.com/k8gb-io/k8gb/controllers/depresolver"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestFactory(t *testing.T) {
	cl := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
	var tests = []struct {
		name          string
		config        depresolver.Config
		expectedType  string
		expectedError bool
		client        client.Client
	}{
		{
			name: "ExternalDNS",
			config: depresolver.Config{
				EdgeDNSType: depresolver.DNSTypeExternal,
			},
			expectedType: "EXTDNS",
			client:       cl,
		},
		{
			name: "Infoblox",
			config: depresolver.Config{
				EdgeDNSType: depresolver.DNSTypeInfoblox,
			},
			expectedType: "Infoblox",
			client:       cl,
		},
		{
			name: "No EdgeDNS",
			config: depresolver.Config{
				EdgeDNSType: depresolver.DNSTypeNoEdgeDNS,
			},
			expectedType: "EMPTY",
			client:       cl,
		},
		{
			name: "Invalid",
			config: depresolver.Config{
				EdgeDNSType: depresolver.DNSTypeNoEdgeDNS,
			},
			expectedType:  "EMPTY",
			expectedError: true,
			client:        nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := NewDNSProviderFactory(context.Background(), test.client, test.config)
			if test.expectedError {
				assert.Error(t, err)
				return
			}
			p := f.Provider()
			assert.Equal(t, test.expectedType, p.String())
		})
	}
}
