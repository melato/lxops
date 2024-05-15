package cli

import (
	"testing"
)

func TestParsePrefixNameDateTime(t *testing.T) {
	var opt ImageMetadataOptions
	err := opt.ParsePrefixNameDateTime("a-nginx-20240203-0834")
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	expect := func(name, expected, value string) {
		if value != expected {
			t.Fatalf("%s=%s, expected: %s", name, value, expected)
		}
	}
	expect("os", "nginx", opt.OS)
	expect("variant", "a", opt.Variant)
	expect("date", "20240203-0834", opt.Date.UTC().Format("20060102-1504"))
}
