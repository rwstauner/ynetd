package procman

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Process represents a command managed by a ProcessManager.
type Process struct {
	argv           []string
	cmd            *exec.Cmd
	manager        *ProcessManager
	started        time.Time
	stopSignal     os.Signal
	waitAfterStart time.Duration
	waitUntil      time.Time
}

type launchRequest struct {
	process *Process
	ready   chan bool
}

// LaunchOnce launches the process if it isn't already running.
func (p *Process) LaunchOnce() {
	ready := make(chan bool)
	req := &launchRequest{process: p, ready: ready}
	p.manager.launcher <- req
	<-ready
}

// Stop sends the configured signal to the process.
func (p *Process) Stop() {
	p.manager.stopper <- p
}

func (p *Process) String() string {
	return fmt.Sprintf("%s (running: %t)", p.argv, p.cmd != nil)
}
