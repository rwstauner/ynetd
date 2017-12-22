package service

import (
	"io"
	"net"
	"time"

	"github.com/rwstauner/ynetd/procman"
)

// Service represents a single service proxy.
type Service struct {
	Proxy   map[string]string
	Command *procman.Process
	Timeout time.Duration
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
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *Service) handleConnection(src *net.TCPConn, dst string) {
	if s.Command != nil {
		s.Command.LaunchOnce()
	}

	conn, err := dialWithRetries("tcp", dst, s.Timeout)
	if err != nil {
		src.Close()
		logger.Printf("connect to %s failed: %s", dst, err.Error())
		return
	}

	src.SetKeepAlive(true)
	src.SetKeepAlivePeriod(60 * time.Second)

	fwd := conn.(*net.TCPConn)
	go forward(src, fwd)
	forward(fwd, src)
}

func (s *Service) proxy(src string, dst string) {
	ln, err := net.Listen("tcp", src)
	if err != nil {
		logger.Printf("listen error: %s", err.Error())
		return
	}
	defer ln.Close()

	logger.Printf("proxy %s -> %s cmd: %s", src, dst, s.Command)

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
		go s.handleConnection(conn.(*net.TCPConn), dst)
	}
}

// Listen starts listening on the defined ports for incoming connections.
func (s *Service) Listen() error {
	for src, dst := range s.Proxy {
		addrs, err := parseAddr(src)
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			go s.proxy(addr, dst)
		}
	}
	return nil
}
