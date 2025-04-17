package common

import (
	"fmt"
	"github.com/k8gb-io/k8gb/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type QueryMode int

const (
	QueryModeList QueryMode = iota
	QueryModeGet
)

type QueryOptions struct {
	Mode     QueryMode
	GetKey   *types.NamespacedName
	ListOpts *client.ListOptions
}

func GetQueryOptions(rr v1beta1.ResourceRef, namespace string) (*QueryOptions, error) {
	if rr.Name != "" && rr.LabelSelector.MatchLabels == nil && len(rr.LabelSelector.MatchExpressions) == 0 {
		return &QueryOptions{
			Mode:   QueryModeGet,
			GetKey: &types.NamespacedName{Name: rr.Name, Namespace: namespace},
		}, nil
	}

	if rr.LabelSelector.MatchLabels != nil || len(rr.LabelSelector.MatchExpressions) > 0 {
		selector, err := metav1.LabelSelectorAsSelector(&rr.LabelSelector)
		if err != nil {
			return nil, err
		}
		opts := &client.ListOptions{
			Namespace:     namespace,
			LabelSelector: selector,
		}
		return &QueryOptions{
			Mode:     QueryModeList,
			ListOpts: opts,
		}, nil
	}

	return nil, fmt.Errorf("gslb spec must contain label selector or name")
}
