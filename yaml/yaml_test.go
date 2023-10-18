package yaml

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestMap(t *testing.T) {
	m := make(map[string]string)
	m["a"] = "a"
	data, err := yaml.Marshal(m)
	if err != nil {
		t.Fail()
		return
	}
	var v interface{}
	err = yaml.Unmarshal(data, &v)
	if err != nil {
		t.Fail()
		return
	}
	switch v.(type) {
	case map[interface{}]interface{}:
		return
	case map[string]interface{}:
		return
	default:
		t.Fatalf("%T", v)
	}
}

type A struct {
	B int `yaml:"b"`
}

func TestDecode(t *testing.T) {
	var a A
	s := `b: 3`
	decoder := yaml.NewDecoder(strings.NewReader(s))
	err := decoder.Decode(&a)
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	if a.B != 3 {
		t.Fail()
	}
}
