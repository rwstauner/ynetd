package procman

import (
	"os/exec"
	"testing"
)

func TestString(t *testing.T) {
	proc := &Process{argv: []string{"sleep", "1"}}

	if proc.String() != "[sleep 1] (running: false)" {
		t.Errorf("unexpected string: %s", proc)
	}

	// fake it
	proc.cmd = &exec.Cmd{}
	if proc.String() != "[sleep 1] (running: true)" {
		t.Errorf("unexpected string: %s", proc)
	}
}
