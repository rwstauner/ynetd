// +build !windows

package main

import (
	"syscall"
)

func init() {
	sigChld = syscall.SIGCHLD
}
