package network

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
)

type Address struct {
	Ip   string
	Mask int
}

var cidr4Pattern = regexp.MustCompile("^([0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+)/([0-9]+)")
var cidr6Pattern = regexp.MustCompile("^([0-9a-fA-F:]+)/([0-9]+)")

func DetectIpv4Addresses(public bool) ([]Ipv4, error) {
	addresses, err := FindIpv4Addresses()
	if err != nil {
		return nil, err
	}
	var result []Ipv4
	for _, a := range addresses {
		ip, err := ParseIpv4(a.Ip)
		if err != nil {
			return nil, err
		}
		if ip.IsPrivate() != public {
			result = append(result, ip)
		}
	}
	return result, nil
}

func FindIpv4Addresses() ([]Address, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var addresses []Address
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range addrs {
			s := a.String()
			parts := cidr4Pattern.FindStringSubmatch(s)
			if len(parts) == 3 {
				var addr Address
				addr.Ip = parts[1]
				addr.Mask, err = strconv.Atoi(parts[2])
				if err != nil {
					return nil, fmt.Errorf("invalid network mask: %s", s)
				}
				if addr.Ip != "127.0.0.1" {
					addresses = append(addresses, addr)
				}
			}
		}
	}
	return addresses, nil
}

func DetectIpv6Addresses(public bool) ([]Ipv6, error) {
	addresses, err := FindIpv6Addresses()
	if err != nil {
		return nil, err
	}
	var result []Ipv6
	for _, a := range addresses {
		ip, err := ParseIpv6(a.Ip)
		if err != nil {
			return nil, err
		}
		if (public && ip.IsPublic()) || (!public && ip.IsLocalUnicast()) {
			result = append(result, ip)
		}
	}
	return result, nil
}

func FindIpv6Addresses() ([]Address, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var addresses []Address
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range addrs {
			s := a.String()
			parts := cidr6Pattern.FindStringSubmatch(s)
			if len(parts) == 3 {
				var addr Address
				addr.Ip = parts[1]
				addr.Mask, err = strconv.Atoi(parts[2])
				if err != nil {
					return nil, fmt.Errorf("invalid network mask: %s", s)
				}
				if addr.Ip != "::1" {
					addresses = append(addresses, addr)
				}
			}
		}
	}
	return addresses, nil
}
