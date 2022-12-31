package utils

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func YamlEncode(v interface{}) (*bytes.Buffer, error) {

	b := &bytes.Buffer{}
	yamlEncoder := yaml.NewEncoder(b)
	yamlEncoder.SetIndent(2)

	if err := yamlEncoder.Encode(v); err != nil {
		return nil, err
	}

	return b, nil
}
