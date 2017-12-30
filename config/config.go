package config

import (
	"time"
)

// Config is a struct representing the ynetd configuration.
type Config struct {
	Services []Service
}

// DefaultTimeout is the default duration to allow new connections
// to attempt to forward to the service.
var DefaultTimeout = 5 * time.Minute

// DefaultStopAfter is the default duration of inactivity after which
// a command will be stopped (zero means the command will not be stopped).
var DefaultStopAfter = time.Duration(0)

// DefaultWaitAfterStart is the default duration to wait after starting a
// command before forwarding connections (zero means as soon as the port is open).
var DefaultWaitAfterStart = time.Duration(0)

// ParseDuration is a wrapper around time.ParseDuration.
// It returns the provided default if the string is blank
// and panics if there is an error in parsing.
func ParseDuration(str string, defaultVal time.Duration) time.Duration {
	if str == "" {
		return defaultVal
	}
	duration, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}
	return duration
}
