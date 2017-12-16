package main

// TODO: kill process after timeout without usage
// TODO: json config file

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rwstauner/ynetd/procman"
	"github.com/rwstauner/ynetd/service"
)

// Version is the program version, filled in from git during build process.
var Version string

var logger = log.New(os.Stdout, "ynetd ", log.Ldate|log.Ltime|log.Lmicroseconds)

func setupSignals(pm *procman.ProcessManager) {
	channel := make(chan os.Signal, 1)

	signal.Notify(channel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	pm.Signal(<-channel)
}

// TODO: make listen/proxy arrays
var listenAddress string
var proxyAddress string
var timeout = 5 * time.Minute
var printVersion = false

func init() {
	const (
		listenUsage  = "Address to listen on"
		proxyUsage   = "Address to proxy to (the address the command should be listening on)"
		timeoutUsage = "Seconds to allow command to start up"
	)

	flag.StringVar(&listenAddress, "listen", "", listenUsage)
	flag.StringVar(&listenAddress, "l", "", listenUsage+" (shorthand)")

	flag.StringVar(&proxyAddress, "proxy", "", proxyUsage)
	flag.StringVar(&proxyAddress, "p", "", proxyUsage+" (shorthand)")

	flag.DurationVar(&timeout, "timeout", timeout, timeoutUsage)
	flag.DurationVar(&timeout, "t", timeout, timeoutUsage+" (shorthand)")

	flag.BoolVar(&printVersion, "version", printVersion, "Print version")

	service.SetLogger(logger)
	procman.SetLogger(logger)
}

func main() {
	os.Exit(frd())
}

func frd() int {
	flag.Parse()
	cmd := flag.Args()
	if printVersion {
		fmt.Println("ynetd", Version)
		return 0
	}
	if listenAddress == "" {
		fmt.Println("listenAddress is required")
		return 1
	}
	if proxyAddress == "" {
		fmt.Println("proxyAddress is required")
		return 1
	}

	pm := procman.New()

	svc := service.Service{
		Proxy: map[string]string{
			listenAddress: proxyAddress,
		},
		Command: pm.Process(cmd),
		Timeout: timeout,
	}

	go svc.Listen()

	go setupSignals(pm)

	// block
	pm.Manage()

	return 0
}
