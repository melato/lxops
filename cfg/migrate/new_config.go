package migrate

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v2"
)

const Comment = "#new"

func MigrateNew(data []byte) ([]byte, error) {
	var m map[string]any
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	os := make(map[string]any)
	ostype := m["ostype"]
	image := m["image"]
	if ostype != nil {
		os["name"] = ostype
	}
	if image != nil {
		os["image"] = image
	}
	delete(m, "ostype")
	delete(m, "image")
	m["os"] = os

	var buf bytes.Buffer
	fmt.Fprintln(&buf, "#lxdops\n")
	encoder := yaml.NewEncoder(&buf)
	err = encoder.Encode(m)
	if err != nil {
		return nil, err
	}
	err = encoder.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
