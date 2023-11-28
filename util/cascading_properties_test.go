package util

import (
	"testing"
)

func TestCascadingProperties(t *testing.T) {
	var properties CascadingProperties
	//properties.AddMap(nil)
	properties.AddMap(map[string]string{
		"a": "A",
	})
	s, err := properties.Substitute("a(a)b(a)c")
	if err != nil {
		t.Fatalf("substitution error: %v", err)
	}
	if s != "aAbAc" {
		t.Fatalf("wrong value: %s", s)
	}
	_, err = properties.Substitute("(x)")
	if err == nil {
		t.Fatalf("should have failed")
	}
}
