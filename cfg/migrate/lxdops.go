package migrate

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

func MigrateLxdops(data []byte) ([]byte, error) {
	var m map[string]any
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	os := m["os"]
	if os != nil {
		v, ok := os.(map[any]any)
		if !ok {
			return nil, fmt.Errorf("unexpected os type: %T", os)
		}
		delete(m, "os")
		m["ostype"] = v["name"]
		m["image"] = v["image"]
	}
	return Marshal("lxops-v1", m)
}
