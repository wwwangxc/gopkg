package driver

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	register(NameMySQL, &mysqlDriver{})
}

type mysqlDriver struct{}

// Open return GORM mysql dialector
func (m *mysqlDriver) Open(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}
