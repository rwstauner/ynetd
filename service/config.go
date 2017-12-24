package service

import (
	"time"

	"github.com/rwstauner/ynetd/config"
	"github.com/rwstauner/ynetd/procman"
)

func parseDuration(str string, defaultVal time.Duration) time.Duration {
	if str == "" {
		return defaultVal
	}
	duration, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}
	return duration
}

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
		Timeout:   parseDuration(c.Timeout, config.DefaultTimeout),
		StopAfter: parseDuration(c.StopAfter, config.DefaultStopAfter),
	}
	return
}
