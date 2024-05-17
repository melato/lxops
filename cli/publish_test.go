package cli

import (
	"fmt"
	"testing"

	"melato.org/lxops/srv"
)

func TestCapFirst(t *testing.T) {
	verify := func(s, expected string) {
		c := capFirst(s)
		if c != expected {
			t.Fail()
			fmt.Printf("capFirst(%s)=%s, expected: %s\n", s, c, expected)
		}
	}
	verify("", "")
	verify("a", "A")
	verify("one", "One")
	verify("αβ", "Αβ")
}

func TestMergeStructs(t *testing.T) {
	source := srv.ImageFields{
		Name:    "name1",
		Variant: "var1",
	}
	target := srv.ImageFields{
		Release: "release2",
		Variant: "var2",
	}
	var p PublishOps
	p.mergeStructs(&target, &source)
	expect := func(expected, actual string) {
		if actual != expected {
			t.Fatalf("expected: %s, actual: %s", expected, actual)
		}
	}
	expect("name1", target.Name)
	expect("var1", target.Variant)
	expect("release2", target.Release)
}
