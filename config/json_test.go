package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestParseConfigFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "ynetdjson")
	if err != nil {
		t.Errorf("failed to create tempfile: %s", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(`{
		"Services": [
			{
				"Proxy": {":3000": "localhost:4000"},
				"Command": ["sleep", "1"],
				"StopAfter": "10m",
				"StopSignal": "INT",
				"Timeout": "150ms"
			},
			{
				"Proxy": {":3001": "localhost:4001"},
				"Command": ["sleep", "2"],
				"StopAfter": "11m",
				"StopSignal": "TERM",
				"Timeout": "151ms"
			}
		]
	}`))
	tmpfile.Close()

	cfg, err := parseConfigFile(tmpfile.Name())

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(cfg.Services) != 2 {
		t.Errorf("got wrong number of services: %d", len(cfg.Services))
	}

	for i, svc := range cfg.Services {
		if len(svc.Proxy) != 1 {
			t.Errorf("unexpected proxy: %v", svc.Proxy)
		}
		if svc.Proxy[fmt.Sprintf(":%d", i+3000)] != fmt.Sprintf("localhost:%d", i+4000) {
			t.Errorf("unexpected proxy: %v", svc.Proxy)
		}
		if fmt.Sprintf("%s", svc.Command) != fmt.Sprintf("[sleep %d]", i+1) {
			t.Errorf("unexpected command: %s", svc.Command)
		}
		if svc.Timeout != fmt.Sprintf("%dms", i+150) {
			t.Errorf("unexpected timeout: %s", svc.Timeout)
		}
		if svc.StopAfter != fmt.Sprintf("%dm", i+10) {
			t.Errorf("unexpected StopAfter: %s", svc.StopAfter)
		}
		if svc.StopSignal != []string{"INT", "TERM"}[i] {
			t.Errorf("unexpected StopSignal: %s", svc.StopSignal)
		}
	}
}

func TestParseConfigFileError(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "ynetdjson")
	if err != nil {
		t.Errorf("failed to create tempfile: %s", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(`{`))
	tmpfile.Close()

	cfg, err := parseConfigFile(tmpfile.Name())

	if err == nil {
		t.Errorf("expected error, got none")
	} else if !strings.Contains(err.Error(), "unexpected EOF") {
		t.Errorf("unexpected error: %s", err)
	}

	if len(cfg.Services) > 0 {
		t.Errorf("got services: %d", len(cfg.Services))
	}
}
