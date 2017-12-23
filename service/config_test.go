package service

import (
	"strings"
	"testing"
	"time"

	"github.com/rwstauner/ynetd/config"
	"github.com/rwstauner/ynetd/procman"
)

func TestNew(t *testing.T) {
	svc, err := New(config.Service{
		Proxy:   map[string]string{"hello": "goodbye"},
		Command: []string{"sleep", "10"},
		Timeout: "4s",
	}, procman.New())

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(svc.Proxy) != 1 || svc.Proxy["hello"] != "goodbye" {
		t.Errorf("proxy incorrect: %q", svc.Proxy)
	}

	if svc.Command == nil {
		t.Errorf("unexpected nil Command")
	} else if svc.Command.String() != "[sleep 10] (running: false)" {
		t.Errorf("command incorrect: %s", svc.Command)
	}

	if svc.Timeout != (4 * time.Second) {
		t.Errorf("timeout incorrect: %s", svc.Timeout)
	}
}

func TestNewError(t *testing.T) {
	_, err := New(config.Service{Timeout: "foo"}, procman.New())

	if err == nil {
		t.Errorf("expected error, got none")
	}
	if !strings.Contains(err.Error(), "invalid duration") {
		t.Errorf("unexpected error: %s", err)
	}
}
