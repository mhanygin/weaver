package logger

import (
    "os"

    "github.com/op/go-logging"
)

func init() {
    logger := logging.NewLogBackend(os.Stderr, "", 0)
    formater := logging.NewBackendFormatter(logger, format)
    loggerLeveled := logging.AddModuleLevel(logger)
    loggerLeveled.SetLevel(logging.DEBUG, "")

    logging.SetBackend(loggerLeveled, formater)
}

var Log = logging.MustGetLogger("serve-runner")

var format = logging.MustStringFormatter(
    `%{time:15:04:05.000} %{shortfunc} %{level:.5s} %{message}`,
)
