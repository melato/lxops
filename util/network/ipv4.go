package network

import (
	"fmt"
	"strconv"
	"strings"
)

type Ipv4 []uint8

func (dd Ipv4) String() string {
	ss := make([]string, len(dd))
	for i, d := range dd {
		ss[i] = strconv.Itoa(int(d))
	}
	return strings.Join(ss, ".")
}

func (dd Ipv4) IsPrivate() bool {
	if dd[0] == 10 {
		return true
	}
	if dd[0] == 192 && dd[1] == 168 {
		return true
	}
	if dd[0] == 172 && (dd[1]&0xf0) == 16 {
		return true
	}
	return false
}

func (dd Ipv4) SetLast(d uint8) {
	dd[len(dd)-1] = d
}

func (dd Ipv4) WithLast(d uint8) Ipv4 {
	result := make(Ipv4, len(dd))
	copy(result, dd)
	result.SetLast(d)
	return result
}

func (dd Ipv4) Next() {
	for i := len(dd) - 1; i >= 0; i-- {
		d := dd[i]
		if d < 255 {
			dd[i] = d + 1
			return
		}
		dd[i] = 1
	}
}

func ParseIpv4(ip string) (Ipv4, error) {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return nil, fmt.Errorf("cannot parse ipv4 address: %s", ip)
	}
	dd := make([]uint8, 4)
	for i, _ := range dd {
		d, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, fmt.Errorf("cannot parse ipv4 address: %s", ip)
		}
		dd[i] = uint8(d)
	}
	return dd, nil
}

func IsPrivateIpv4(s string) (bool, error) {
	ip, err := ParseIpv4(s)
	if err != nil {
		return false, err
	}
	return ip.IsPrivate(), nil
}
