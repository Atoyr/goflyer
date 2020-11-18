package controller

import (
	"log"
	"os"
	"sync"
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags)
	logMu  sync.Mutex
)

func SetLogger(l *log.Logger) {
	if l == nil {
		l = log.New(os.Stderr, "", log.LstdFlags)
	}
	logMu.Lock()
	logger = l
	logMu.Unlock()
}

func logf(format string, v ...interface{}) {
	logMu.Lock()
	logger.Printf(format, v...)
	logMu.Unlock()
}
