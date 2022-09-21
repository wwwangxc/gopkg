package config

import (
	"fmt"
	"log"
)

const (
	packageName = "gopkg/config"

	logStatusError = "[ERROR]"
	logStatusInfo  = "[INFO]"
)

func logErrorf(format string, args ...interface{}) {
	logf(logStatusError, format, args...)
}

func logInfo(format string, args ...interface{}) {
	logf(logStatusInfo, format, args...)
}

func logf(logStatus, format string, args ...interface{}) {
	log.Printf("%s %s %s", packageName, logStatus, fmt.Sprintf(format, args...))
}
