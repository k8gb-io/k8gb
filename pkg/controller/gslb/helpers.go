package gslb

import (
	"encoding/json"

	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
	yamlConv "github.com/ghodss/yaml"
)

//YamlToGslb takes yaml and returns Gslb object
func YamlToGslb(yaml []byte) (*ohmyglbv1beta1.Gslb, error) {
	// yamlBytes contains a []byte of my yaml job spec
	// convert the yaml to json
	jsonBytes, err := yamlConv.YAMLToJSON(yaml)
	if err != nil {
		return &ohmyglbv1beta1.Gslb{}, err
	}
	// unmarshal the json into the kube struct
	var gslb = &ohmyglbv1beta1.Gslb{}
	err = json.Unmarshal(jsonBytes, &gslb)
	if err != nil {
		return &ohmyglbv1beta1.Gslb{}, err
	}
	return gslb, nil
}
