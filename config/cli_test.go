package config

import (
	"fmt"
	"testing"
	"time"
)

func TestLoadNoArgs(t *testing.T) {
	listenAddress = ""
	proxyAddress = ""
	cfg, err := Load([]string{})

	if len(cfg.Services) != 0 {
		t.Errorf("Service configured without args")
	}

	if err != nil {
		t.Errorf("got error: %s", err)
	}
}

func TestLoadBasicArgs(t *testing.T) {
	listenAddress = ":5000"
	proxyAddress = "localhost:5001"
	timeout = 2 * time.Second
	cfg, err := Load([]string{"foo", "bar"})

	if err != nil {
		t.Errorf("got error: %s", err)
	}

	if len(cfg.Services) != 1 {
		t.Errorf("Service not configured from args")
	}

	svc := cfg.Services[0]
	if len(svc.Proxy) != 1 {
		t.Errorf("Proxy incorrect")
	}
	if svc.Proxy[":5000"] != "localhost:5001" {
		t.Errorf("Proxy incorrect")
	}
	if fmt.Sprintf("%s", svc.Command) != "[foo bar]" {
		t.Errorf("Command incorrect")
	}
	if svc.Timeout != "2s" {
		t.Errorf("Timeout incorrect")
	}
}

func TestLoadNoListen(t *testing.T) {
	listenAddress = ""
	proxyAddress = "localhost:5001"
	_, err := Load([]string{"foo", "bar"})

	if err == nil {
		t.Errorf("expected error, got none")
	} else if err.Error() != "listenAddress is required" {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestLoadNoProxy(t *testing.T) {
	listenAddress = ":5000"
	proxyAddress = ""
	_, err := Load([]string{"foo", "bar"})

	if err == nil {
		t.Errorf("expected error, got none")
	} else if err.Error() != "proxyAddress is required" {
		t.Errorf("unexpected error: %s", err)
	}
}
