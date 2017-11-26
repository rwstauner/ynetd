package main

// TODO: move to pkg in case we wanted multiple
// TODO: kill process after timeout without usage
// TODO: json config file

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

// Version is the program version, filled in from git during build process.
var Version string

var logger = log.New(os.Stdout, "ynetd ", log.Ldate|log.Ltime|log.Lmicroseconds)

func flog(spec string, args ...interface{}) {
	logger.Printf(spec, args...)
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
	flog("child started: %d", cmd.Process.Pid)

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

var sigChld os.Signal

func setupSignals() {
	channel := make(chan os.Signal, 1)

	// All signals.
	signal.Notify(channel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	if sigChld != nil {
		signal.Notify(channel, sigChld)
	}

	for sig := range channel {
		switch sig {
		case sigChld:
			if process == nil {
				continue
			}

			process.Wait()
			// Next client can attempt to restart the command.
			// FIXME: reap all child processes and judge by pid when to restart.
			process = nil
			flog("child reaped")
		default:
			if process == nil {
				os.Exit(0)
			}

			flog("sending %s to %d", sig, process.Process.Pid)
			if err := process.Process.Signal(sig); err != nil {
				flog("error: %s", err)
			}
			// TODO: Allow configuration for which signals to exit with.
			err := process.Wait()
			status := 0
			if err != nil {
				if frdErr, ok := err.(*exec.ExitError); ok {
					flog("process state: %s", frdErr.ProcessState)
					status = frdErr.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
				}
			}
			flog("waited (%d): %s", status, err)
			os.Exit(status)
		}
	}
}

func forward(src *net.TCPConn, dst *net.TCPConn) {
	defer src.CloseRead()
	defer dst.CloseWrite()
	io.Copy(dst, src)
}

func dialWithRetries(network string, address string, timeout time.Duration) (conn net.Conn, err error) {
	timer := time.After(timeout)
	dialer := net.Dialer{Timeout: timeout}
	for {
		select {
		case <-timer:
			flog("timed out after %s", timeout)
			return
		default:
			if conn, err = dialer.Dial(network, address); err == nil {
				return
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func handleConnection(src *net.TCPConn, dst string, cmd []string, timeout time.Duration) {
	// TODO: make cmd optional
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

	go listen(listenAddress, proxyAddress, cmd, timeout)

	// block
	setupSignals()
	return 0
}
