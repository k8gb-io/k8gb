package gslb

import (
	"context"

	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileGslb) gslbIngress(gslb *ohmyglbv1beta1.Gslb) (*v1beta1.Ingress, error) {
	ingress := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gslb.Name,
			Namespace: gslb.Namespace,
		},
		Spec: gslb.Spec.Ingress,
	}

	err := controllerutil.SetControllerReference(gslb, ingress, r.scheme)
	if err != nil {
		return nil, err
	}
	return ingress, err
}

func (r *ReconcileGslb) ensureIngress(request reconcile.Request,
	instance *ohmyglbv1beta1.Gslb,
	i *v1beta1.Ingress,
) (*reconcile.Result, error) {
	found := &v1beta1.Ingress{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      instance.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new Ingress", "Ingress.Namespace", i.Namespace, "Ingress.Name", i.Name)
		err = r.client.Create(context.TODO(), i)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Ingress", "Ingress.Namespace", i.Namespace, "Ingress.Name", i.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Ingress")
		return &reconcile.Result{}, err
	}

	// Update existing object with new spec
	found.Spec = i.Spec
	r.client.Update(context.TODO(), found)

	return nil, nil
}
