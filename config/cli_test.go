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
	proxySpec = ""
	cfg, err := Load([]string{})

	if len(cfg.Services) != 0 {
		t.Errorf("Service configured without args")
	}

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestLoadArgs(t *testing.T) {
	proxySpec = ":5000 localhost:5001 some:6001 some:7001"
	timeout = 2 * time.Second
	waitAfterStart = 500 * time.Millisecond
	cfg, err := Load([]string{"foo", "bar"})

	if err != nil {
		t.Errorf("got error: %s", err)
	}

	if len(cfg.Services) != 1 {
		t.Fatalf("Service not configured from args")
	}

	svc := cfg.Services[0]
	if len(svc.Proxy) != 2 || svc.Proxy[":5000"] != "localhost:5001" || svc.Proxy["some:6001"] != "some:7001" {
		t.Errorf("Proxy incorrect: %q", svc.Proxy)
	}
	if fmt.Sprintf("%s", svc.Command) != "[foo bar]" {
		t.Errorf("Command incorrect: %s", svc.Command)
	}
	if svc.Timeout != "2s" {
		t.Errorf("Timeout incorrect: %s", svc.Timeout)
	}
	if svc.WaitAfterStart != "500ms" {
		t.Errorf("WaitAfterStart incorrect: %s", svc.WaitAfterStart)
	}
}

func TestLoadProxySep(t *testing.T) {
	proxySep = "+"
	proxySpec = ":5000+localhost:5001+some:6001+some:7001"
	cfg, err := Load([]string{})

	if err != nil {
		t.Errorf("got error: %s", err)
	}

	if len(cfg.Services) != 1 {
		t.Fatalf("Service not configured from args")
	}

	svc := cfg.Services[0]
	if len(svc.Proxy) != 2 || svc.Proxy[":5000"] != "localhost:5001" || svc.Proxy["some:6001"] != "some:7001" {
		t.Errorf("Proxy incorrect: %q", svc.Proxy)
	}
}

func TestLoadOddProxy(t *testing.T) {
	tests := []string{"localhost:5001", ":5001 :5002 :5003"}

	for _, val := range tests {
		t.Run(fmt.Sprintf("-proxy '%s'", val), func(t *testing.T) {
			proxySpec = val
			_, err := Load([]string{"foo", "bar"})

			if err == nil {
				t.Errorf("expected error, got none")
			} else if !strings.Contains(err.Error(), "-proxy must contain pairs") {
				t.Errorf("unexpected error: %s", err)
			}
		})
	}
}

func TestLoadNoProxy(t *testing.T) {
	proxySpec = ""
	_, err := Load([]string{"foo", "bar"})

	if err == nil {
		t.Errorf("expected error, got none")
	} else if err.Error() != "-proxy is required" {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestLoadConfigFile(t *testing.T) {
	proxySpec = ""

	tmpfile, err := ioutil.TempFile("", "ynetdjson")
	if err != nil {
		t.Errorf("failed to create tempfile: %s", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(`
"services":
  - proxy: {":5000": "localhost:5001"}
    command: ["3", "4"]
    timeout: "15ms"
    wait_after_start: "25ms"
`))
	tmpfile.Close()

	configfile = tmpfile.Name()

	cfg, err := Load([]string{})

	if err != nil {
		t.Errorf("got error: %s", err)
	}

	if len(cfg.Services) != 1 {
		t.Fatalf("services incorrect: %d", len(cfg.Services))
	}

	svc := cfg.Services[0]
	if len(svc.Proxy) != 1 || svc.Proxy[":5000"] != "localhost:5001" {
		t.Errorf("Proxy incorrect: %q", svc.Proxy)
	}
	if fmt.Sprintf("%s", svc.Command) != "[3 4]" {
		t.Errorf("Command incorrect: %s", svc.Command)
	}
	if svc.Timeout != "15ms" {
		t.Errorf("Timeout incorrect: %s", svc.Timeout)
	}
	if svc.WaitAfterStart != "25ms" {
		t.Errorf("WaitAfterStart incorrect: %s", svc.WaitAfterStart)
	}
}

func TestLoadConfigFileError(t *testing.T) {
	proxySpec = ""

	tmpfile, err := ioutil.TempFile("", "ynetdjson")
	if err != nil {
		t.Errorf("failed to create tempfile: %s", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(`{`))
	tmpfile.Close()

	configfile = tmpfile.Name()

	cfg, err := Load([]string{})

	if err == nil {
		t.Errorf("expected error, got none")
	} else if !strings.Contains(err.Error(), fmt.Sprintf("parsing config file '%s': yaml: line 1: did not find expected node content", tmpfile.Name())) {
		t.Errorf("unexpected error: %s", err)
	}

	if len(cfg.Services) != 0 {
		t.Errorf("got %d services", len(cfg.Services))
	}
}

func TestLoadStopAfter(t *testing.T) {
	setup := func(t *testing.T) Service {
		configfile = ""
		proxySep = " "
		proxySpec = ":5000 localhost:5001"

		cfg, err := Load([]string{})

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if len(cfg.Services) != 1 {
			t.Fatalf("unexpected services: %q", cfg.Services)
		}

		return cfg.Services[0]
	}

	t.Run("default", func(t *testing.T) {
		stopAfter = 0
		stopSignal = "TERM"

		svc := setup(t)

		if svc.StopAfter != "0s" {
			t.Errorf("incorrect StopAfter: %q", svc.StopAfter)
		}
		if svc.StopSignal != "TERM" {
			t.Errorf("incorrect StopAfter: %q", svc.StopSignal)
		}
	})

	t.Run("custom", func(t *testing.T) {
		stopAfter = 200 * time.Millisecond
		stopSignal = "INT"

		svc := setup(t)

		if svc.StopAfter != "200ms" {
			t.Errorf("incorrect StopAfter: %q", svc.StopAfter)
		}
		if svc.StopSignal != "INT" {
			t.Errorf("incorrect StopAfter: %q", svc.StopSignal)
		}
	})
}
