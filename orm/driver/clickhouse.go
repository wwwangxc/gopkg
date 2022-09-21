package driver

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

func init() {
	register(NameClickhouse, &clickhouseDriver{})
}

type clickhouseDriver struct{}

// Open return GORM clickhouse dialector
func (c *clickhouseDriver) Open(dsn string) gorm.Dialector {
	return clickhouse.Open(dsn)
}
