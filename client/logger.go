package client

import(
  "log"
)

var Logger *log.Logger

func logf(format string, v ...interface{}) {
  if Logger == nil {
    log.Printf(format, v...)
    return
  }
  Logger.Printf(format, v...)
}
