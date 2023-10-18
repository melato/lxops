package network

import (
	"testing"
)

func TestIpv4String(t *testing.T) {
	ip, err := ParseIpv4("10.0.0.1")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if ip.String() != "10.0.0.1" {
		t.Fail()
	}
}

func TestIpv4Next(t *testing.T) {
	ip, err := ParseIpv4("10.0.0.1")
	if err != nil {
		t.Fatalf("%v", err)
	}
	ip.Next()
	if ip.String() != "10.0.0.2" {
		t.Fail()
	}
}
