// +build windows

package procman

import (
	"os"
	"os/exec"
)

func (proc *Process) signal(sig os.Signal) error {
	return proc.cmd.Process.Signal(sig)
}

func prepareCommand(cmd *exec.Cmd) {
	// noop
}
