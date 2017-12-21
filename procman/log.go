package procman

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "procman ", log.Ldate|log.Ltime|log.Lmicroseconds)

// SetLogger sets the logger used in this package.
func SetLogger(l *log.Logger) {
	logger = l
}
