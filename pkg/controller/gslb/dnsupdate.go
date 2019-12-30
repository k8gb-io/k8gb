package gslb

import (
	"context"

	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
	externaldns "github.com/kubernetes-incubator/external-dns/endpoint"
	"github.com/txn2/txeh"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileGslb) updateHostsConfigMap(gslb *ohmyglbv1beta1.Gslb, cmName string) error {
	workerIPs, err := r.getWorkerIPs()
	if err != nil {
		return err
	}

	managedHosts, err := r.getServiceHealthStatus(gslb)
	if err != nil {
		return err
	}

	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		return err
	}

	for _, ip := range workerIPs {
		for host, health := range managedHosts {
			if health == "Healthy" {
				hosts.AddHost(ip, host)
			}
		}
	}

	hfData := hosts.RenderHostsFile()

	nn := types.NamespacedName{
		Name:      cmName,
		Namespace: gslbOperatorNamespace,
	}

	cm := &corev1.ConfigMap{}

	err = r.client.Get(context.TODO(), nn, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Can't find CoreDNS configmap. Did you install helm chart?")
			return nil
		}
		return err
	}

	cm.Data["gslb.hosts"] = hfData

	err = r.client.Update(context.TODO(), cm)
	return err
}

func (r *ReconcileGslb) getWorkerIPs() ([]string, error) {
	workers, err := r.getHealthyWorkers()
	if err != nil {
		return nil, err
	}
	var IPs []string
	for _, address := range workers {
		IPs = append(IPs, address)
	}
	return IPs, err
}

func (r *ReconcileGslb) gslbDNSEndpoint(gslb *ohmyglbv1beta1.Gslb) (*externaldns.DNSEndpoint, error) {
	dnsEndpointSpec := externaldns.DNSEndpointSpec{}
	dnsEndpoint := &externaldns.DNSEndpoint{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gslb.Name,
			Namespace: gslb.Namespace,
		},
		Spec: dnsEndpointSpec,
	}

	err := controllerutil.SetControllerReference(gslb, dnsEndpoint, r.scheme)
	if err != nil {
		return nil, err
	}
	return dnsEndpoint, err
}

func (r *ReconcileGslb) ensureDNSEndpoint(request reconcile.Request,
	instance *ohmyglbv1beta1.Gslb,
	i *externaldns.DNSEndpoint,
) (*reconcile.Result, error) {
	found := &externaldns.DNSEndpoint{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      instance.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new DNSEndpoint", "DNSEndpoint.Namespace", i.Namespace, "DNSEndpoint.Name", i.Name)
		err = r.client.Create(context.TODO(), i)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new DNSEndpoint", "DNSEndpoint.Namespace", i.Namespace, "DNSEndpoint.Name", i.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get DNSEndpoint")
		return &reconcile.Result{}, err
	}

	return nil, nil
}
