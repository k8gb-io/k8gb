/*
Copyright 2021 The k8gb Contributors.

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
package controllers

import (
	"fmt"
	"sort"
	"strings"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

func sortTargets(targets []string) []string {
	sort.Slice(targets, func(i, j int) bool {
		return targets[i] < targets[j]
	})
	return targets
}

func (r *GslbReconciler) gslbDNSEndpoint(gslb *k8gbv1beta1.Gslb) (*externaldns.DNSEndpoint, error) {
	var gslbHosts []*externaldns.Endpoint
	var ttl = externaldns.TTL(gslb.Spec.Strategy.DNSTtlSeconds)

	serviceHealth, err := r.getServiceHealthStatus(gslb)
	if err != nil {
		return nil, err
	}

	localTargets, err := r.DNSProvider.GslbIngressExposedIPs(gslb)
	if err != nil {
		return nil, err
	}

	for host, health := range serviceHealth {
		var finalTargets []string

		if !strings.Contains(host, r.Config.EdgeDNSZone) {
			return nil, fmt.Errorf("ingress host %s does not match delegated zone %s", host, r.Config.EdgeDNSZone)
		}

		isPrimary := gslb.Spec.Strategy.PrimaryGeoTag == r.Config.ClusterGeoTag
		isHealthy := health == k8gbv1beta1.Healthy

		if isHealthy {
			finalTargets = append(finalTargets, localTargets...)
			localTargetsHost := fmt.Sprintf("localtargets-%s", host)
			dnsRecord := &externaldns.Endpoint{
				DNSName:    localTargetsHost,
				RecordTTL:  ttl,
				RecordType: "A",
				Targets:    localTargets,
			}
			gslbHosts = append(gslbHosts, dnsRecord)
		}

		// Check if host is alive on external Gslb
		externalTargets := r.DNSProvider.GetExternalTargets(host)

		sortTargets(externalTargets)

		if len(externalTargets) > 0 {
			switch gslb.Spec.Strategy.Type {
			case roundRobinStrategy, geoStrategy:
				finalTargets = append(finalTargets, externalTargets...)
			case failoverStrategy:
				// If cluster is Primary
				if isPrimary {
					// If cluster is Primary and Healthy return only own targets
					// If cluster is Primary and Unhealthy return Secondary external targets
					if !isHealthy {
						finalTargets = externalTargets
						log.Info().
							Str("gslb", gslb.Name).
							Str("cluster", gslb.Spec.Strategy.PrimaryGeoTag).
							Str("targets", fmt.Sprintf("%v", finalTargets)).
							Str("workload", k8gbv1beta1.Unhealthy.String()).
							Msg("Executing failover strategy for primary cluster")
					}
				} else {
					// If cluster is Secondary and Primary external cluster is Healthy
					// then return Primary external targets.
					// Return own targets by default.
					finalTargets = externalTargets
					log.Info().
						Str("gslb", gslb.Name).
						Str("cluster", gslb.Spec.Strategy.PrimaryGeoTag).
						Str("targets", fmt.Sprintf("%v", finalTargets)).
						Str("workload", k8gbv1beta1.Healthy.String()).
						Msg("Executing failover strategy for secondary cluster")
				}
			}
		} else {
			log.Info().
				Str("host", host).
				Msg("No external targets have been found for host")
		}

		r.updateRuntimeStatus(gslb, isPrimary, health, finalTargets)
		log.Info().
			Str("gslb", gslb.Name).
			Str("targets", fmt.Sprintf("%v", finalTargets)).
			Msg("Final target list")

		if len(finalTargets) > 0 {
			dnsRecord := &externaldns.Endpoint{
				DNSName:    host,
				RecordTTL:  ttl,
				RecordType: "A",
				Targets:    finalTargets,
				Labels: externaldns.Labels{
					"strategy": gslb.Spec.Strategy.Type,
				},
			}
			gslbHosts = append(gslbHosts, dnsRecord)
		}
	}
	dnsEndpointSpec := externaldns.DNSEndpointSpec{
		Endpoints: gslbHosts,
	}

	dnsEndpoint := &externaldns.DNSEndpoint{
		ObjectMeta: metav1.ObjectMeta{
			Name:        gslb.Name,
			Namespace:   gslb.Namespace,
			Annotations: map[string]string{"k8gb.absa.oss/dnstype": "local"},
			Labels:      map[string]string{"k8gb.absa.oss/dnstype": "local"},
		},
		Spec: dnsEndpointSpec,
	}

	err = controllerutil.SetControllerReference(gslb, dnsEndpoint, r.Scheme)
	if err != nil {
		return nil, err
	}
	return dnsEndpoint, err
}

func (r *GslbReconciler) updateRuntimeStatus(gslb *k8gbv1beta1.Gslb, isPrimary bool, isHealthy k8gbv1beta1.HealthStatus, finalTargets []string) {
	switch gslb.Spec.Strategy.Type {
	case roundRobinStrategy:
		m.UpdateRoundrobinStatus(gslb, isHealthy, finalTargets)
	case geoStrategy:
		m.UpdateGeoIPStatus(gslb, isHealthy, finalTargets)
	case failoverStrategy:
		m.UpdateFailoverStatus(gslb, isPrimary, isHealthy, finalTargets)
	}
}
