package procman

import (
	"fmt"
	"syscall"
	"testing"

	"github.com/rwstauner/ynetd/config"
)

var stopSignal = syscall.SIGINT

func TestProcess(t *testing.T) {
	pm := New()
	proc := pm.Process(config.Service{Command: []string{"foo", "bar"}})

	if fmt.Sprintf("%s", proc.argv) != "[foo bar]" {
		t.Errorf("Unexpected argv: %s", proc.argv)
	}
	if proc.manager != pm {
		t.Errorf("who dis?")
	}
	if proc.stopSignal != syscall.SIGINT {
		t.Errorf("incorrect stop signal: %s", proc.stopSignal)
	}
}

func TestProcessStopSignal(t *testing.T) {
	pm := New()
	proc := pm.Process(config.Service{
		Command:    []string{"siggy"},
		StopSignal: "TERM",
	})

	if fmt.Sprintf("%s", proc.argv) != "[siggy]" {
		t.Errorf("Unexpected argv: %s", proc.argv)
	}
	if proc.stopSignal != syscall.SIGTERM {
		t.Errorf("incorrect stop signal: %s", proc.stopSignal)
	}
}

func TestProcessEmpty(t *testing.T) {
	pm := New()
	proc := pm.Process(config.Service{Command: []string{}})

	if proc != nil {
		t.Errorf("expected nil, got %s", proc)
	}
}
