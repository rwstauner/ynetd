package procman

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

// Process represents a command managed by a ProcessManager.
type Process struct {
	argv       []string
	cmd        *exec.Cmd
	manager    *ProcessManager
	mutex      *sync.Mutex
	stopSignal os.Signal
}

// LaunchOnce launches the process if it isn't already running.
func (p *Process) LaunchOnce() {
	p.manager.launcher <- p
}

// Stop sends the configured signal to the process.
func (p *Process) Stop() {
	p.manager.stopper <- p
}

func (p *Process) String() string {
	return fmt.Sprintf("%s (running: %t)", p.argv, p.cmd != nil)
}
