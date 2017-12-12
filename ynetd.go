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
	"os/signal"
	"syscall"
	"time"

	"github.com/rwstauner/ynetd/procman"
)

// Version is the program version, filled in from git during build process.
var Version string

var logger = log.New(os.Stdout, "ynetd ", log.Ldate|log.Ltime|log.Lmicroseconds)

func setupSignals(pm *procman.ProcessManager) {
	channel := make(chan os.Signal, 1)

	signal.Notify(channel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	pm.Signal(<-channel)
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
			logger.Printf("timed out after %s", timeout)
			return
		default:
			if conn, err = dialer.Dial(network, address); err == nil {
				return
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func handleConnection(src *net.TCPConn, dst string, cmd *procman.Process, timeout time.Duration) {
	if cmd != nil {
		cmd.LaunchOnce()
	}

	conn, err := dialWithRetries("tcp", dst, timeout)
	if err != nil {
		src.Close()
		logger.Printf("connect to %s failed: %s", dst, err.Error())
		return
	}

	src.SetKeepAlive(true)
	src.SetKeepAlivePeriod(time.Second * 60)

	fwd := conn.(*net.TCPConn)
	go forward(src, fwd)
	forward(fwd, src)
}

func listen(src string, dst string, cmd *procman.Process, timeout time.Duration) {
	ln, err := net.Listen("tcp", src)
	if err != nil {
		logger.Printf("listen error: %s", err.Error())
		return
	}
	defer ln.Close()

	logger.Printf("listen %s proxy %s cmd: %s", src, dst, cmd)

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Printf("accept error: %s", err.Error())
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

	go listen(listenAddress, proxyAddress, pm.Process(cmd), timeout)

	go setupSignals(pm)

	// block
	pm.Manage()

	return 0
}
