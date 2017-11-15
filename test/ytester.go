package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var (
	after  = 0 * time.Millisecond
	before = 0 * time.Millisecond
	knock  = false
	loop   = true
	port   = ""
	send   = ""
	serve  = ""
)

func init() {
	flag.DurationVar(&after, "after", after, "after")
	flag.DurationVar(&before, "before", before, "before")
	flag.BoolVar(&knock, "knock", knock, "knock")
	flag.BoolVar(&loop, "loop", loop, "loop")
	flag.StringVar(&port, "port", port, "port")
	flag.StringVar(&send, "send", send, "send")
	flag.StringVar(&serve, "serve", serve, "serve")
}

func flog(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func listen(addr string) {
	ln, _ := net.Listen("tcp", addr)
	defer ln.Close()

	for {
		conn, _ := ln.Accept()

		flog("serving:", serve)
		conn.Write([]byte(serve + "\n"))
		conn.Close()

		time.Sleep(after)
		if !loop {
			break
		}
	}
}

func main() {
	flog("ytester:", os.Args[:])
	flag.Parse()

	time.Sleep(before)

	switch {
	case knock:
		net.Dial("tcp", "localhost:"+port)
	case send != "":
		conn, err := net.DialTimeout("tcp", "localhost:"+port, 10*time.Second)
		if err != nil {
			flog("send error:", err)
			os.Exit(1)
		}
		conn.Write([]byte(send + "\n"))
		io.Copy(os.Stdout, conn)
		conn.Close()
	default:
		if serve != "" {
			listen(":" + port)
		}
	}
}
