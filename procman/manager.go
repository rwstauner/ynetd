package procman

import (
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/hashicorp/go-reap"
	"github.com/rwstauner/ynetd/config"
)

// ProcessManager manages launching and reaping of processes.
type ProcessManager struct {
	procs    map[int]*Process // pid -> proc
	launcher chan *launchRequest
	stopper  chan *Process
	signals  chan os.Signal
}

// New returns a new ProcessManager.
func New() *ProcessManager {
	return &ProcessManager{
		launcher: make(chan *launchRequest),
		stopper:  make(chan *Process),
		procs:    make(map[int]*Process),
		signals:  make(chan os.Signal),
	}
}

// Process returns the address to a new Process
// wrapping the command with the ProcessManager embedded.
func (m *ProcessManager) Process(cfg config.Service) *Process {
	if len(cfg.Command) == 0 {
		return nil
	}
	return &Process{
		argv:           cfg.Command,
		manager:        m,
		stopSignal:     getSignal(cfg.StopSignal, syscall.SIGINT),
		waitAfterStart: config.ParseDuration(cfg.WaitAfterStart, config.DefaultWaitAfterStart),
	}
}

func (m *ProcessManager) launch(proc *Process) *exec.Cmd {
	cmd := exec.Command(proc.argv[0], proc.argv[1:]...)
	prepareCommand(cmd)

	logger.Printf("starting: %s", proc)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		logger.Printf("start error: %s", err)
		return nil
	}

	return cmd
}

// Manage starts listening for processes to launch and reap.
func (m *ProcessManager) Manage() {
	reaper := make(reap.PidCh, 5)
	errors := make(reap.ErrorCh, 1)
	done := make(chan struct{})

	if reap.IsSupported() {
		go reap.ReapChildren(reaper, errors, done, nil)
	} else {
		logger.Printf("Child reaping is not currently supported on this platform.")
	}

	for {
		select {
		case sig := <-m.signals:
			// Tell reaper to stop.
			close(done)
			// Signal them all first.
			for _, process := range m.procs {
				if err := process.signal(sig); err != nil {
					logger.Printf("signal error: %s", err)
				}
			}
			// Then wait for them.
			// TODO: timeout?
			for _, process := range m.procs {
				err := process.cmd.Wait()
				status := 0
				if err != nil {
					if frdErr, ok := err.(*exec.ExitError); ok {
						status = frdErr.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
					}
				}
				logger.Printf("exited %d (%s) %s", status, err, process.argv)
			}
			return

		case req := <-m.launcher:
			proc := req.process
			if proc.cmd == nil {
				proc.cmd = m.launch(proc)
				proc.started = time.Now()
				if proc.waitAfterStart > 0 {
					proc.waitUntil = proc.started.Add(proc.waitAfterStart)
				}
				m.procs[proc.cmd.Process.Pid] = proc
			}
			if proc.waitAfterStart > 0 && proc.waitUntil.After(time.Now()) {
				time.AfterFunc(proc.waitUntil.Sub(time.Now()), func() { req.ready <- true })
			} else {
				req.ready <- true
			}

		case proc := <-m.stopper:
			if proc.cmd != nil {
				logger.Printf("stopping %s", proc.argv)
				if err := proc.signal(proc.stopSignal); err != nil {
					logger.Printf("signal error: %s", err)
				}
				// The reaper will take it from here.
			}

		case pid := <-reaper:
			proc, ok := m.procs[pid]
			if ok {
				proc.cmd = nil
				delete(m.procs, pid)
			}

		case err := <-errors:
			logger.Printf("reap error: %s", err)
		}
	}
}

// Signal sends a signal to each managed process.
func (m *ProcessManager) Signal(sig os.Signal) {
	m.signals <- sig
}
