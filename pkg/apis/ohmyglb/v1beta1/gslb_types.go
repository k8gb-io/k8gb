package v1beta1

import (
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Strategy defines Gslb behavior
// +k8s:openapi-gen=true
type Strategy struct {
	Type          string `json:"type"`
	PrimaryGeoTag string `json:"primaryGeoTag,omitempty"`
	// Defines DNS record TTL in seconds
	DNSTtlSeconds int `json:"dnsTtlSeconds,omitempty"`
	// Split brain TXT record expiration in seconds
	SplitBrainThresholdSeconds int `json:"splitBrainThresholdSeconds,omitempty"`
}
// GslbSpec defines the desired state of Gslb
// +k8s:openapi-gen=true
type GslbSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Ingress  v1beta1.IngressSpec `json:"ingress"`
	Strategy Strategy            `json:"strategy"`
}

// GslbStatus defines the observed state of Gslb
// +k8s:openapi-gen=true
type GslbStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +listType=set
	ManagedHosts   []string            `json:"managedHosts"`
	ServiceHealth  map[string]string   `json:"serviceHealth"`
	HealthyRecords map[string][]string `json:"healthyRecords"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Gslb is the Schema for the gslbs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=gslbs,scope=Namespaced
type Gslb struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GslbSpec   `json:"spec,omitempty"`
	Status GslbStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GslbList contains a list of Gslb
type GslbList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gslb `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Gslb{}, &GslbList{})
}
