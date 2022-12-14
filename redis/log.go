package redis

import (
	"fmt"
	"log"
)

const (
	packageName = "gopkg/redis"

	logStatusError = "[ERROR]"
)

func logErrorf(format string, args ...interface{}) {
	logf(logStatusError, format, args...)
}

func logf(logStatus, format string, args ...interface{}) {
	log.Printf("%s %s %s", packageName, logStatus, fmt.Sprintf(format, args...))
}
