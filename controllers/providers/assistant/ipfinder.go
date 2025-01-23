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
	"fmt"
	"strings"

	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IPFinder interface {
	Find(path string) ([]string, error)
}

type IngressFiner struct {
	client  client.Client
	context context.Context
}

func NewIngressFiner(context context.Context, client client.Client) *IngressFiner {
	return &IngressFiner{context: context, client: client}
}

func (f *IngressFiner) Find(path string) (ips []string, err error) {
	arr := strings.Split(path, ".")
	if len(arr) != 2 {
		return nil, fmt.Errorf("path format error (namespace.name): %s", path)
	}
	selector := client.ObjectKey{Namespace: arr[0], Name: arr[1]}
	ingress := &netv1.Ingress{}
	err = f.client.Get(f.context, selector, ingress)
	if err != nil {
		return nil, fmt.Errorf("finding ingress %s (%s)", path, err)
	}
	for _, lb := range ingress.Status.LoadBalancer.Ingress {
		if lb.IP != "" {
			ips = append(ips, lb.IP)
		}
	}
	return ips, nil
}
