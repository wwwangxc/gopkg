# gopkg ORM Package

[![Go Report Card](https://goreportcard.com/badge/github.com/wwwangxc/gopkg/orm)](https://goreportcard.com/report/github.com/wwwangxc/gopkg/orm)
[![GoDoc](https://pkg.go.dev/badge/github.com/wwwangxc/gopkg/orm?status.svg)](https://pkg.go.dev/github.com/wwwangxc/gopkg/orm)
[![OSCS Status](https://www.oscs1024.com/platform/badge/wwwangxc/gopkg.svg?size=small)](https://www.murphysec.com/dr/c1TuOdJ62DzT0agLwg)

gopkg/orm is an componentized orm package.

it provides an easy way to configure and manage gorm.DB.

Based on [gorm.io/gorm](https://github.com/go-gorm/gorm).

## Install

```sh
go get github.com/wwwangxc/gopkg/orm
```

## Quick Start

**main.go**

```go
package main

import (
	"fmt"

	// gopkg/orm will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/orm"
	"gorm.io/gorm"
)

func main() {
    // default config is read from ./app.yaml
    // you can also set config by code
    // orm.LoadConfig("./app.yaml")

	db, err := orm.NewGORMProxy("db_mysql",
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
	db, err = orm.NewGORMProxy("db_mysql")
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}

	// PostgreSQL automatically uses the extended protocol
	db, err = orm.NewGORMProxy("db_postgresql", orm.WithDriver("postgresql"))
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}

	// PostgreSQL disables implicit prepared statement usage
	db, err = orm.NewGORMProxy("db_postgresql", orm.WithDriver("postgresql.simple"))
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}

	// SQLite
	db, err = orm.NewGORMProxy("db_sqlite", orm.WithDriver("sqlite"))
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}

	// SQL Server
	db, err = orm.NewGORMProxy("db_sqlserver", orm.WithDriver("sqlserver"))
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}

	// Clickhouse
	db, err = orm.NewGORMProxy("db_clickhouse", orm.WithDriver("clickhouse"))
	if err != nil {
		fmt.Printf("new gorm proxy fail. error:%v", err)
	}
}
```

**app.yaml***

```yaml
client:
  mysql:
    max_idle: 11
    max_open: 22
    max_idle_time: 33
  postgresql:
    max_idle: 11
    max_open: 22
    max_idle_time: 33
  sqlite:
    max_idle: 11
    max_open: 22
    max_idle_time: 33
  sqlserver:
    max_idle: 11
    max_open: 22
    max_idle_time: 33
  clickhouse:
    max_idle: 11
    max_open: 22
    max_idle_time: 33
  service:
    - name: db_mysql
      dsn: root:root@tcp(127.0.0.1:3306)/db1?charset=utf8&parseTime=True
      driver: mysql
      max_idle: 111
      max_open: 222
      max_idle_time: 333
    - name: db_postgresql
      dsn: "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
      driver: postgresql
    - name: db_sqlite
      dsn: gopkg.db
      driver: sqlite
    - name: db_sqlserver
      dsn: sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm
      driver: sqlserver
    - name: db_clickhouse
      dsn: tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20
      driver: clickhouse
```
