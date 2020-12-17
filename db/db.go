package db

import (
  "time"
  "log"
  "sync"
  "os"

  "github.com/atoyr/goflyer/client/bitflyer"
)

type Database interface {
  InsertExecution(executions []bitflyer.Execution)
  SelectExecution(after time.Time, count int) ([]bitflyer.Execution)
}

var (
  logger = log.New(os.Stderr, "", log.LstdFlags)
  logMu sync.Mutex
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
