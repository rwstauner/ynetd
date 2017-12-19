package service

import (
	"time"

	"github.com/rwstauner/ynetd/procman"
)

// DefaultTimeout is the default timeout duration for new connections
// to proxy to the service.
var DefaultTimeout = 5 * time.Minute

// Config holds string representations of Service attributes.
type Config struct {
	Proxy   map[string]string
	Command []string
	Timeout string
}

func parseTimeout(timeout string) time.Duration {
	if timeout == "" {
		return DefaultTimeout
	}
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		panic(err)
	}
	return duration
}

// New returns a new Service based on the provided Config.
func New(c Config, pm *procman.ProcessManager) Service {
	return Service{
		Proxy:   c.Proxy,
		Command: pm.Process(c.Command),
		Timeout: parseTimeout(c.Timeout),
	}
}
