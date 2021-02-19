package controllers

import (
	"context"
	"reflect"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *GslbReconciler) gslbIngress(gslb *k8gbv1beta1.Gslb) (*v1beta1.Ingress, error) {
	metav1.SetMetaDataAnnotation(&gslb.ObjectMeta, strategyAnnotation, gslb.Spec.Strategy.Type)
	if gslb.Spec.Strategy.PrimaryGeoTag != "" {
		metav1.SetMetaDataAnnotation(&gslb.ObjectMeta, primaryGeoTagAnnotation, gslb.Spec.Strategy.PrimaryGeoTag)
	}
	ingress := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        gslb.Name,
			Namespace:   gslb.Namespace,
			Annotations: gslb.Annotations,
		},
		Spec: gslb.Spec.Ingress,
	}

	err := controllerutil.SetControllerReference(gslb, ingress, r.Scheme)
	if err != nil {
		return nil, err
	}
	return ingress, err
}

func (r *GslbReconciler) saveIngress(instance *k8gbv1beta1.Gslb, i *v1beta1.Ingress) error {
	found := &v1beta1.Ingress{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      instance.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new Ingress", "Ingress.Namespace", i.Namespace, "Ingress.Name", i.Name)
		err = r.Create(context.TODO(), i)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Ingress", "Ingress.Namespace", i.Namespace, "Ingress.Name", i.Name)
			return err
		}
		// Creation was successful
		return nil
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Ingress")
		return err
	}

	// Update existing object with new spec and annotations
	if !ingressCompare(found, i) {
		found.Spec = i.Spec
		found.Annotations = i.Annotations
		err = r.Update(context.TODO(), found)
		if errors.IsConflict(err) {
			r.Log.Info("CONFLICT, Ingress updated by third-party")
			return nil
		}
		if err != nil {
			// Update failed
			log.Error(err, "Failed to update Ingress", "Ingress.Namespace", found.Namespace, "Ingress.Name", found.Name)
			return err
		}
	}

	return nil
}

func ingressCompare(found *v1beta1.Ingress, i *v1beta1.Ingress) bool {
	for k, v := range i.Annotations {
		if found.Annotations[k] != v {
			return false
		}
	}
	return reflect.DeepEqual(found.Spec, i.Spec)
}
