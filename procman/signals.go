package procman

import (
	"fmt"
	"os"
	"syscall"
)

func getSignal(sig string, defaultSig os.Signal) os.Signal {
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
