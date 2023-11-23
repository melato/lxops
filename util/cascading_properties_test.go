package util

import (
	"testing"
)

func TestCascadingProperties(t *testing.T) {
	var properties CascadingProperties
	properties.AddMap(map[string]string{
		"a": "A",
	})
	s, err := properties.Substitute("a(a)b(a)c")
	if err != nil {
		t.Fail()
	}
	if s != "aAbAc" {
		t.Fail()
	}
}
