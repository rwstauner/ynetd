package service

import (
	"strings"
	"testing"
	"time"

	"github.com/rwstauner/ynetd/procman"
)

func TestNew(t *testing.T) {
	svc := New(Config{
		Proxy:   map[string]string{"hello": "goodbye"},
		Command: []string{"sleep", "10"},
		Timeout: "4s",
	}, procman.New())

	if svc.Proxy["hello"] != "goodbye" {
		t.Errorf("proxy incorrect: %s", svc.Proxy)
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
	func() {
		defer func() {
			if r := recover(); r != nil {
				if !strings.Contains(r.(error).Error(), "invalid duration") {
					t.Errorf("unexpected error: %s", r)
				}
			} else {
				t.Errorf("expected error, got none")
			}
		}()
		_ = New(Config{Timeout: "foo"}, procman.New())
	}()
}
