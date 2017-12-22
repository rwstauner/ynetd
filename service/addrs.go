package service

import (
	"fmt"
	"net"
	"strings"
)

func ipString(ip net.IP, zone string) string {
	if ip.To4() == nil {
		return fmt.Sprintf("%s%%%s", ip.String(), zone)
	}
	return ip.String()
}

func interfaceAddrs(name string) ([]string, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	ips := []string{}
	for _, a := range addrs {
		ip, _, err := net.ParseCIDR(a.String())
		if err != nil {
			return nil, err
		}
		ips = append(ips, ipString(ip, name))
	}
	return ips, nil
}

func parseAddr(addr string) ([]string, error) {
	spec := strings.Split(addr, ":")

	if len(spec) == 3 && spec[0] == "interface" {
		name, port := spec[1], spec[2]
		ips, err := interfaceAddrs(name)
		if err != nil {
			return nil, err
		}
		addrs := []string{}
		for _, ip := range ips {
			addrs = append(addrs, net.JoinHostPort(ip, port))
		}
		if len(addrs) == 0 {
			return nil, fmt.Errorf("no addresses found for interface: %s", name)
		}
		return addrs, nil
	}

	return []string{addr}, nil
}
