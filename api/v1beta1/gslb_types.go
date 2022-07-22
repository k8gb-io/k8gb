package v1beta1

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
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Strategy defines Gslb behavior
// +k8s:openapi-gen=true
type Strategy struct {
	// Load balancing strategy type:(roundRobin|failover)
	Type string `json:"type"`
	// roundrobin and in the future consistent may (but also may not) contain a weight
	Weight Weight `json:"weight,omitempty"`
	// Primary Geo Tag. Valid for failover strategy only
	PrimaryGeoTag string `json:"primaryGeoTag,omitempty"`
	// Defines DNS record TTL in seconds
	DNSTtlSeconds int `json:"dnsTtlSeconds,omitempty"`
	// Split brain TXT record expiration in seconds
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
	// Comma-separated list of hosts. Duplicating the value from range .spec.ingress.rules[*].host for printer column
	Hosts string `json:"hosts,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Gslb is the Schema for the gslbs API
// +kubebuilder:printcolumn:name="strategy",type=string,JSONPath=`.spec.strategy.type`
// +kubebuilder:printcolumn:name="geoTag",type=string,JSONPath=`.status.geoTag`
// +kubebuilder:printcolumn:name="hosts",type=string,JSONPath=`.status.hosts`,priority=1
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

type Percentage string

type Weight map[string]Percentage

func (p Percentage) String() string {
	return string(p)
}

func (p Percentage) IsEmpty() bool {
	return string(p) == ""
}

func (p Percentage) TryParse() (v int, err error) {
	return strconv.Atoi(strings.TrimSuffix(strings.ReplaceAll(p.String(), " ", ""), "%"))
}

func (p Percentage) Int() int {
	v, _ := p.TryParse()
	return v
}

func (w Weight) IsEmpty() bool {
	return len(w) == 0
}
