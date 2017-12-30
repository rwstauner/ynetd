package service

import (
	"github.com/rwstauner/ynetd/config"
	"github.com/rwstauner/ynetd/procman"
)

// New returns the address to a new Service based on the provided Config.
func New(c config.Service, pm *procman.ProcessManager) (svc *Service, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	svc = &Service{
		Proxy:     c.Proxy,
		Command:   pm.Process(c),
		Timeout:   config.ParseDuration(c.Timeout, config.DefaultTimeout),
		StopAfter: config.ParseDuration(c.StopAfter, config.DefaultStopAfter),
	}
	return
}
