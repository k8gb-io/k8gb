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

package main

import (
	"os"

	str "github.com/AbsaOSS/gopkg/strings"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	"github.com/k8gb-io/k8gb/controllers"
	"github.com/k8gb-io/k8gb/controllers/depresolver"
	"github.com/k8gb-io/k8gb/controllers/logging"
	"github.com/k8gb-io/k8gb/controllers/providers/dns"
	"github.com/k8gb-io/k8gb/controllers/providers/metrics"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	externaldns "sigs.k8s.io/external-dns/endpoint"
	// +kubebuilder:scaffold:imports
)

var (
	runtimescheme = runtime.NewScheme()
	version       = "development"
	commit        = "none"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(runtimescheme))

	utilruntime.Must(k8gbv1beta1.AddToScheme(runtimescheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var exitCode = 1
	defer func() {
		os.Exit(exitCode)
	}()
	var f *dns.ProviderFactory
	resolver := depresolver.NewDependencyResolver()
	config, err := resolver.ResolveOperatorConfig()
	deprecations := resolver.GetDeprecations()
	// Initialize desired log or default log in case of configuration failed.
	logging.Init(config)
	log := logging.Logger()
	log.Info().
		Str("version", version).
		Str("commit", commit).
		Msg("K8gb status")
	if err != nil {
		log.Err(err).Msg("can't resolve environment variables")
		return
	}
	log.Debug().
		Str("config", str.ToString(config)).
		Msg("Resolved config")

	ctrl.SetLogger(logging.NewLogrAdapter(log))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             runtimescheme,
		MetricsBindAddress: config.MetricsAddress,
		Port:               9443,
		LeaderElection:     false,
		LeaderElectionID:   "8020e9ff.absa.oss",
	})
	if err != nil {
		log.Err(err).Msg("Unable to start manager")
		return
	}

	for _, d := range deprecations {
		log.Warn().Msg(d)
	}

	log.Info().Msg("Registering components")

	// Add external-dns DNSEndpoints resource
	// https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#adding-3rd-party-resources-to-your-operator
	schemeBuilder := &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}}
	schemeBuilder.Register(&externaldns.DNSEndpoint{}, &externaldns.DNSEndpointList{})
	if err := schemeBuilder.AddToScheme(mgr.GetScheme()); err != nil {
		log.Err(err).Msg("Extending scheme")
		return
	}

	reconciler := &controllers.GslbReconciler{
		Config:      config,
		Client:      mgr.GetClient(),
		DepResolver: resolver,
		Scheme:      mgr.GetScheme(),
	}

	log.Info().Msg("Starting metrics")
	metrics.Init(config)
	defer metrics.Metrics().Unregister()
	err = metrics.Metrics().Register()
	if err != nil {
		log.Err(err).Msg("Register metrics error")
		return
	}

	log.Info().Msg("Resolving DNS provider")
	f, err = dns.NewDNSProviderFactory(reconciler.Client, *reconciler.Config)
	if err != nil {
		log.Err(err).Msg("Unable to create factory")
		return
	}
	reconciler.DNSProvider = f.Provider()
	log.Info().Str("provider", reconciler.DNSProvider.String()).Msg("Started")

	if err = reconciler.SetupWithManager(mgr); err != nil {
		log.Err(err).Msg("Unable to create controller Gslb")
		return
	}
	metrics.Metrics().SetRuntimeInfo(version, commit)
	// +kubebuilder:scaffold:builder
	log.Info().Msg("Starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Err(err).Msg("Problem running manager")
		return
	}
	// time to call deferred functions including the exit one with code=0
	exitCode = 0
}
