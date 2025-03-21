# gopkg MySQL Package

[![Go Report Card](https://goreportcard.com/badge/github.com/wwwangxc/gopkg/mysql)](https://goreportcard.com/report/github.com/wwwangxc/gopkg/mysql)
[![GoDoc](https://pkg.go.dev/badge/github.com/wwwangxc/gopkg/mysql?status.svg)](https://pkg.go.dev/github.com/wwwangxc/gopkg/mysql)
[![OSCS Status](https://www.oscs1024.com/platform/badge/wwwangxc/gopkg.svg?size=small)](https://www.murphysec.com/dr/c1TuOdJ62DzT0agLwg)

gopkg/config is an componentized mysql package.

It provides an easy way to configre and manage mysql client.

Based on [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) and [jmoiron/sqlx](https://github.com/jmoiron/sqlx).

## Required

Go >= 1.15

## Install

```sh
go get github.com/wwwangxc/gopkg/mysql
```

## Quick Start

### Client

**main.go**

```go
package main

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    
    // gopkg/mysql will automatically read configuration
    // files (./app.yaml) when package loaded
    "github.com/wwwangxc/gopkg/mysql"
)

func main() {
	// default config is read from ./app.yaml
	// you can also set config by code
	// mysql.LoadConfig("./app.yaml")

    // new mysql client proxy.
    cli := mysql.NewClientProxy("client1",
        mysql.WithDSN(""),            // set dsn
        mysql.WithMaxIdle(20),        // sets the maximum number of connections in the idle connection pool.
        mysql.WithMaxIdleTime(1000),  // sets the maximum amount of time a connection may be reused. uint: milliseconds
        mysql.WithMaxOpen(20))        // sets the maximum number of open connections to the database.
    
    // Insert
    result, err := cli.Exec(context.TODO(), "INSERT INTO user (name) VALUES (?)", "wwwangxc")
    if err != nil {
        fmt.Printf("exec fail. error:%v", err)
        return
    }
    
    lastID, _ := result.LastInsertId()
    fmt.Printf("last id: %d", lastID)
    
    // Transaction
    err = cli.Transaction(context.TODO(), func(t *sql.Tx) error {
        // do somothing...
        // return error will rollback transaction
        return nil
    })
    if err != nil {
        fmt.Printf("transaction fail. error:%v", err)
        return
    }
    
    // Query
    var users []*User
    scanFunc := func(rows *sql.Rows) error {
        user := &User{}
        if err := rows.Scan(&user.Name); err != nil {
    	    return err
        }
    
        users = append(users, user)
        return nil
    }

    err = cli.Query(context.TODO(), scanFunc, "SELECT name FROM user")
    if err != nil {
        fmt.Printf("query fail. error:%v", err)
        return
    }
    
    // QueryRow
    // get first result
    // return database/sql.ErrNoRows when record not found
    user := &User{}
    err = cli.QueryRow(context.TODO(), []interface{}{&user.Name}, "SELECT name FROM user")
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            fmt.Println("not found")
            return
        }
    
        fmt.Printf("query row fail. error:%v", err)
        return
    }
    
    // Select
    // storing results in a []*User
    err = cli.Select(context.TODO(), &users, "SELECT name FROM user")
    if err != nil {
        fmt.Printf("select fail. error:%v", err)
        return
    }
    
    // Get
    // get a single result and storing result in a *User
    // return database/sql.ErrNoRows when record not found
    err = cli.Get(context.TODO(), &user, "SELECT name FROM user")
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            fmt.Println("not found")
            return
        }
    
    	fmt.Printf("get fail. error:%v", err)
    	return
    }
}

type User struct {
    Name string `db:"name"`
}
```

**app.yaml**

```yaml
client:
  mysql:
    max_idle: 11
    max_open: 22
    max_idle_time: 33
  service:
    - name: client1
      dsn: root:root@tcp(127.0.0.1:3306)/db1?charset=utf8&parseTime=True
    - name: client2
      dsn: root:root@tcp(127.0.0.1:3306)/db2?charset=utf8&parseTime=True
      max_idle: 111
      max_open: 222
      max_idle_time: 333
```

## How To Mock

```go
package tests

import (
    "testing"
    
    "github.com/agiledragon/gomonkey"
    "github.com/golang/mock/gomock"

    "github.com/wwwangxc/gopkg/mysql"
    "github.com/wwwangxc/gopkg/mysql/mockmysql"
)

func TestMock(t *testing.T){
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockCli := mockmysql.NewMockClientProxy(ctrl)
    mockCli.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
        Return(nil).AnyTimes()
    
    mockCli.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
        DoAndReturn(func(_ context.Context, _ mysql.TxFunc, opts ...mysql.TxOption) error {
            return nil
        })
    
    patches := gomonkey.ApplyFunc(mysql.NewClientProxy,
        func(string, ...mysql.ClientProxyOption) (mysql.ClientProxy, error) {
            return mockCli, nil
        })
    defer patches.Reset()
}
```
