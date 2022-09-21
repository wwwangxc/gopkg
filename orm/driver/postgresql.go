package driver

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	register(NamePostgreSQL, &postgresqlDriver{})
}

type postgresqlDriver struct{}

// Open return GORM postgre sql dialector
//
// automatically uses the extended protocol
func (p *postgresqlDriver) Open(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}
