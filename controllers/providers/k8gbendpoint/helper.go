package k8gbendpoint

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

	"github.com/k8gb-io/k8gb/controllers/utils"
	"github.com/rs/zerolog"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

func saveDNSEndpoint(ctx context.Context, client client.Client, namespace string, i *externaldns.DNSEndpoint, logger *zerolog.Logger) error {
	found := &externaldns.DNSEndpoint{}
	err := client.Get(ctx, types.NamespacedName{
		Name:      i.Name,
		Namespace: namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the DNSEndpoint
		logger.Info().
			Interface("DNSEndpoint", i).
			Msgf("Creating a new DNSEndpoint")
		err = client.Create(ctx, i)

		if err != nil {
			// Creation failed
			logger.Err(err).
				Str("namespace", i.Namespace).
				Str("name", i.Name).
				Msg("Failed to create new DNSEndpoint")
			return err
		}
		// Creation was successful
		return nil
	} else if err != nil {
		// Error that isn't due to the service not existing
		logger.Err(err).Msg("Failed to get DNSEndpoint")
		return err
	}

	// Update existing object with new spec, labels and annotations
	found.Spec = i.Spec
	found.ObjectMeta.Annotations = i.ObjectMeta.Annotations
	found.ObjectMeta.Labels = i.ObjectMeta.Labels
	err = client.Update(context.TODO(), found)

	if err != nil {
		// Update failed
		logger.Err(err).
			Str("namespace", found.Namespace).
			Str("name", found.Name).
			Msg("Failed to update DNSEndpoint")
		return err
	}
	return nil
}

func removeEndpoint(ctx context.Context, client client.Client, endpointKey client.ObjectKey, logger *zerolog.Logger) error {
	logger.Info().
		Str("namespace", endpointKey.Namespace).
		Str("name", endpointKey.Name).
		Msg("Removing endpoint")
	dnsEndpoint := &externaldns.DNSEndpoint{}
	err := client.Get(ctx, endpointKey, dnsEndpoint)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Warn().
				Str("namespace", endpointKey.Namespace).
				Str("name", endpointKey.Name).
				Err(err).
				Msg("Endpoint not found")
			return nil
		}
		return err
	}
	err = client.Delete(ctx, dnsEndpoint)
	return err
}

func getNSCombinations(original []utils.DNSServer, hostToUse string) []utils.DNSServer {
	portToUse := original[0].Port
	nameServerToUse := []utils.DNSServer{
		{
			Host: hostToUse,
			Port: portToUse,
		},
	}
	defaultPortAdded := false
	for _, s := range original {
		if s.Port != 53 {
			nameServerToUse = append(nameServerToUse, utils.DNSServer{
				Host: hostToUse,
				Port: s.Port,
			})
		} else {
			defaultPortAdded = true
		}
	}
	if !defaultPortAdded {
		nameServerToUse = append(nameServerToUse, utils.DNSServer{
			Host: hostToUse,
			Port: 53,
		})
	}
	return nameServerToUse
}
