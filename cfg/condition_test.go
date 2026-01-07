package cfg

import (
	"testing"
)

func TestCondition(t *testing.T) {
	vars := map[string]string{"one": "1"}
	get := func(name string) (string, bool) {
		value, found := vars[name]
		return value, found
	}
	verify := func(path string, expectedPath string, expectedPass bool) {
		result, pass := filterPath(HostPath(path), get)
		if string(result) != expectedPath || pass != expectedPass {
			t.Fatalf("%s: %s, %v", path, result, pass)
		}
	}
	verify("x", "x", true)
	verify("one|x", "x", true)
	verify("a|x", "", false)
}
