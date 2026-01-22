package util

import (
	"testing"
)

func TestEqualArrays(t *testing.T) {
	a := []string{"a", "b"}
	if !StringSlice(a).Equals(a) {
		t.Fail()
	}
}

func TestStringSliceDiff(t *testing.T) {
	a := []string{"a", "b"}
	b := []string{"b", "c"}
	c := StringSlice(a).Diff(b)
	if !StringSlice(c).Equals([]string{"a"}) {
		t.Fatalf("union: %v", c)
	}
}

func TestStringSliceUnion(t *testing.T) {
	a := []string{"a", "b"}
	b := []string{"b", "c"}
	c := StringSlice(a).Union(b)
	if !StringSlice(c).Equals([]string{"a", "b", "c"}) {
		t.Fatalf("union: %v", c)
	}
}
