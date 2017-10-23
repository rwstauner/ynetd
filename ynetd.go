package main

// TODO: move to pkg in case we wanted multiple
// TODO: kill process after timeout without usage

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var retries = 60 * 5 * 4 // 60s * 5min * 4 = every 250ms
var process = (*exec.Cmd)(nil)

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

func setupSignals(cmd *exec.Cmd) {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

	for sig := range channel {
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

func forward(src *net.TCPConn, dst *net.TCPConn) {
	defer src.CloseRead()
	defer dst.CloseWrite()
	io.Copy(dst, src)
}

func dialWithRetries(network string, address string, retries int) (net.Conn, error) {
	conn, err := (net.Conn)(nil), (error)(nil)
	for i := 0; i < retries; i++ {
		conn, err = net.Dial(network, address)
		if err == nil {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
	return conn, err
}

func handleConnection(src *net.TCPConn, dst string, cmd []string) {
	if process == nil {
		process = launch(cmd)
		time.Sleep(250 * time.Millisecond)
	}

	conn, err := dialWithRetries("tcp", dst, retries)
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

func listen(src string, dst string, cmd []string) {
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
		go handleConnection(conn.(*net.TCPConn), dst, cmd)
	}
}

var listenAddress string
var proxyAddress string

func init() {
	const (
		listenUsage = "Address to listen on"
		proxyUsage  = "Address to proxy to (the address the command should be listening on)"
	)
	flag.StringVar(&listenAddress, "listen", "", listenUsage)
	flag.StringVar(&listenAddress, "l", "", listenUsage+" (shorthand)")
	flag.StringVar(&proxyAddress, "proxy", "", proxyUsage)
	flag.StringVar(&proxyAddress, "p", "", proxyUsage+" (shorthand)")
}

func main() {
	flag.Parse()
	cmd := flag.Args()
	if listenAddress == "" {
		log.Fatal("listenAddress is required")
	}
	if proxyAddress == "" {
		log.Fatal("proxyAddress is required")
	}
	listen(listenAddress, proxyAddress, cmd)
}
