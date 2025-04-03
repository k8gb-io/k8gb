package controllers

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"context"
	"fmt"
	"strconv"

	"github.com/k8gb-io/k8gb/controllers/utils"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	"github.com/k8gb-io/k8gb/controllers/depresolver"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

// SetupWithManager configures controller manager
func (r *GslbReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Figure out Gslb resource name to Reconcile when non controlled Name is updated

	r.Recorder = mgr.GetEventRecorderFor("gslb-controller")

	endpointMapHandler := handler.EnqueueRequestsFromMapFunc(
		func(_ context.Context, a client.Object) []reconcile.Request {
			gslbList := &k8gbv1beta1.GslbList{}
			opts := []client.ListOption{
				client.InNamespace(a.GetNamespace()),
			}
			c := mgr.GetClient()
			err := c.List(context.TODO(), gslbList, opts...)
			if err != nil {
				log.Info().Msg("Can't fetch gslb objects")
				return nil
			}
			reconcileRequests := []reconcile.Request{}
		GslbLoop:
			for _, gslb := range gslbList.Items {
				for _, server := range gslb.Status.Servers {
					for _, service := range server.Services {
						if service.Name == a.GetName() {
							reconcileRequests = append(reconcileRequests, reconcile.Request{
								NamespacedName: types.NamespacedName{
									Name:      gslb.Name,
									Namespace: a.GetNamespace(),
								},
							})
							continue GslbLoop
						}
					}
				}
			}
			return reconcileRequests
		})

	ingressMapHandler := handler.EnqueueRequestsFromMapFunc(
		func(_ context.Context, a client.Object) []reconcile.Request {
			annotations := a.GetAnnotations()
			if annotationValue, found := annotations[strategyAnnotation]; found && a.GetOwnerReferences() != nil {
				c := mgr.GetClient()
				r.createGSLBFromIngress(c, a, annotationValue)
			}
			return nil
		})

	return ctrl.NewControllerManagedBy(mgr).
		For(&k8gbv1beta1.Gslb{}).
		Owns(&netv1.Ingress{}).
		Owns(&externaldns.DNSEndpoint{}).
		Watches(&corev1.Endpoints{}, endpointMapHandler).
		Watches(&netv1.Ingress{}, ingressMapHandler).
		WithEventFilter(predicate.Funcs{
			UpdateFunc: func(e event.TypedUpdateEvent[client.Object]) bool {
				if e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration() {
					return true
				}

				// endpoints don't have state, therefore they don't have a generation
				// but when their subsets change they must be be reconciled
				gvk, err := apiutil.GVKForObject(e.ObjectOld, r.Scheme)
				if err != nil {
					log.Warn().Msg("could not fetch GroupVersionKind for object")
				} else if gvk.Kind == "Endpoints" {
					return true
				}

				// Ignore reconciliation in case nothing has changed in k8gb annotations
				oldAnnotations := e.ObjectOld.GetAnnotations()
				newAnnotations := e.ObjectNew.GetAnnotations()
				reconcile := !utils.EqualPredefinedAnnotations(oldAnnotations, newAnnotations, k8gbAnnotations...)

				return reconcile
			},
		}).
		Complete(r)
}

func (r *GslbReconciler) createGSLBFromIngress(c client.Client, a client.Object, strategy string) {
	log.Info().
		Str("annotation", fmt.Sprintf("(%s:%s)", strategyAnnotation, strategy)).
		Str("ingress", a.GetName()).
		Msg("Detected strategy annotation on ingress")

	ingressToReuse := &netv1.Ingress{}
	err := c.Get(context.Background(), client.ObjectKey{
		Namespace: a.GetNamespace(),
		Name:      a.GetName(),
	}, ingressToReuse)
	if err != nil {
		log.Info().
			Str("ingress", a.GetName()).
			Msg("Ingress does not exist anymore. Skipping Glsb creation...")
		return
	}
	gslbExist := &k8gbv1beta1.Gslb{}
	err = c.Get(context.Background(), client.ObjectKey{
		Namespace: a.GetNamespace(),
		Name:      a.GetName(),
	}, gslbExist)
	if err == nil {
		log.Info().
			Str("gslb", gslbExist.Name).
			Msg("Gslb already exists. Skipping Gslb creation...")
		return
	}
	gslb := &k8gbv1beta1.Gslb{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: a.GetNamespace(),
			Name:      a.GetName(),
		},
		Spec: k8gbv1beta1.GslbSpec{
			Ingress: k8gbv1beta1.FromV1IngressSpec(ingressToReuse.Spec),
		},
	}

	gslb.Spec.Strategy, err = r.parseStrategy(a.GetAnnotations(), strategy)
	if err != nil {
		log.Err(err).
			Str("gslb", gslbExist.Name).
			Msg("can't parse Gslb strategy")
		return
	}

	err = controllerutil.SetControllerReference(ingressToReuse, gslb, r.Scheme)
	if err != nil {
		log.Err(err).
			Str("ingress", ingressToReuse.Name).
			Str("gslb", gslb.Name).
			Msg("Cannot set the Ingress as the owner of the Gslb")
	}

	log.Info().
		Str("gslb", gslb.Name).
		Msg(fmt.Sprintf("Creating a new Gslb out of Ingress with '%s' annotation", strategyAnnotation))
	err = c.Create(context.Background(), gslb)
	if err != nil {
		log.Err(err).Msg("Glsb creation failed")
	}
}

func (r *GslbReconciler) parseStrategy(annotations map[string]string, strategy string) (result k8gbv1beta1.Strategy, err error) {
	toInt := func(k string, v string) (int, error) {
		intValue, err := strconv.Atoi(v)
		if err != nil {
			return -1, fmt.Errorf("can't parse annotation value %s to int for key %s", v, k)
		}
		return intValue, nil
	}

	result = k8gbv1beta1.Strategy{
		Type: strategy,
	}

	for annotationKey, annotationValue := range annotations {
		switch annotationKey {
		case dnsTTLSecondsAnnotation:
			if result.DNSTtlSeconds, err = toInt(annotationKey, annotationValue); err != nil {
				return result, err
			}
		case primaryGeoTagAnnotation:
			result.PrimaryGeoTag = annotationValue
		}
	}

	if strategy == depresolver.FailoverStrategy {
		if len(result.PrimaryGeoTag) == 0 {
			return result, fmt.Errorf("%s strategy requires annotation %s", depresolver.FailoverStrategy, primaryGeoTagAnnotation)
		}
	}

	return result, nil
}
