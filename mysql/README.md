# gopkg MySQL Package


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
    // new mysql client proxy.
    // mysql.WithDSN(""): set dsn
    // mysql.WithMaxIdle(20): sets the maximum number of connections in the idle connection pool.
    // mysql.WithMaxIdleTime(1000) sets the maximum amount of time a connection may be reused. uint: milliseconds
    // mysql.WithMaxOpen(20) sets the maximum number of open connections to the database.
    cli := mysql.NewClientProxy("client1", mysql.WithDSN(""), mysql.WithMaxIdle(20), mysql.WithMaxIdleTime(1000), mysql.WithMaxOpen(20))
    
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

### SQLBuilder

```go
package main

import (
	"fmt"

	"github.com/wwwangxc/gopkg/mysql/sqlbuilder"
)

func main() {
    builder := sqlbuilder.NewSelect("table_name",
        sqlbuilder.WithFieldEqual("f_equal", 1),
        sqlbuilder.WithFieldLike("f_like", "val"),
        sqlbuilder.WithFieldLessThan("f_less", 111),
        sqlbuilder.WithFieldLessOrEqualThan("f_less_eq", 222),
        sqlbuilder.WithFieldGreaterThan("f_greater", 333),
        sqlbuilder.WithFieldGreaterOrEqualThan("f_greater_equal", 444),
        sqlbuilder.WithOrderBy("f_asc"),
        sqlbuilder.WithOrderByDESC("f_desc"),
        sqlbuilder.WithLimit(100),
        sqlbuilder.WithOffset(1),
        sqlbuilder.WithForceIndex("index_name"),
    )
    
    // sql:
    //     SELECT select_field FROM table_name
    //     FORCE INDEX(index_name)
    //     WHERE f_equal=? AND f_like LIKE ? AND f_less<? AND f_less_equal<=? AND f_greater>? AND f_greater_equal>=?
    //     ORDER BY f_asc ASC, f_desc DESC
    //     LIMIT 100 OFFSET 1
    //
    // args:
    //     [1, "%val%", 111, 222, 333, 444]
    sql, args := builder.Build("select_field")
    fmt.Printf("generate sql: %s\n", sql)
    fmt.Println(args)
}
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
