package common

import (
	"fmt"
	"github.com/k8gb-io/k8gb/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetListOptions(rr v1beta1.ResourceRef, namespace string) (*client.ListOptions, error) {
	var filterApplied = false
	opts := &client.ListOptions{
		Namespace: namespace,
	}
	if rr.LabelSelector.MatchLabels != nil || len(rr.LabelSelector.MatchExpressions) > 0 {
		filterApplied = true
		selector, err := metav1.LabelSelectorAsSelector(&rr.LabelSelector)
		if err != nil {
			return nil, err
		}
		opts.LabelSelector = selector
	}
	if rr.Name != "" {
		filterApplied = true
		opts.FieldSelector = fields.OneTermEqualSelector("metadata.name", rr.Name)
	}
	if !filterApplied {
		return nil, fmt.Errorf("gslb spec must contain label selector or name")
	}
	return opts, nil
}
