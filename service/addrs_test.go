package service

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

func TestParseAddrNonInterface(t *testing.T) {
	addrs, err := parseAddr("host:port")

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(addrs) != 1 || addrs[0] != "host:port" {
		t.Errorf("unexpected result: %q", addrs)
	}
}

func findInterface(t *testing.T, names ...string) (string, bool) {
	for _, name := range names {
		iface, err := net.InterfaceByName(name)
		if err == nil {
			isv6 := false
			addrs, err := iface.Addrs()
			if err == nil {
				for _, a := range addrs {
					if strings.Contains(a.String(), ":") {
						isv6 = true
					}
				}
			}
			return name, isv6
		}
	}
	t.Skipf("interfaces '%s' not found", names)
	return "", false
}

func TestParseAddrInterface(t *testing.T) {
	iface, isv6 := findInterface(t, "lo", "lo0")
	addrs, err := parseAddr(fmt.Sprintf("interface:%s:5001", iface))

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(addrs) == 0 {
		t.Errorf("no addresses")
	}

	found4, found6 := false, false
	for _, addr := range addrs {
		if addr == "127.0.0.1:5001" {
			found4 = true
		}
		if addr == fmt.Sprintf("[::1%%%s]:5001", iface) {
			found6 = true
		}
	}
	if !found4 {
		t.Errorf("localhost address not found in %s!", iface)
	}
	if isv6 && !found6 {
		t.Errorf("localhost v6 address not found in %s!", iface)
	}
}
