package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rwstauner/ynetd/service"
)

var configfile string
var listenAddress string
var proxySpec string
var timeout = service.DefaultTimeout

func init() {
	const (
		configUsage  = "Path to configuration file"
		listenUsage  = "Address to listen on (deprecated)"
		proxyUsage   = "Addresses to proxy, separated by spaces (\"fromhost:port tohost:port from to\")"
		timeoutUsage = "Duration of time to allow command to start up"
	)

	flag.StringVar(&configfile, "config", "", configUsage)
	flag.StringVar(&configfile, "c", "", configUsage+" (shorthand)")

	flag.StringVar(&listenAddress, "listen", "", listenUsage)
	flag.StringVar(&listenAddress, "l", "", listenUsage+" (shorthand)")

	flag.StringVar(&proxySpec, "proxy", "", proxyUsage)
	flag.StringVar(&proxySpec, "p", "", proxyUsage+" (shorthand)")

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

	proxy := make(map[string]string)
	if listenAddress != "" {
		fmt.Fprintln(os.Stderr, "-listen is deprecated.  Use -proxy 'from:port to:port'")
		if proxySpec == "" {
			err = fmt.Errorf("-proxy is required")
			return
		}
		proxy[listenAddress] = proxySpec
	} else if proxySpec != "" {
		addrs := strings.Split(proxySpec, " ")
		if len(addrs)%2 != 0 {
			err = fmt.Errorf("-proxy must contain pairs of addresses: \"from1 to1 from2 to2\"")
			return
		}
		var key string
		for i, s := range addrs {
			if i%2 == 0 {
				key = s
			} else {
				proxy[key] = s
			}
		}
	} else if len(args) > 0 {
		err = fmt.Errorf("-proxy is required")
		return
	}

	if len(proxy) > 0 {
		cfg.Services = append(cfg.Services, service.Config{
			Proxy:   proxy,
			Command: args,
			Timeout: timeout.String(),
		})
	}

	return
}
