package utils

import (
	"time"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ReconcileResultHandler struct {
	log           logr.Logger
	delayedResult ctrl.Result
}

func NewReconcileResultHandler(reconcileAfter int, log logr.Logger) *ReconcileResultHandler {
	return &ReconcileResultHandler{
		delayedResult: ctrl.Result{RequeueAfter: time.Second * time.Duration(reconcileAfter)},
		log:           log,
	}
}

// Stop stops reconciliation loop
func (r *ReconcileResultHandler) Stop() (ctrl.Result, error) {
	r.log.Info("reconciler exit")
	return ctrl.Result{}, nil
}

// RequeueError requeue loop immediately
// see default controller limiter: https://danielmangum.com/posts/controller-runtime-client-go-rate-limiting/
func (r *ReconcileResultHandler) RequeueError(err error) (ctrl.Result, error) {
	// logging error is handled in caller function
	return ctrl.Result{}, err
}

// Requeue requeue loop after config.ReconcileRequeueSeconds
// this apply in case you didn't modify request resources.
// If so, reconciliation starts immediately
// see: https://github.com/operator-framework/operator-sdk/issues/1164
func (r *ReconcileResultHandler) Requeue() (ctrl.Result, error) {
	return r.delayedResult, nil
}

func (r *ReconcileResultHandler) RequeueNow() (ctrl.Result, error) {
	return ctrl.Result{Requeue: true}, nil
}
