package gslb

import (
	"context"

	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
)

func (r *ReconcileGslb) updateManagedHostsStatus(gslb *ohmyglbv1beta1.Gslb) error {
	var hosts []string
	for _, rule := range gslb.Spec.Ingress.Rules {
		hosts = append(hosts, rule.Host)
	}
	gslb.Status.ManagedHosts = hosts
	err := r.client.Status().Update(context.TODO(), gslb)
	return err
}
