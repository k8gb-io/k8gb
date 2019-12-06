package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GslbSpec defines the desired state of Gslb
// +k8s:openapi-gen=true
type GslbSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Host string `json:"host"`
	HTTP struct {
		Paths []struct {
			Backend struct {
				ServiceName string `json:"serviceName"`
				ServicePort string `json:"servicePort"`
			} `json:"backend"`
			Path string `json:"path"`
		} `json:"paths"`
	} `json:"http"`
	Strategy string `json:"strategy"`
	TLS      struct {
		SecretName string `json:"secretName"`
	} `json:"tls"`
}

// GslbStatus defines the observed state of Gslb
// +k8s:openapi-gen=true
type GslbStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +listType=set
	Hosts []string `json:"hosts"`
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
