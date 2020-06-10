//Package utils implements common, reusable helpers
package utils

import (
	"encoding/json"

	kgbv1beta1 "github.com/AbsaOSS/kgb/pkg/apis/kgb/v1beta1"
	yamlConv "github.com/ghodss/yaml"
)

//YamlToGslb takes yaml and returns Gslb object
func YamlToGslb(yaml []byte) (*kgbv1beta1.Gslb, error) {
	// yamlBytes contains a []byte of my yaml job spec
	// convert the yaml to json
	jsonBytes, err := yamlConv.YAMLToJSON(yaml)
	if err != nil {
		return &kgbv1beta1.Gslb{}, err
	}
	// unmarshal the json into the kube struct
	gslb := &kgbv1beta1.Gslb{}
	err = json.Unmarshal(jsonBytes, &gslb)
	if err != nil {
		return &kgbv1beta1.Gslb{}, err
	}
	return gslb, nil
}
