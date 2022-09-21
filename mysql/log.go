package mysql

import (
	"fmt"
	"log"
)

const (
	packageName = "gopkg/mysql"

	logStatusError = "[ERROR]"
)

func logErrorf(format string, args ...interface{}) {
	logf(logStatusError, format, args...)
}

func logf(logStatus, format string, args ...interface{}) {
	log.Printf("%s %s %s", packageName, logStatus, fmt.Sprintf(format, args...))
}
