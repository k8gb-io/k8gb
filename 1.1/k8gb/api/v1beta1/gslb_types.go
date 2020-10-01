/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	// Important: Run "make" to regenerate code after modifying this file

	Ingress  v1beta1.IngressSpec `json:"ingress"`
	Strategy Strategy            `json:"strategy"`
}

// GslbStatus defines the observed state of Gslb
type GslbStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ServiceHealth  map[string]string   `json:"serviceHealth"`
	HealthyRecords map[string][]string `json:"healthyRecords"`
	GeoTag         string              `json:"geoTag"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Gslb is the Schema for the gslbs API
type Gslb struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GslbSpec   `json:"spec,omitempty"`
	Status GslbStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GslbList contains a list of Gslb
type GslbList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gslb `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Gslb{}, &GslbList{})
}
