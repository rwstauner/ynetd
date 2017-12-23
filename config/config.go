package config

import (
	"time"
)

// Config is a struct representing the ynetd configuration.
type Config struct {
	Services []Service
}

// DefaultTimeout is the default timeout duration for new connections
// to proxy to the service.
var DefaultTimeout = 5 * time.Minute

