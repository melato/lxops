package network

import (
	"fmt"
	"strconv"
	"strings"
)

type Ipv6 [8]uint16

func (dd Ipv6) String() string {
	ss := make([]string, len(dd))
	for i, d := range dd {
		ss[i] = strconv.FormatUint(uint64(d), 16)
	}
	return strings.Join(ss, ":")
}

func (ip Ipv6) IsLocalUnicast() bool {
	// fc00::/7 - Unique Local Address
	// 0xfc = 11111100
	// 0xfd = 11111101
	// 0xfe = 11111110
	return ip[0]&0xfe00 == 0xfc00
}

func (ip Ipv6) IsLinkLocal() bool {
	return ip[0] == 0xfe80
}

func (ip Ipv6) IsPublic() bool {
	return !ip.IsLocalUnicast() && !ip.IsLinkLocal()
}

func (ip Ipv6) WithLast(d uint16) Ipv6 {
	ip[len(ip)-1] = d
	return ip
}

func (ip Ipv6) SetLast(d uint16) Ipv6 {
	return ip.WithLast(d)
}

func (ip Ipv6) Next() Ipv6 {
	for i := len(ip) - 1; i >= 0; i-- {
		d := ip[i]
		if d < 0xffff {
			ip[i] = d + 1
			break
		}
		ip[i] = 1
	}
	return ip
}

func ParseIpv6(s string) (Ipv6, error) {
	var ip Ipv6
	parts := strings.Split(s, ":")
	n := len(parts)
	if n < 3 || n > 8 {
		return Ipv6{}, fmt.Errorf("cannot parse Ipv6 address: %s", s)
	}
	var doubleColon bool
	var j int
	for _, p := range parts {
		if p != "" {
			u, err := strconv.ParseUint(p, 16, 32)
			if err != nil {
				return Ipv6{}, fmt.Errorf("cannot parse Ipv6 address: %s", s)
			}
			ip[j] = uint16(u)
			j++
		} else {
			if j == 0 || j == 7 {
				ip[j] = 0
				break
			}
			if doubleColon {
				return Ipv6{}, fmt.Errorf("cannot parse Ipv6 address: %s", s)
			}
			doubleColon = true
			for k := n; k <= 8; k++ {
				ip[j] = 0
				j++
			}
		}
	}
	return ip, nil
}

func IsPrivateIpv6(s string) (bool, error) {
	ip, err := ParseIpv6(s)
	if err != nil {
		return false, err
	}
	return ip.IsLocalUnicast(), nil
}
