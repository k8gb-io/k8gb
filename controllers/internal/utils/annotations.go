package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MergeAnnotations adds or updates annotations from defaultSource to defaultTarget
func MergeAnnotations(target *metav1.ObjectMeta, source *metav1.ObjectMeta) {
	if target.Annotations == nil {
		target.Annotations = make(map[string]string)
	}
	if source.Annotations == nil {
		source.Annotations = make(map[string]string)
	}
	for k, v := range source.Annotations {
		if target.Annotations[k] != v {
			target.Annotations[k] = v
		}
	}
}

// ContainsAnnotations checks if defaultTarget contains all annotations of defaultSource.
func ContainsAnnotations(target *metav1.ObjectMeta, source *metav1.ObjectMeta) bool {
	if source.Annotations == nil {
		return true
	}
	if target.Annotations == nil {
		return false
	}
	for k, v := range source.Annotations {
		if target.Annotations[k] != v {
			return false
		}
	}
	return true
}
