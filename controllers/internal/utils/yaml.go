// Package utils implements common, reusable helpers
package utils

import (
	"encoding/json"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	yamlConv "github.com/ghodss/yaml"
)

// YamlToGslb takes yaml and returns Gslb object
func YamlToGslb(yaml []byte) (*k8gbv1beta1.Gslb, error) {
	// yamlBytes contains a []byte of my yaml job spec
	// convert the yaml to json
	jsonBytes, err := yamlConv.YAMLToJSON(yaml)
	if err != nil {
		return &k8gbv1beta1.Gslb{}, err
	}
	// unmarshal the json into the kube struct
	gslb := &k8gbv1beta1.Gslb{}
	err = json.Unmarshal(jsonBytes, &gslb)
	if err != nil {
		return &k8gbv1beta1.Gslb{}, err
	}
	return gslb, nil
}
