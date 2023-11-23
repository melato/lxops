package util

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	re := regexp.MustCompile(`\([^()]*\)`)
	var err error
	fn := func(key string) string {
		key = key[1 : len(key)-1]
		fmt.Printf("key: %s", key)
		switch key {
		case "a":
			return "A"
		case "b":
			return "B"
		default:
			err = fmt.Errorf("missing key: %s", key)
			return ""
		}
	}
	s := re.ReplaceAllStringFunc("a(b)c", fn)
	if s != "aBc" {
		t.Fatalf("%s", s)
	}
	s = re.ReplaceAllStringFunc("a(x)c", fn)
	if err == nil {
		t.Fatalf("expected error")
	}
}
