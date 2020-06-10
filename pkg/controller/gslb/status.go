package gslb

import (
	"context"
	"regexp"

	kgbv1beta1 "github.com/AbsaOSS/kgb/pkg/apis/kgb/v1beta1"
	externaldns "github.com/kubernetes-incubator/external-dns/endpoint"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	types "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ReconcileGslb) updateGslbStatus(gslb *kgbv1beta1.Gslb) error {
	var err error

	gslb.Status.ServiceHealth, err = r.getServiceHealthStatus(gslb)
	if err != nil {
		return err
	}

	err = r.updateIngressHostsPerStatusMetric(gslb, gslb.Status.ServiceHealth)
	if err != nil {
		return err
	}

	gslb.Status.HealthyRecords, err = r.getHealthyRecords(gslb)
	if err != nil {
		return err
	}

	gslb.Status.GeoTag = r.config.ClusterGeoTag

	err = r.updateHealthyRecordsMetric(gslb, gslb.Status.HealthyRecords)
	if err != nil {
		return err
	}

	err = r.client.Status().Update(context.TODO(), gslb)
	return err
}

func (r *ReconcileGslb) getServiceHealthStatus(gslb *kgbv1beta1.Gslb) (map[string]string, error) {
	serviceHealth := make(map[string]string)
	for _, rule := range gslb.Spec.Ingress.Rules {
		for _, path := range rule.HTTP.Paths {
			service := &corev1.Service{}
			finder := client.ObjectKey{
				Namespace: gslb.Namespace,
				Name:      path.Backend.ServiceName,
			}
			err := r.client.Get(context.TODO(), finder, service)
			if err != nil {
				if errors.IsNotFound(err) {
					serviceHealth[rule.Host] = "NotFound"
					continue
				}
				return serviceHealth, err
			}

			endpoints := &corev1.Endpoints{}

			nn := types.NamespacedName{
				Name:      path.Backend.ServiceName,
				Namespace: gslb.Namespace,
			}

			err = r.client.Get(context.TODO(), nn, endpoints)
			if err != nil {
				return serviceHealth, err
			}

			serviceHealth[rule.Host] = "Unhealthy"
			if len(endpoints.Subsets) > 0 {
				for _, subset := range endpoints.Subsets {
					if len(subset.Addresses) > 0 {
						serviceHealth[rule.Host] = "Healthy"
					}
				}
			}
		}
	}
	return serviceHealth, nil
}

func (r *ReconcileGslb) getHealthyRecords(gslb *kgbv1beta1.Gslb) (map[string][]string, error) {

	dnsEndpoint := &externaldns.DNSEndpoint{}

	nn := types.NamespacedName{
		Name:      gslb.Name,
		Namespace: gslb.Namespace,
	}

	err := r.client.Get(context.TODO(), nn, dnsEndpoint)
	if err != nil {
		return nil, err
	}

	healthyRecords := make(map[string][]string)

	serviceRegex, _ := regexp.Compile("^localtargets")
	for _, endpoint := range dnsEndpoint.Spec.Endpoints {
		local := serviceRegex.Match([]byte(endpoint.DNSName))
		if !local && endpoint.RecordType == "A" {
			if len(endpoint.Targets) > 0 {
				healthyRecords[endpoint.DNSName] = endpoint.Targets
			}
		}
	}

	return healthyRecords, nil
}
