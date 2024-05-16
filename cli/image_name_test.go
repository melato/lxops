package cli

import (
	"testing"
)

func TestPrefixNameParse(t *testing.T) {
	var p PrefixNameTime
	err := p.Parse("a-nginx-20240203-0834")
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	expect := func(name, expected, value string) {
		if value != expected {
			t.Fatalf("%s=%s, expected: %s", name, value, expected)
		}
	}
	expect("prefix", "a", p.Prefix)
	expect("base", "nginx", p.Name)
	expect("date", "20240203-0834", p.Timestamp)
}
