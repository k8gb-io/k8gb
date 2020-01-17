package gslb

import (
	"context"

	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
	externaldns "github.com/kubernetes-incubator/external-dns/endpoint"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	types "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ReconcileGslb) updateGslbStatus(gslb *ohmyglbv1beta1.Gslb) error {
	gslb.Status.ManagedHosts = getGslbManagedHosts(gslb)

	var err error

	gslb.Status.ServiceHealth, err = r.getServiceHealthStatus(gslb)
	if err != nil {
		return err
	}

	gslb.Status.HealthyRecords, err = r.getHealthyRecords(gslb)
	if err != nil {
		return err
	}

	err = r.client.Status().Update(context.TODO(), gslb)
	return err
}

func getGslbManagedHosts(gslb *ohmyglbv1beta1.Gslb) []string {
	var hosts []string
	for _, rule := range gslb.Spec.Ingress.Rules {
		hosts = append(hosts, rule.Host)
	}
	return hosts
}

func (r *ReconcileGslb) getServiceHealthStatus(gslb *ohmyglbv1beta1.Gslb) (map[string]string, error) {
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

func (r *ReconcileGslb) getHealthyHosts(gslb *ohmyglbv1beta1.Gslb) ([]string, error) {
	serviceHealth, err := r.getServiceHealthStatus(gslb)
	if err != nil {
		return nil, err
	}
	var healthyHosts []string
	for host, status := range serviceHealth {
		if status == "Healthy" {
			healthyHosts = append(healthyHosts, host)
		}
	}

	return healthyHosts, err
}

func (r *ReconcileGslb) getHealthyRecords(gslb *ohmyglbv1beta1.Gslb) (map[string][]string, error) {
	healthyHosts, err := r.getHealthyHosts(gslb)
	if err != nil {
		return nil, err
	}

	dnsEndpoint := &externaldns.DNSEndpoint{}

	nn := types.NamespacedName{
		Name:      gslb.Name,
		Namespace: gslb.Namespace,
	}

	err = r.client.Get(context.TODO(), nn, dnsEndpoint)
	if err != nil {
		return nil, err
	}

	healthyRecords := make(map[string][]string)

	for _, host := range healthyHosts {
		for _, endpoint := range dnsEndpoint.Spec.Endpoints {
			if endpoint.DNSName == host {
				healthyRecords[host] = endpoint.Targets
			}
		}
	}

	return healthyRecords, nil
}
