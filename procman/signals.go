package procman

import (
	"fmt"
	"syscall"
)

func getSignal(sig string, defaultSig syscall.Signal) syscall.Signal {
	switch sig {
	case "":
		return defaultSig
	case "INT":
		return syscall.SIGINT
	case "TERM":
		return syscall.SIGTERM
	default:
		panic(fmt.Sprintf("Unknown signal '%s'", sig))
	}
}
