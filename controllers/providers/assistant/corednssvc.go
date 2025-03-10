package assistant

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
	coreerrors "errors"

	"github.com/k8gb-io/k8gb/controllers/depresolver"

	"github.com/k8gb-io/k8gb/controllers/utils"

	"github.com/k8gb-io/k8gb/controllers/logging"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const coreDNSServiceLabel = "app.kubernetes.io/name=coredns"

// CoreDNSService is common wrapper operating on GSLB instance.
// It uses apimachinery client to call kubernetes API
type CoreDNSService struct {
	client client.Client
	config depresolver.Config
}

var log = logging.Logger()

func NewCoreDNSServiceAssistant(client client.Client, config depresolver.Config) *CoreDNSService {
	return &CoreDNSService{
		client: client,
		config: config,
	}
}

// GetResource returns the CoreDNS Service
func (r *CoreDNSService) GetResource() (*corev1.Service, error) {
	serviceList := &corev1.ServiceList{}
	sel, err := labels.Parse(coreDNSServiceLabel)
	if err != nil {
		log.Err(err).Msg("Badly formed label selector")
		return nil, err
	}
	listOption := &client.ListOptions{
		LabelSelector: sel,
		Namespace:     r.config.K8gbNamespace,
	}

	err = r.client.List(context.TODO(), serviceList, listOption)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Warn().Err(err).Msg("Can't find CoreDNS service")
		}
	}
	if len(serviceList.Items) != 1 {
		log.Warn().Msg("More than 1 CoreDNS service was found")
		for _, service := range serviceList.Items {
			log.Info().
				Str("serviceName", service.Name).
				Msg("Found CoreDNS service")
		}
		err := coreerrors.New("more than 1 CoreDNS service was found. Check if CoreDNS exposed correctly")
		return nil, err
	}
	coreDNSService := &serviceList.Items[0]
	return coreDNSService, nil
}

// GetExposedIPs retrieves list of IP's exposed by CoreDNS
func (r *CoreDNSService) GetExposedIPs() ([]string, error) {
	coreDNSService, err := r.GetResource()
	if err != nil {
		return nil, err
	}
	if coreDNSService.Spec.Type == "ClusterIP" {
		if len(coreDNSService.Spec.ClusterIPs) == 0 {
			errMessage := "no ClusterIPs found"
			log.Warn().
				Str("serviceName", coreDNSService.Name).
				Msg(errMessage)
			err := coreerrors.New(errMessage)
			return nil, err
		}
		return coreDNSService.Spec.ClusterIPs, nil
	}
	// LoadBalancer / ExternalName / NodePort service
	if len(coreDNSService.Status.LoadBalancer.Ingress) == 0 {
		errMessage := "no LoadBalancer ExternalIPs are found"
		log.Warn().
			Str("serviceName", coreDNSService.Name).
			Msg(errMessage)
		err := coreerrors.New(errMessage)
		return nil, err
	}

	var ipList []string
	for _, ingressStatusIP := range coreDNSService.Status.LoadBalancer.Ingress {
		var confirmedIPs, err = extractIPFromLB(ingressStatusIP, r.config.EdgeDNSServers)
		if err != nil {
			return nil, err
		}
		ipList = append(ipList, confirmedIPs...)
	}
	return ipList, nil

}

func extractIPFromLB(lb corev1.LoadBalancerIngress, ns utils.DNSList) (ips []string, err error) {
	if lb.Hostname != "" {
		IPs, err := utils.Dig(lb.Hostname, 8, ns...)
		if err != nil {
			log.Warn().Err(err).
				Str("loadBalancerHostname", lb.Hostname).
				Msg("Can't dig CoreDNS service LoadBalancer FQDN")
			return nil, err
		}
		return IPs, nil
	}
	if lb.IP != "" {
		return []string{lb.IP}, nil
	}
	return nil, nil
}
