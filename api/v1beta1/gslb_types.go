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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Strategy defines Gslb behavior
// +k8s:openapi-gen=true
type Strategy struct {
	// Load balancing strategy type:(roundRobin|failover|geoip)
	// +kubebuilder:validation:Enum=failover;geoip;roundRobin
	Type string `json:"type"`
	// Primary Geo Tag. Valid for failover strategy only
	PrimaryGeoTag string `json:"primaryGeoTag,omitempty"`
	// List of additional geo tags that should be tried in order if PrimaryGeoTag fails. Valid for failover strategy only
	FailoverOrder []string `json:"failoverOrder,omitempty"`
	// Defines DNS record TTL in seconds
	// +kubebuilder:validation:Minimum=0
	DNSTtlSeconds int `json:"dnsTtlSeconds,omitempty"`
	// Split brain TXT record expiration in seconds
	// +kubebuilder:validation:Minimum=0
	SplitBrainThresholdSeconds int `json:"splitBrainThresholdSeconds,omitempty"`
}

// GslbSpec defines the desired state of Gslb
// +k8s:openapi-gen=true
type GslbSpec struct {
	// Gslb-enabled Ingress Spec
	Ingress IngressSpec `json:"ingress"`
	// Gslb Strategy spec
	Strategy Strategy `json:"strategy"`
}

// GslbStatus defines the observed state of Gslb
type GslbStatus struct {
	// Associated Service status
	ServiceHealth map[string]HealthStatus `json:"serviceHealth"`
	// Current Healthy DNS record structure
	HealthyRecords map[string][]string `json:"healthyRecords"`
	// Cluster Geo Tag
	GeoTag string `json:"geoTag"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Gslb is the Schema for the gslbs API
// +kubebuilder:printcolumn:name="strategy",type=string,JSONPath=`.spec.strategy.type`
// +kubebuilder:printcolumn:name="geoTag",type=string,JSONPath=`.status.geoTag`
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

type HealthStatus string

const (
	Healthy   HealthStatus = "Healthy"
	Unhealthy HealthStatus = "Unhealthy"
	NotFound  HealthStatus = "NotFound"
)

func (h HealthStatus) String() string {
	return string(h)
}

func init() {
	SchemeBuilder.Register(&Gslb{}, &GslbList{})
}
