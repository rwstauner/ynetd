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

	"github.com/rwstauner/ynetd/config"
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

var printVersion = false

func init() {
	flag.BoolVar(&printVersion, "version", printVersion, "Print version")
	procman.SetLogger(logger)
	service.SetLogger(logger)
}

func main() {
	os.Exit(frd())
}

func frd() int {
	flag.Parse()
	if printVersion {
		fmt.Println("ynetd", Version)
		return 0
	}

	cfg, err := config.Load(flag.Args())
	if err != nil {
		fmt.Println(err)
		return 1
	}

	pm := procman.New()
	services, err := config.MakeServices(cfg, pm)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	if len(services) == 0 {
		fmt.Println("no services configured!")
		return 1
	}

	for _, svc := range services {
		go svc.Listen()
	}

	go setupSignals(pm)

	// block
	pm.Manage()

	return 0
}
