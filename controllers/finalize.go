package controllers

import (
	"context"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
)

func (r *GslbReconciler) finalizeGslb(gslb *k8gbv1beta1.Gslb) (err error) {
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.
	err = r.DNSProvider.Finalize(gslb)
	if err != nil {
		log.Error(err, "Can't finalize GSLB (%s)")
		return
	}
	log.Info("Successfully finalized Gslb")
	return
}

func (r *GslbReconciler) addFinalizer(gslb *k8gbv1beta1.Gslb) error {
	log.Info("Adding Finalizer for the Gslb")
	gslb.SetFinalizers(append(gslb.GetFinalizers(), gslbFinalizer))

	// Update CR
	err := r.Update(context.TODO(), gslb)
	if err != nil {
		log.Error(err, "Failed to update Gslb with finalizer")
		return err
	}
	return nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func remove(list []string, s string) []string {
	for i, v := range list {
		if v == s {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}
