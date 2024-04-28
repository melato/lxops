package cli

import (
	"fmt"
	"testing"
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
