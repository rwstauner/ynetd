package procman

import (
	"fmt"
	"testing"
)

func TestProcess(t *testing.T) {
	pm := New()
	proc := pm.Process([]string{"foo", "bar"})

	if fmt.Sprintf("%s", proc.argv) != "[foo bar]" {
		t.Errorf("Unexpected argv: %s", proc.argv)
	}
	if proc.manager != pm {
		t.Errorf("who dis?")
	}
}

func TestProcessEmpty(t *testing.T) {
	pm := New()
	proc := pm.Process([]string{})

	if proc != nil {
		t.Errorf("expected nil, got %s", proc)
	}
}
