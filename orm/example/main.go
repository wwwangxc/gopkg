package main

import (
	"fmt"

	// gopkg/orm will automatically read configuration
	// files (./app.yaml) when package loaded
	"gorm.io/gorm"

	"github.com/wwwangxc/gopkg/orm"
)

func main() {
	db, err := orm.NewGORM("db_mysql",
		orm.WithDSN(""),                    // set dsn
		orm.WithMaxIdle(20),                // set the maximum number of connections in the idle connection pool.
		orm.WithMaxIdle(1000),              // set the maximum amount of time aconnection may be reused. uint: milliseconds
		orm.WithMaxOpen(20),                // set the maximum number of open connections to the database.
		orm.WithGORMConfig(&gorm.Config{}), // set gorm config. see: https://gorm.io/docs/gorm_config.html
		orm.WithDriver("mysql"))            // set database driver, default mysql
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}

	// MySQL
	db, err = orm.NewGORM("db_mysql")
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}

	// PostgreSQL automatically uses the extended protocol
	db, err = orm.NewGORM("db_postgresql", orm.WithDriver("postgresql"))
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}

	// PostgreSQL disables implicit prepared statement usage
	db, err = orm.NewGORM("db_postgresql", orm.WithDriver("postgresql.simple"))
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}

	// SQLite
	db, err = orm.NewGORM("db_sqlite", orm.WithDriver("sqlite"))
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}

	// SQL Server
	db, err = orm.NewGORM("db_sqlserver", orm.WithDriver("sqlserver"))
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}

	// Clickhouse
	db, err = orm.NewGORM("db_clickhouse", orm.WithDriver("clickhouse"))
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}

	fmt.Println(db.NowFunc().String())
}
