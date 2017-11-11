package main

// TODO: move to pkg in case we wanted multiple
// TODO: kill process after timeout without usage

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func flog(spec string, args ...interface{}) {
	log.Printf("ynetd: "+spec, args...)
}

func launch(args []string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)

	flog("Starting: %s", args)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	go setupSignals(cmd)

	return cmd
}

var process *exec.Cmd
var launchMux = sync.Mutex{}

func launchOnce(cmd []string) {
	launchMux.Lock()
	if process == nil {
		process = launch(cmd)
		time.Sleep(250 * time.Millisecond)
	}
	launchMux.Unlock()
}

func setupSignals(cmd *exec.Cmd) {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel,
		syscall.SIGCHLD,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	for sig := range channel {
		switch sig {
		case syscall.SIGCHLD:
			cmd.Wait()
			// Next client can attempt to restart the command.
			process = nil
		default:
			cmd.Process.Signal(sig)
			// TODO: Allow configuration for which signals to exit with.
			err := cmd.Wait()
			status := 0
			if err != nil {
				if frdErr, ok := err.(*exec.ExitError); ok {
					status = frdErr.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
				}
			}
			os.Exit(status)
		}
	}
}

func forward(src *net.TCPConn, dst *net.TCPConn) {
	defer src.CloseRead()
	defer dst.CloseWrite()
	io.Copy(dst, src)
}

func dialWithRetries(network string, address string, timeout time.Duration) (net.Conn, error) {
	timer := time.After(timeout)
	dialer := net.Dialer{Timeout: timeout}
	var err error
	for {
		select {
		case <-timer:
			flog("timed out after %s", timeout)
			return nil, err
		default:
			if conn, e := dialer.Dial(network, address); e == nil {
				return conn, e
			} else {
				err = e
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func handleConnection(src *net.TCPConn, dst string, cmd []string, timeout time.Duration) {
	launchOnce(cmd)

	conn, err := dialWithRetries("tcp", dst, timeout)
	if err != nil {
		src.Close()
		flog("connect to %s failed: %s", dst, err.Error())
		return
	}

	src.SetKeepAlive(true)
	src.SetKeepAlivePeriod(time.Second * 60)

	fwd := conn.(*net.TCPConn)
	go forward(src, fwd)
	forward(fwd, src)
}

func listen(src string, dst string, cmd []string, timeout time.Duration) {
	ln, err := net.Listen("tcp", src)
	if err != nil {
		flog("listen error: %s", err.Error())
		return
	}
	defer ln.Close()

	flog("listen %s proxy %s cmd: %s", src, dst, cmd)

	for {
		conn, err := ln.Accept()
		if err != nil {
			flog("accept error: %s", err.Error())
			if opErr, ok := err.(*net.OpError); ok {
				if !opErr.Temporary() {
					break
				}
			}
			continue
		}
		go handleConnection(conn.(*net.TCPConn), dst, cmd, timeout)
	}
}

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
}

func main() {
	flag.Parse()
	cmd := flag.Args()
	if printVersion {
		fmt.Println("ynetd", Version)
		os.Exit(0)
	}
	if listenAddress == "" {
		log.Fatal("listenAddress is required")
	}
	if proxyAddress == "" {
		log.Fatal("proxyAddress is required")
	}
	listen(listenAddress, proxyAddress, cmd, timeout)
}
