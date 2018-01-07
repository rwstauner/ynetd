package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestParseConfigJson(t *testing.T) {
	t.Run("minimal json", func(t *testing.T) {
		bytes := []byte(`{
			"Services": [
				{
					"Proxy": {":3000": "localhost:4000"},
					"Command": ["sleep", "1"]
				},
				{
					"Command": ["sleep", "2"],
					"Proxy": {":3001": "localhost:4001"}
				}
			]
		}`)
		assertValidJSON(t, bytes)
		assertFileParsed(t, bytes, false)
	})

	t.Run("full json", func(t *testing.T) {
		bytes := []byte(`{
			"Services": [
				{
					"Proxy": {":3000": "localhost:4000"},
					"Command": ["sleep", "1"],
					"StopAfter": "10m",
					"StopSignal": "INT",
					"Timeout": "150ms",
					"WaitAfterStart": "250ms"
				},
				{
					"Command": ["sleep", "2"],
					"Proxy": {":3001": "localhost:4001"},
					"Timeout": "151ms",
					"StopAfter": "11m",
					"StopSignal": "TERM",
					"WaitAfterStart": "251ms"
				}
			]
		}`)
		assertValidJSON(t, bytes)
		assertFileParsed(t, bytes, true)
	})
}

func assertValidJSON(t *testing.T, bytes []byte) {
	var obj Config
	err := json.Unmarshal(bytes, &obj)
	if err != nil {
		t.Errorf("invalid json: %s", err)
	}
}

func TestParseConfigYaml(t *testing.T) {
	t.Run("valid json", func(t *testing.T) {
		bytes := []byte(`{
"services": [
	{ "proxy": {":3000": "localhost:4000"}, "command": ["sleep", "1"] },
	{
		"command": ["sleep", "2"],
		"proxy": {":3001": "localhost:4001"}
	}
]}`)
		assertValidJSON(t, bytes)
		assertFileParsed(t, bytes, false)
	})
	t.Run("minimal yaml", func(t *testing.T) {
		assertFileParsed(t, []byte(`---
services:
  - { proxy: {":3000": "localhost:4000"}, command: ["sleep", "1"] }
  -
    command:
      - sleep
      - "2"
    proxy:
      ":3001": "localhost:4001"
`),
			false)
	})

	t.Run("full yaml", func(t *testing.T) {
		assertFileParsed(t, []byte(`---
services:
  -
    proxy: {":3000": "localhost:4000"}
    command: ["sleep", "1"]
    stop_after: "10m"
    stop_signal: "INT"
    timeout: "150ms"
    wait_after_start: "250ms"
  -
    command:
      - sleep
      - 2
    proxy:
      ":3001": "localhost:4001"
    stop_after: 11m
    stop_signal: TERM
    timeout: 151ms
    wait_after_start: 251ms
`),
			true)
	})
}

func assertFileParsed(t *testing.T, bytes []byte, everything bool) {
	tmpfile, err := ioutil.TempFile("", "ynetd")
	if err != nil {
		t.Errorf("failed to create tempfile: %s", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write(bytes)
	tmpfile.Close()

	cfg, err := parseConfigFile(tmpfile.Name())

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(cfg.Services) != 2 {
		t.Errorf("got wrong number of services: %d", len(cfg.Services))
	}

	for i, svc := range cfg.Services {
		if len(svc.Proxy) != 1 || svc.Proxy[fmt.Sprintf(":%d", i+3000)] != fmt.Sprintf("localhost:%d", i+4000) {
			t.Errorf("unexpected Proxy: %v", svc.Proxy)
		}
		if fmt.Sprintf("%s", svc.Command) != fmt.Sprintf("[sleep %d]", i+1) {
			t.Errorf("unexpected Command: %s", svc.Command)
		}
		var (
			expTimeout        = ""
			expStopAfter      = ""
			expStopSignal     = ""
			expWaitAfterStart = ""
		)
		if everything {
			expTimeout = fmt.Sprintf("%dms", i+150)
			expStopAfter = fmt.Sprintf("%dm", i+10)
			expStopSignal = []string{"INT", "TERM"}[i]
			expWaitAfterStart = fmt.Sprintf("%dms", i+250)
		}
		if svc.Timeout != expTimeout {
			t.Errorf("unexpected Timeout: %s", svc.Timeout)
		}
		if svc.StopAfter != expStopAfter {
			t.Errorf("unexpected StopAfter: %s", svc.StopAfter)
		}
		if svc.StopSignal != expStopSignal {
			t.Errorf("unexpected StopSignal: %s", svc.StopSignal)
		}
		if svc.WaitAfterStart != expWaitAfterStart {
			t.Errorf("unexpected WaitAfterStart: %s", svc.WaitAfterStart)
		}
	}
}

func TestParseConfigFileError(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "ynetd")
	if err != nil {
		t.Errorf("failed to create tempfile: %s", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(`{`))
	tmpfile.Close()

	cfg, err := parseConfigFile(tmpfile.Name())

	if err == nil {
		t.Errorf("expected error, got none")
	} else if !strings.Contains(err.Error(), "yaml: ") {
		t.Errorf("unexpected error: %s", err)
	}

	if len(cfg.Services) > 0 {
		t.Errorf("got services: %d", len(cfg.Services))
	}
}
