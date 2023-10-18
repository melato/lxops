package lxops

import (
	"testing"
)

func TestBaseName(t *testing.T) {
	name := BaseName("a.yaml")
	if name != "a" {
		t.Errorf("name: %s expected:%s", name, "a")
	}
}
