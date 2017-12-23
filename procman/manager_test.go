package procman

import (
	"fmt"
	"testing"

	"github.com/rwstauner/ynetd/config"
)

func TestProcess(t *testing.T) {
	pm := New()
	proc := pm.Process(config.Service{Command: []string{"foo", "bar"}})

	if fmt.Sprintf("%s", proc.argv) != "[foo bar]" {
		t.Errorf("Unexpected argv: %s", proc.argv)
	}
	if proc.manager != pm {
		t.Errorf("who dis?")
	}
}

func TestProcessEmpty(t *testing.T) {
	pm := New()
	proc := pm.Process(config.Service{Command: []string{}})

	if proc != nil {
		t.Errorf("expected nil, got %s", proc)
	}
}
