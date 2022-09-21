package driver

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	register(NamePostgreSQLSimple, &postgreSimpleDriver{})
}

type postgreSimpleDriver struct{}

// Open return GORM postgre sql dialector
//
// uses the simple protocol
func (p *postgreSimpleDriver) Open(dsn string) gorm.Dialector {
	return postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	})
}
