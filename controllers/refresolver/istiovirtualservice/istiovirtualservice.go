package istiovirtualservice

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
	"strings"

	"github.com/k8gb-io/k8gb/controllers/refresolver/ingress/common"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	"github.com/k8gb-io/k8gb/controllers/logging"
	"github.com/k8gb-io/k8gb/controllers/utils"
	istio "istio.io/client-go/pkg/apis/networking/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var log = logging.Logger()

const (
	// comma separated list of external IP addresses
	externalIPsAnnotation = "k8gb.io/exposed-ip-addresses"
)

type ReferenceResolver struct {
	virtualService *istio.VirtualService
	lbService      *corev1.Service
}

// NewReferenceResolver creates a new reference resolver capable of understanding `networking.istio.io/v1` resources
func NewReferenceResolver(gslb *k8gbv1beta1.Gslb, k8sClient client.Client) (*ReferenceResolver, error) {
	virtualServiceList, err := getGslbVirtualServiceRef(gslb, k8sClient)
	if err != nil {
		return nil, err
	}

	if len(virtualServiceList) != 1 {
		return nil, fmt.Errorf("exactly 1 VirtualService resource expected but %d were found", len(virtualServiceList))
	}
	virtualService := virtualServiceList[0]

	gateway, err := getGateway(virtualService, k8sClient)
	if err != nil {
		return nil, err
	}

	lbService, err := getLbService(gateway, k8sClient)
	if err != nil {
		return nil, err
	}

	return &ReferenceResolver{
		virtualService: virtualService,
		lbService:      lbService,
	}, nil
}

// getGslbVirtualServiceRef resolves an istio virtual service resource referenced by the Gslb spec
func getGslbVirtualServiceRef(gslb *k8gbv1beta1.Gslb, k8sClient client.Client) ([]*istio.VirtualService, error) {
	virtualServiceList := &istio.VirtualServiceList{}
	opts, err := common.GetListOptions(gslb.Spec.ResourceRef, gslb.Namespace)
	if err != nil {
		return nil, err
	}
	err = k8sClient.List(context.TODO(), virtualServiceList, opts)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info().
				Str("gslb", gslb.Name).
				Msg("Can't find referenced VirtualService resource")
		}
		return nil, err
	}

	return virtualServiceList.Items, err
}

// getGateway retrieves the istio gateway referenced by the istio virtual service
func getGateway(virtualService *istio.VirtualService, k8sClient client.Client) (*istio.Gateway, error) {
	var ingressGateways []string
	for _, gateway := range virtualService.Spec.Gateways {
		// count only dedicated ingress gateways
		if gateway != "mesh" {
			ingressGateways = append(ingressGateways, gateway)
		}
	}

	if len(ingressGateways) != 1 {
		return nil, fmt.Errorf("expected exactly 1 Gateway to be referenced by the VirtualService but %d were found", len(ingressGateways))
	}
	gatewayRef := strings.Split(ingressGateways[0], "/")
	gatewayNamespace := gatewayRef[0]
	gatewayName := gatewayRef[1]

	gateway := &istio.Gateway{}
	err := k8sClient.Get(context.TODO(), types.NamespacedName{
		Namespace: gatewayNamespace,
		Name:      gatewayName,
	}, gateway)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info().
				Str("gatewayNamespace", gatewayNamespace).
				Str("gatewayName", gatewayName).
				Msg("Can't find Gateway resource referenced by VirtualService")
		}
		return nil, err
	}

	return gateway, nil
}

// getLbService retrieves the kubernetes service referenced by an istio gateway
func getLbService(gateway *istio.Gateway, k8sClient client.Client) (*corev1.Service, error) {
	gatewayServiceList := &corev1.ServiceList{}
	opts := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(gateway.Spec.Selector),
	}

	err := k8sClient.List(context.TODO(), gatewayServiceList, opts)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info().
				Str("gateway", gateway.Name).
				Msg("Can't find any service with the gateway's selector")
		}
		return nil, err
	}

	if len(gatewayServiceList.Items) != 1 {
		return nil, fmt.Errorf("exactly 1 Service resource expected but %d were found", len(gatewayServiceList.Items))
	}

	gatewayService := &gatewayServiceList.Items[0]
	return gatewayService, nil
}

// GetServers retrieves the GSLB server configuration from the istio virtual service resource
func (rr *ReferenceResolver) GetServers() ([]*k8gbv1beta1.Server, error) {
	hosts := rr.virtualService.Spec.Hosts
	if len(hosts) < 1 {
		return nil, fmt.Errorf("can't find hosts in VirtualService %s", rr.virtualService.Name)
	}

	servers := []*k8gbv1beta1.Server{}
	for _, host := range hosts {
		server := &k8gbv1beta1.Server{
			Host:     host,
			Services: []*k8gbv1beta1.NamespacedName{},
		}
		for _, http := range rr.virtualService.Spec.Http {
			for _, route := range http.Route {
				server.Services = append(server.Services, &k8gbv1beta1.NamespacedName{
					Name:      route.Destination.Host,
					Namespace: rr.virtualService.Namespace,
				})
			}
		}
		for _, tls := range rr.virtualService.Spec.Tls {
			for _, route := range tls.Route {
				server.Services = append(server.Services, &k8gbv1beta1.NamespacedName{
					Name:      route.Destination.Host,
					Namespace: rr.virtualService.Namespace,
				})
			}
		}
		for _, tcp := range rr.virtualService.Spec.Tcp {
			for _, route := range tcp.Route {
				server.Services = append(server.Services, &k8gbv1beta1.NamespacedName{
					Name:      route.Destination.Host,
					Namespace: rr.virtualService.Namespace,
				})
			}
		}
		servers = append(servers, server)
	}

	return servers, nil
}

// GetGslbExposedIPs retrieves the load balancer IP address of the GSLB
func (rr *ReferenceResolver) GetGslbExposedIPs(gslbAnnotations map[string]string, edgeDNSServers utils.DNSList) ([]string, error) {
	// fetch the IP addresses of the reverse proxy from an annotation if it exists
	if ingressIPsFromAnnotation, ok := gslbAnnotations[externalIPsAnnotation]; ok {
		return utils.ParseIPAddresses(ingressIPsFromAnnotation)
	}

	// if there is no annotation -> fetch the IP addresses from the Status of the Ingress resource
	gslbIngressIPs := []string{}
	for _, ip := range rr.lbService.Status.LoadBalancer.Ingress {
		if len(ip.IP) > 0 {
			gslbIngressIPs = append(gslbIngressIPs, ip.IP)
		}
		if len(ip.Hostname) > 0 {
			IPs, err := utils.Dig(ip.Hostname, 8, edgeDNSServers...)
			if err != nil {
				log.Warn().Err(err).Msg("Dig error")
				return nil, err
			}
			gslbIngressIPs = append(gslbIngressIPs, IPs...)
		}
	}

	return gslbIngressIPs, nil
}
