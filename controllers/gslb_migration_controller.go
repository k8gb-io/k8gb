package controllers

import (
	"context"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	k8gbv1beta1io "github.com/k8gb-io/k8gb/api/v1beta1io"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LegacyGslbMigrationReconciler struct {
	client.Client
	Recorder record.EventRecorder
}

func (r *LegacyGslbMigrationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	legacy := &k8gbv1beta1.Gslb{}
	if err := r.Get(ctx, req.NamespacedName, legacy); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if legacy.Labels != nil && legacy.Labels[migrationLabelKey] == "true" {
		r.Recorder.Eventf(legacy, corev1.EventTypeWarning, "LegacyIgnored", "Legacy Gslb ignored. Edit k8gb.io Gslb %s/%s instead.", legacy.Namespace, legacy.Name)
		return ctrl.Result{}, nil
	}

	desired := convertGslbLegacyToIO(legacy)

	existing := &k8gbv1beta1io.Gslb{}
	if err := r.Get(ctx, types.NamespacedName{Name: legacy.Name, Namespace: legacy.Namespace}, existing); err != nil {
		if apierrors.IsNotFound(err) {
			if err := r.Create(ctx, desired); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			return ctrl.Result{}, err
		}
	} else {
		patch := client.MergeFrom(existing.DeepCopy())
		existing.Spec = desired.Spec
		if err := r.Patch(ctx, existing, patch); err != nil {
			return ctrl.Result{}, err
		}
	}

	legacyPatch := client.MergeFrom(legacy.DeepCopy())
	if legacy.Labels == nil {
		legacy.Labels = map[string]string{}
	}
	legacy.Labels[migrationLabelKey] = "true"
	if err := r.Patch(ctx, legacy, legacyPatch); err != nil {
		return ctrl.Result{}, err
	}

	r.Recorder.Eventf(legacy, corev1.EventTypeNormal, "LegacyMigrated", "Legacy Gslb migrated to k8gb.io Gslb %s/%s. Edit the k8gb.io object going forward.", legacy.Namespace, legacy.Name)
	return ctrl.Result{}, nil
}

func (r *LegacyGslbMigrationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8gbv1beta1.Gslb{}).
		Complete(r)
}
