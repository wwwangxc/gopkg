package driver

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func init() {
	register(NameSQLServer, &sqlserverDriver{})
}

type sqlserverDriver struct{}

// Open return GORM sqlserver dialector
func (s *sqlserverDriver) Open(dsn string) gorm.Dialector {
	return sqlserver.Open(dsn)
}
