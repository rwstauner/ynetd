package service

import (
	"io"
	"net"
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

func (s *Service) shouldStop() bool {
	return s.StopAfter > 0
}

func (s *Service) proxy(src string, dst string) error {
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

		for {
			select {
			case conn := <-conns:
				go s.handleConnection(conn, dst)
				if s.shouldStop() {
					if stopTimer == nil {
						stopTimer = time.NewTimer(s.StopAfter)
						stopTimeChan = stopTimer.C
					} else {
						if !stopTimer.Stop() {
							// If the timer expired when we didn't handle it.
							if !stopped {
								<-stopTimer.C // drain
							}
						}
						// TODO: We should probably reset this when a connection closes
						// rather than when it opens.
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
