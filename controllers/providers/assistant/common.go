package assistant

/*
Copyright 2021-2025 The k8gb Contributors.

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

	"github.com/k8gb-io/k8gb/controllers/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getObjectList[T client.ObjectList](ctx context.Context, cl client.Client, list T, label, namespace string) (T, error) {
	sel, err := labels.Parse(label)
	if err != nil {
		return list, fmt.Errorf("badly formed label selector %s (%s)", label, err)
	}
	listOption := &client.ListOptions{
		LabelSelector: sel,
		Namespace:     namespace,
	}
	err = cl.List(ctx, list, listOption)
	if err != nil {
		return list, fmt.Errorf("listing %s (%s)", label, err)
	}
	return list, nil
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
