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

	"github.com/k8gb-io/k8gb/controllers/boot"
	"github.com/k8gb-io/k8gb/controllers/depresolver"
	"github.com/k8gb-io/k8gb/controllers/providers/dns"
	"github.com/k8gb-io/k8gb/controllers/utils"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CoreDNSReconciler struct {
	client.Client
	Config      *depresolver.Config
	Recorder    record.EventRecorder
	Tracer      trace.Tracer
	Scheme      *runtime.Scheme
	DNSProvider dns.Provider
	Logger      *zerolog.Logger
	Bootstrap   *boot.Bootstrap
}

func (c *CoreDNSReconciler) Reconcile(_ context.Context, req ctrl.Request) (ctrl.Result, error) {
	// todo: introduce variable for reconciliation interval
	result := utils.NewReconcileResultHandler(0)
	c.Logger.Info().Msgf("Reconciling '%s' %s", req, c.Bootstrap.IPs)
	for _, zoneInfo := range c.Config.DelegationZones {
		err := c.DNSProvider.CreateZoneDelegation(zoneInfo)
		if err != nil {
			c.Logger.Err(err).Msg("Error creating zone delegation")
		}
	}
	return result.Stop()
}
