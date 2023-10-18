package lxops

import (
	"testing"
)

func TestIsDescendant(t *testing.T) {
	if !Path("a/b/c").IsDescendantOf("a/b") {
		t.Fail()
	}
}

func TestPathRoots(t *testing.T) {
	roots := RootPaths([]string{
		"a/b",
		"a/b/c",
		"a/d",
	})
	if len(roots) != 2 {
		t.Fail()
	}
	if roots[0] != "a/b" || roots[1] != "a/d" {
		t.Fail()
	}
}
