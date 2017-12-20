package config

import (
	"flag"
	"fmt"

	"github.com/rwstauner/ynetd/service"
)

var configfile string
var listenAddress string
var proxyAddress string
var timeout = service.DefaultTimeout

func init() {
	const (
		configUsage  = "Path to configuration file"
		listenUsage  = "Address to listen on"
		proxyUsage   = "Address to proxy to (the address the command should be listening on)"
		timeoutUsage = "Duration of time to allow command to start up"
	)

	flag.StringVar(&configfile, "config", "", configUsage)
	flag.StringVar(&configfile, "c", "", configUsage+" (shorthand)")

	flag.StringVar(&listenAddress, "listen", "", listenUsage)
	flag.StringVar(&listenAddress, "l", "", listenUsage+" (shorthand)")

	flag.StringVar(&proxyAddress, "proxy", "", proxyUsage)
	flag.StringVar(&proxyAddress, "p", "", proxyUsage+" (shorthand)")

	flag.DurationVar(&timeout, "timeout", timeout, timeoutUsage)
	flag.DurationVar(&timeout, "t", timeout, timeoutUsage+" (shorthand)")
}

// Load config from cli arguments.
func Load(args []string) (cfg Config, err error) {
	if configfile != "" {
		var e error
		cfg, e = parseConfigFile(configfile)
		if e != nil {
			err = fmt.Errorf("error parsing config file '%s': %s", configfile, e)
			return
		}
	}

	if listenAddress != "" {
		if proxyAddress == "" {
			err = fmt.Errorf("proxyAddress is required")
		}
		cfg.Services = append(cfg.Services, service.Config{
			Proxy: map[string]string{
				listenAddress: proxyAddress,
			},
			Command: args,
			Timeout: timeout.String(),
		})
	} else if proxyAddress != "" {
		err = fmt.Errorf("listenAddress is required")
	}

	return
}
