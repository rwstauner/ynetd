// +build !windows

package procman

import (
	"os"
	"os/exec"
	"syscall"
)

func (proc *Process) signal(sig os.Signal) error {
	return syscall.Kill(-proc.cmd.Process.Pid, sig.(syscall.Signal))
}

func prepareCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
