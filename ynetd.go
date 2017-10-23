package main

// TODO: move to pkg in case we wanted multiple
// TODO: listen for signal, pass to child, wait, exit
// TODO: kill process after timeout without usage
// TODO: stream process output

import (
	"flag"
	"io"
	"log"
	"net"
	"os/exec"
	"time"
)

var retries = 60 * 5 * 4 // 60s * 5min * 4 = every 250ms
var process = (*exec.Cmd)(nil)

func flog(spec string, args ...interface{}) {
	log.Printf(spec, args...)
}

func launch(args []string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)

	flog("Starting: %s", args)

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return cmd
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

func main() {
	flag.Parse()
	args := flag.Args()
	src, dst, cmd := args[0], args[1], args[2:]
	listen(src, dst, cmd)
}
