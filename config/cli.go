package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var configfile string
var listenAddress string
var proxySep = " "
var proxySpec string
var timeout = DefaultTimeout
var stopAfter = DefaultStopAfter
var stopSignal = "INT"
var waitAfterStart = DefaultWaitAfterStart

func init() {
	const (
		configUsage         = "Path to configuration file"
		listenUsage         = "Address to listen on (deprecated)"
		proxySepUsage       = "Separator character for -proxy"
		proxyUsage          = "Addresses to proxy, separated by spaces (\"fromhost:port tohost:port from to\")"
		timeoutUsage        = "Duration of time to allow connections to attempt to forward"
		stopAfterUsage      = "Duration of time after the last client connection to stop the command"
		stopSignalUsage     = "Signal to send to stop"
		waitAfterStartUsage = "Duration of time to wait while command starts before forwarding"
	)

	flag.StringVar(&configfile, "config", "", configUsage)
	flag.StringVar(&configfile, "c", "", configUsage+" (shorthand)")

	flag.DurationVar(&stopAfter, "stop-after", stopAfter, stopAfterUsage)
	flag.StringVar(&stopSignal, "stop-signal", stopSignal, stopSignalUsage)

	flag.StringVar(&listenAddress, "listen", "", listenUsage)
	flag.StringVar(&listenAddress, "l", "", listenUsage+" (shorthand)")

	flag.StringVar(&proxySep, "proxy-sep", proxySep, proxyUsage)
	flag.StringVar(&proxySep, "ps", proxySep, proxyUsage+" (shorthand)")

	flag.StringVar(&proxySpec, "proxy", "", proxyUsage)
	flag.StringVar(&proxySpec, "p", "", proxyUsage+" (shorthand)")

	flag.DurationVar(&timeout, "timeout", timeout, timeoutUsage)
	flag.DurationVar(&timeout, "t", timeout, timeoutUsage+" (shorthand)")

	flag.DurationVar(&waitAfterStart, "wait-after-start", waitAfterStart, waitAfterStartUsage)
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
		addrs := strings.Split(proxySpec, proxySep)
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
		cfg.Services = append(cfg.Services, Service{
			Proxy:          proxy,
			Command:        args,
			Timeout:        timeout.String(),
			StopAfter:      stopAfter.String(),
			StopSignal:     stopSignal,
			WaitAfterStart: waitAfterStart.String(),
		})
	}

	return
}
