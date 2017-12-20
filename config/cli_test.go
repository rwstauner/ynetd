package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

func TestLoadConfigFile(t *testing.T) {
	listenAddress = ""
	proxyAddress = ""

	tmpfile, err := ioutil.TempFile("", "ynetdjson")
	if err != nil {
		t.Errorf("failed to create tempfile: %s", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(`{
		"Services": [
			{
				"Proxy": {":5000": "localhost:5001"},
				"Command": ["3", "4"],
				"Timeout": "15ms"
			}
		]
	}`))
	tmpfile.Close()

	configfile = tmpfile.Name()

	cfg, err := Load([]string{})

	if err != nil {
		t.Errorf("got error: %s", err)
	}

	if len(cfg.Services) != 1 {
		t.Errorf("services incorrect: %d", len(cfg.Services))
	}

	svc := cfg.Services[0]
	if len(svc.Proxy) != 1 {
		t.Errorf("Proxy incorrect")
	}
	if svc.Proxy[":5000"] != "localhost:5001" {
		t.Errorf("Proxy incorrect")
	}
	if fmt.Sprintf("%s", svc.Command) != "[3 4]" {
		t.Errorf("Command incorrect")
	}
	if svc.Timeout != "15ms" {
		t.Errorf("Timeout incorrect")
	}
}

func TestLoadConfigFileError(t *testing.T) {
	listenAddress = ""
	proxyAddress = ""

	tmpfile, err := ioutil.TempFile("", "ynetdjson")
	if err != nil {
		t.Errorf("failed to create tempfile: %s", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(`{ "Services": [], }`))
	tmpfile.Close()

	configfile = tmpfile.Name()

	cfg, err := Load([]string{})

	if err == nil {
		t.Errorf("expected error, got none")
	} else if !strings.Contains(err.Error(), fmt.Sprintf("parsing config file '%s': invalid char", tmpfile.Name())) {
		t.Errorf("unexpected error: %s", err)
	}

	if len(cfg.Services) != 0 {
		t.Errorf("got %d services", len(cfg.Services))
	}
}
