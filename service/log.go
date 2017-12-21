package service

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "ynetd ", log.Ldate|log.Ltime|log.Lmicroseconds)

// SetLogger sets the logger used in this package.
func SetLogger(l *log.Logger) {
	logger = l
}
