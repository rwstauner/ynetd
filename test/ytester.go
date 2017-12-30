package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	after      = 0 * time.Millisecond
	before     = 0 * time.Millisecond
	knock      = false
	loop       = false
	port       = ""
	send       = ""
	serve      = ""
	serveAfter = 0 * time.Second
	timeout    = 2 * time.Second
	logger     = log.New(os.Stderr, "ytester ", log.Ldate|log.Ltime|log.Lmicroseconds)
)

func init() {
	flag.DurationVar(&after, "after", after, "after")
	flag.DurationVar(&before, "before", before, "before")
	flag.BoolVar(&knock, "knock", knock, "knock")
	flag.BoolVar(&loop, "loop", loop, "loop")
	flag.StringVar(&port, "port", port, "port")
	flag.StringVar(&send, "send", send, "send")
	flag.StringVar(&serve, "serve", serve, "serve")
	flag.DurationVar(&serveAfter, "serve-after", serveAfter, "serve-after")
	flag.DurationVar(&timeout, "timeout", timeout, "timeout")
}

func flog(spec string, args ...interface{}) {
	logger.Printf(spec, args...)
}

func listen(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		flog("listen error: %s", err)
		return
	}
	defer ln.Close()

	flog("listening %s", port)

	msg := "not yet"
	time.AfterFunc(serveAfter, func() { msg = serve })

	for {
		conn, err := ln.Accept()
		if err != nil {
			flog("listen error: %s", err)
			continue
		}

		flog("serving: %s", msg)
		conn.Write([]byte(msg + "\n"))
		conn.Close()

		time.Sleep(after)
		if !loop {
			break
		}
	}
}

func main() {
	os.Exit(frd())
}

func frd() int {
	flog("args: %s", os.Args[:])
	flag.Parse()

	time.Sleep(before)

	switch {
	case knock:
		flog("knocking %d", port)
		net.Dial("tcp", "localhost:"+port)
	case send != "":
		c := make(chan net.Conn)
		go func() {
			for {
				flog("dialing %d", port)
				conn, err := net.DialTimeout("tcp", "localhost:"+port, timeout)
				if err == nil {
					c <- conn
					break
				}
				time.Sleep(250 * time.Millisecond)
			}
		}()
		select {
		case conn := <-c:
			flog("sending: %s", send)
			conn.Write([]byte(send + "\n"))
			io.Copy(os.Stdout, conn)
			conn.Close()
		case <-time.After(timeout):
			flog("timed out after %s", timeout)
			return 1
		}
	case serve != "":
		if port == "" {
			flog("port is required")
			return 1
		}
		listen(":" + port)
	}
	return 0
}
