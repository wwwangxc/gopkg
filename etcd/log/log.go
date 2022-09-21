package log

import (
	"fmt"
	"log"
)

const (
	packageName = "go-pkg/etcd"

	logStatusError = "[ERROR]"
	logStatusWarn  = "[WARN]"
)

func Errorf(format string, args ...interface{}) {
	logf(logStatusError, format, args...)
}

func Warn(format string, args ...interface{}) {
	logf(logStatusWarn, format, args...)
}

func logf(logStatus, format string, args ...interface{}) {
	log.Printf("%s %s %s", packageName, logStatus, fmt.Sprintf(format, args...))
}
