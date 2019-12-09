package gslb

import (
	"context"

	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_gslb")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Gslb Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileGslb{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("gslb-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Gslb
	err = c.Watch(&source.Kind{Type: &ohmyglbv1beta1.Gslb{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Ingress
	err = c.Watch(&source.Kind{Type: &v1beta1.Ingress{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &ohmyglbv1beta1.Gslb{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileGslb implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileGslb{}

// ReconcileGslb reconciles a Gslb object
type ReconcileGslb struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Gslb object and makes changes based on the state read
// and what is in the Gslb.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileGslb) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Gslb")

	// Fetch the Gslb instance
	gslb := &ohmyglbv1beta1.Gslb{}
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

	// == Ingress ==========
	result, err = r.ensureIngress(
		request,
		gslb,
		r.gslbIngress(gslb))
	if result != nil {
		return *result, err
	}

	// == Status =
	err = r.updateManagedHostsStatus(gslb)
	if err != nil {
		// Requeue the request
		return reconcile.Result{}, err
	}

	// == Finish ==========
	// Everything went fine, don't requeue
	return reconcile.Result{}, nil
}
