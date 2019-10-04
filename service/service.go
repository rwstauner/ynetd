package service

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/rwstauner/ynetd/procman"
)

// Service represents a single service proxy.
type Service struct {
	Proxy     map[string]string
	Command   *procman.Process
	Timeout   time.Duration
	StopAfter time.Duration
}

func forward(src *net.TCPConn, dst *net.TCPConn, done chan bool) {
	defer src.CloseRead()
	defer dst.CloseWrite()
	io.Copy(dst, src)
	done <- true
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
			// Will retrying help?
			if opErr, ok := err.(*net.OpError); ok {
				if !opErr.Temporary() {
					switch opErr.Err.(type) {
					// Not for bad ports, etc.
					case *net.AddrError, net.InvalidAddrError:
						return
					// Some errors are not considered Temporary
					// but are precisely the sort of errors we want to retry:
					// - connection refused (*os.SyscallError)
					// - no such host (*.netDNSError)
					default:
					}
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *Service) handleConnection(src *net.TCPConn, dst string, done chan bool) {
	if len(dst) > 5 && dst[0:6] == "exec:/" {
		var stdout, stderr bytes.Buffer
		path := dst[5:len(dst)]
		cmd := exec.Command(path)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		dst = strings.TrimSpace(stdout.String())
		if err != nil || dst == "" {
			src.Close()
			logger.Printf("failed to get address from %s (%s): %s", path, err.Error(), stderr.String())
			return
		}
	}

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
	c := make(chan bool, 2)
	go forward(src, fwd, c)
	go forward(fwd, src, c)
	// Wait for input and output to close.
	<-c
	<-c

	done <- true
}

func (s *Service) shouldStop() bool {
	return s.StopAfter > 0
}

func (s *Service) proxy(src string, dst string) error {
	if dst == "" {
		return fmt.Errorf("destination address is required")
	}

	ln, err := net.Listen("tcp", src)
	if err != nil {
		return err
	}
	go func() {
		defer ln.Close()

		logger.Printf("proxy %s -> %s cmd: %s", src, dst, s.Command)

		conns := make(chan *net.TCPConn)
		go func() {
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
				conns <- conn.(*net.TCPConn)
			}
		}()

		var stopTimeChan <-chan time.Time
		var stopTimer *time.Timer
		stopped := false
		clientFinished := make(chan bool)
		clients := 0

		for {
			select {
			case conn := <-conns:
				// Stop any previous timer as soon as a connection is made.
				if stopTimer != nil {
					if !stopTimer.Stop() {
						// If the timer expired when we didn't handle it.
						if !stopped {
							<-stopTimer.C // drain
						}
					}
					stopped = true
				}
				clients++
				go s.handleConnection(conn, dst, clientFinished)

			case <-clientFinished:
				clients--
				// If configured to stop and all clients have finished, (re)start timer.
				if s.shouldStop() && clients == 0 {
					if stopTimer == nil {
						stopTimer = time.NewTimer(s.StopAfter)
						stopTimeChan = stopTimer.C
					} else {
						stopTimer.Reset(s.StopAfter)
					}
					stopped = false
				}

			case <-stopTimeChan:
				stopped = true
				s.Command.Stop()
			}
		}
	}()
	return nil
}

// Listen starts listening on the defined ports for incoming connections.
func (s *Service) Listen() error {
	for src, dst := range s.Proxy {
		addrs, err := parseAddr(src)
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			err := s.proxy(addr, dst)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
