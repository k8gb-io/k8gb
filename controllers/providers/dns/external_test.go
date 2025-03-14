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
	"context"
	"os"
	"testing"

	"github.com/k8gb-io/k8gb/controllers/depresolver"
	"github.com/k8gb-io/k8gb/controllers/logging"
	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

func TestCreateZoneDelegation(t *testing.T) {
	// arrange
	getFakeClient := func(ctx context.Context, namespace string, names ...string) client.Client {
		cl := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
		for _, v := range names {
			_ = cl.Get(ctx, client.ObjectKey{Name: v, Namespace: namespace}, &externaldns.DNSEndpoint{})
			_ = cl.Get(ctx, client.ObjectKey{Name: v, Namespace: namespace}, &externaldns.DNSEndpoint{})
			_ = cl.Create(ctx, &externaldns.DNSEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      v,
					Namespace: namespace},
			})
		}
		return cl
	}

	ctx := context.TODO()
	var tests = []struct {
		name          string
		config        depresolver.Config
		expectedError bool
		client        client.Client
	}{
		{
			name: "new cloud.example.com",
			config: depresolver.Config{
				K8gbNamespace: "k8gb",
				NSRecordTTL:   60,
				DelegationZones: []*depresolver.DelegationZoneInfo{
					{
						Domain:        "cloud.example.com",
						Zone:          "example.com",
						NegativeTTL:   60,
						IPs:           []string{"10.0.0.1", "10.0.0.2"},
						ClusterNSName: "gslb-ns-eu-k8gb-test-gslb.cloud.example.com",
						ExtClusterNSNames: map[string]string{
							"us": "gslb-ns-us-k8gb-test-gslb.cloud.example.com",
						},
					},
				},
			},
			client:        getFakeClient(ctx, "k8gb", "k8gb-ns-extdns-cloud-example-com"),
			expectedError: false,
		},
		{
			name: "new cloud.example.com, cloud.example.org",
			config: depresolver.Config{
				K8gbNamespace: "k8gb",
				NSRecordTTL:   60,
				DelegationZones: []*depresolver.DelegationZoneInfo{
					{
						Domain:        "cloud.example.com",
						Zone:          "example.com",
						NegativeTTL:   60,
						IPs:           []string{"10.0.0.1", "10.0.0.2"},
						ClusterNSName: "gslb-ns-eu-k8gb-test-gslb.cloud.example.com",
						ExtClusterNSNames: map[string]string{
							"us": "gslb-ns-us-k8gb-test-gslb.cloud.example.com",
						},
					},
					{
						Domain:        "cloud.example.org",
						Zone:          "example.org",
						NegativeTTL:   60,
						IPs:           []string{"10.0.0.1", "10.0.0.2"},
						ClusterNSName: "gslb-ns-eu-k8gb-test-gslb.cloud.example.org",
						ExtClusterNSNames: map[string]string{
							"us": "gslb-ns-us-k8gb-test-gslb.cloud.example.org",
						},
					},
				},
			},
			client:        getFakeClient(ctx, "k8gb", "k8gb-ns-extdns-cloud-example-com", "k8gb-ns-extdns-cloud-example-org"),
			expectedError: false,
		},
	}

	// act
	// assert
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			externalDNSProvider := NewExternalDNS(ctx, test.client, test.config, logging.TestLogger())
			for _, zone := range test.config.DelegationZones {
				err := externalDNSProvider.CreateZoneDelegation(zone)
				assert.Equal(t, test.expectedError, err != nil)
			}
		})
	}
}

func TestMain(m *testing.M) {
	scheme.Scheme.AddKnownTypes(schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}, &externaldns.DNSEndpoint{})
	os.Exit(m.Run())
}
