package driver

import (
	"sync"

	"gorm.io/gorm"
)

const (
	NameMySQL            = "mysql"
	NamePostgreSQL       = "postgresql"
	NamePostgreSQLSimple = "postgresql.simple"
	NameSQLite           = "sqlite"
	NameSQLServer        = "sqlserver"
	NameClickhouse       = "clickhouse"
)

var (
	driverMap   = map[string]Driver{}
	driverMapRW sync.RWMutex
)

// Driver database driver
type Driver interface {

	// Open return GORM database dialector
	Open(dsn string) gorm.Dialector
}

func register(name string, d Driver) {
	driverMapRW.Lock()
	defer driverMapRW.Unlock()
	driverMap[name] = d
}

// Get database driver
func Get(name string) (Driver, bool) {
	driverMapRW.RLock()
	defer driverMapRW.RUnlock()

	d, exist := driverMap[name]
	return d, exist
}
