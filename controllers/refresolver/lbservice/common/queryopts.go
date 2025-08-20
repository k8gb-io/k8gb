package common

import (
	"fmt"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type QueryMode int

const (
	QueryModeGet QueryMode = iota
	QueryModeList
)

type QueryOptions struct {
	Mode     QueryMode
	GetKey   *types.NamespacedName
	ListOpts []client.ListOption
}

// GetQueryOptions creates query options from GSLB ResourceRef
func GetQueryOptions(resourceRef k8gbv1beta1.ResourceRef, namespace string) (*QueryOptions, error) {
	if resourceRef.Name != "" {
		// Direct reference by name
		return &QueryOptions{
			Mode: QueryModeGet,
			GetKey: &types.NamespacedName{
				Name:      resourceRef.Name,
				Namespace: resourceRef.Namespace,
			},
		}, nil
	}

	if resourceRef.MatchLabels != nil || resourceRef.MatchExpressions != nil {
		// Reference by label selector
		listOpts := []client.ListOption{
			client.InNamespace(namespace),
		}

		if resourceRef.MatchLabels != nil || resourceRef.MatchExpressions != nil {
			selector, err := metav1.LabelSelectorAsSelector(&resourceRef.LabelSelector)
			if err != nil {
				return nil, fmt.Errorf("invalid label selector: %w", err)
			}
			listOpts = append(listOpts, client.MatchingLabelsSelector{Selector: selector})
		}

		return &QueryOptions{
			Mode:     QueryModeList,
			ListOpts: listOpts,
		}, nil
	}

	return nil, fmt.Errorf("neither name nor label selector specified in ResourceRef")
}
