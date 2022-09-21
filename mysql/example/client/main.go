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
