/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	externaldns "github.com/kubernetes-incubator/external-dns/endpoint"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
)

// GslbReconciler reconciles a Gslb object
type GslbReconciler struct {
	// New generated fields
	//client.Client
	//Log    logr.Logger
	//Scheme *runtime.Scheme
	// Olde migrated fields
	client      client.Client
	scheme      *runtime.Scheme
	config      *depresolver.Config
	depResolver *depresolver.DependencyResolver
}

// +kubebuilder:rbac:groups=k8gb.absa.oss,resources=gslbs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=k8gb.absa.oss,resources=gslbs/status,verbs=get;update;patch

func (r *GslbReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("gslb", req.NamespacedName)

	// Fetch the Gslb instance
	gslb := &k8gbv1beta1.Gslb{}
	err := r.client.Get(context.TODO(), request.NamespacedName, gslb)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	var result *reconcile.Result

	err = r.depResolver.ResolveGslbSpec(gslb)
	if err != nil {
		log.Error(err, "resolving spec.strategy")
		return reconcile.Result{}, err
	}
	// == Finalizer business ==

	// Check if the Gslb instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isGslbMarkedToBeDeleted := gslb.GetDeletionTimestamp() != nil
	if isGslbMarkedToBeDeleted {
		if contains(gslb.GetFinalizers(), gslbFinalizer) {
			// Run finalization logic for gslbFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeGslb(gslb); err != nil {
				return reconcile.Result{}, err
			}

			// Remove gslbFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			gslb.SetFinalizers(remove(gslb.GetFinalizers(), gslbFinalizer))
			err := r.client.Update(context.TODO(), gslb)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	// Add finalizer for this CR
	if !contains(gslb.GetFinalizers(), gslbFinalizer) {
		if err := r.addFinalizer(gslb); err != nil {
			return reconcile.Result{}, err
		}
	}

	// == Ingress ==========
	ingress, err := r.gslbIngress(gslb)
	if err != nil {
		// Requeue the request
		return reconcile.Result{}, err
	}

	result, err = r.ensureIngress(
		request,
		gslb,
		ingress)
	if result != nil {
		return *result, err
	}

	// == external-dns dnsendpoints CRs ==
	dnsEndpoint, err := r.gslbDNSEndpoint(gslb)
	if err != nil {
		// Requeue the request
		return reconcile.Result{}, err
	}

	result, err = r.ensureDNSEndpoint(
		request,
		gslb,
		dnsEndpoint)
	if result != nil {
		return *result, err
	}

	// == handle delegated zone in Edge DNS

	result, err = r.configureZoneDelegation(gslb)
	if result != nil {
		return *result, err
	}

	// == Status =
	err = r.updateGslbStatus(gslb)
	if err != nil {
		// Requeue the request
		return reconcile.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Second * time.Duration(r.config.ReconcileRequeueSeconds)}, nil
}

func (r *GslbReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8gbv1beta1.Gslb{}).
		Owns(&v1beta1.Ingress{}).
		Owns(&externaldns.DNSEndpoint{}).
		Complete(r)

	// Figure out Gslb resource name to Reconcile when non controlled Endpoint is updated
	mapFn := handler.ToRequestsFunc(
		func(a handler.MapObject) []reconcile.Request {
			gslbList := &k8gbv1beta1.GslbList{}
			opts := []client.ListOption{
				client.InNamespace(a.Meta.GetNamespace()),
			}
			c := mgr.GetClient()
			err := c.List(context.TODO(), gslbList, opts...)
			if err != nil {
				log.Info("Can't fetch gslb objects")
				return nil
			}
			gslbName := ""
			for _, gslb := range gslbList.Items {
				for _, rule := range gslb.Spec.Ingress.Rules {
					for _, path := range rule.HTTP.Paths {
						if path.Backend.ServiceName == a.Meta.GetName() {
							gslbName = gslb.Name
						}
					}
				}
			}
			if len(gslbName) > 0 {
				return []reconcile.Request{
					{NamespacedName: types.NamespacedName{
						Name:      gslbName,
						Namespace: a.Meta.GetNamespace(),
					}},
				}
			}
			return nil
		})

	// Watch for Endpoints that are not controlled directly
	err = c.Watch(
		&source.Kind{Type: &corev1.Endpoints{}},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: mapFn,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
