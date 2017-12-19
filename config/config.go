package config

import (
	"github.com/rwstauner/ynetd/procman"
	"github.com/rwstauner/ynetd/service"
)

// Config is a struct representing the ynetd configuration.
type Config struct {
	Services []service.Config
}

// MakeServices creates Service objects from Config objects.
func MakeServices(cfg Config, pm *procman.ProcessManager) (services []service.Service, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	if len(cfg.Services) > 0 {
		for _, svc := range cfg.Services {
			services = append(services, service.New(svc, pm))
		}
	}

	return
}
