package network

import (
	"testing"
)

func TestIpv6String(t *testing.T) {
	ip, err := ParseIpv6("2001:DB8:3:4:5:6:7:8")
	if err != nil {
		t.Fatalf("%v", err)
	}
	s := ip.String()
	if s != "2001:db8:3:4:5:6:7:8" {
		t.Fatalf("%s", s)
	}
}

func TestIpv6Next(t *testing.T) {
	ip, err := ParseIpv6("2001:DB8:3:4:5:6:7:8")
	if err != nil {
		t.Fatalf("%v", err)
	}
	ip = ip.Next()
	if ip.String() != "2001:db8:3:4:5:6:7:9" {
		t.Fail()
	}
}

func TestIpv6DoubleColon(t *testing.T) {
	ip, err := ParseIpv6("2001:DB8::8")
	if err != nil {
		t.Fatalf("%v", err)
	}
	s := ip.String()
	if s != "2001:db8:0:0:0:0:0:8" {
		t.Fatalf("%s", s)
	}
}

func TestIpv6Public(t *testing.T) {
	ip, err := ParseIpv6("2001:DB8::8")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if !ip.IsPublic() {
		t.Fail()
	}
}

func TestIpv6LocalUnicast(t *testing.T) {
	for _, s := range []string{"fc00::1", "fd00::1"} {
		ip, err := ParseIpv6(s)
		if err != nil {
			t.Fatalf("%v", err)
		}
		if !ip.IsLocalUnicast() {
			t.Fatalf("%v should be local unicast\n", ip)
		}
	}
}

func TestIpv6Last(t *testing.T) {
	ip, err := ParseIpv6("100::1")
	if err != nil {
		t.Fatalf("%v", err)
	}
	ip = ip.WithLast(7)
	if ip.String() != "100:0:0:0:0:0:0:7" {
		t.Fatalf("%s", ip.String())
	}
}
