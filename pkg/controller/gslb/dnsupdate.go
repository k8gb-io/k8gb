package gslb

import (
	"context"

	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
	"github.com/txn2/txeh"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	types "k8s.io/apimachinery/pkg/types"
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
		Namespace: gslb.Namespace,
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
