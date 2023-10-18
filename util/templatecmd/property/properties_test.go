package property

import (
	"testing"
)

func TestPropertiesSet(t *testing.T) {
	properties := make(Properties)
	properties["a"] = "1"
	var options Options
	options.KeyValues = []string{"b=2", "c.d=3"}
	options.Apply(properties)
	if properties["a"] != "1" {
		t.Fail()
	}
	if properties["b"] != "2" {
		t.Fail()
	}
	c := properties["c"].(map[any]any)
	if c["d"] != "3" {
		t.Fail()
	}
}
