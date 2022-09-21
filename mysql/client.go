package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// TxFunc function executed in transaction
type TxFunc func(*sql.Tx) error

// ScanFunc scan function
type ScanFunc func(*sql.Rows) error

// ClientProxy MySQL client proxy
//go:generate mockgen -source=client.go -destination=mockmysql/mysql_mock.go -package=mockmysql
type ClientProxy interface {
	// Exec executes a query without returning any rows
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Transaction auto start and commit transcation
	//
	// Commit the transaction when returns empty error.
	// Rollback the transaction when returns not empty error.
	Transaction(ctx context.Context, f TxFunc, opts ...TxOption) error

	// Query executes a query that returns rows, typically a SELECT.
	//
	// The args are for any placeholder parameters in the query.
	// Loop executes scan function when returns rows not empty.
	Query(ctx context.Context, f ScanFunc, query string, args ...interface{}) error

	// QueryRow executes a query that is expected to return at most one row.
	//
	// Scan the columns from the matched row into the values pointed at by dest.
	// If more than one row matches the query, will uses the first row and discards the rest.
	// sql.ErrNoRows is returned if the result set is empty.
	QueryRow(ctx context.Context, dest []interface{}, query string, args ...interface{}) error

	// Select executes a query and storing the matched row into the
	// struct slice pointed at by dest.
	//
	// If you have null fields and use SELECT *, you must use sql.Null* in your struct.
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Get executes a query that is expected to return at most one row
	// and storing the result set into the struct pointed at by dest.
	//
	// If you have null fields and use SELECT *, you must use sql.Null* in your struct.
	// sql.ErrNoRows is returned if the result set is empty.
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type clientProxyImpl struct {
	name string
	opts []Option
}

// NewClientProxy new myql client proxy
func NewClientProxy(name string, opts ...Option) ClientProxy {
	return &clientProxyImpl{
		name: name,
		opts: opts,
	}
}

// Exec executes a query without returning any rows
func (c *clientProxyImpl) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	db, err := c.getDB()
	if err != nil {
		return nil, err
	}

	return db.ExecContext(ctx, query, args...)
}

// Transaction auto start and commit transcation
//
// Commit the transaction when returns empty error.
// Rollback the transaction when returns not empty error.
func (c *clientProxyImpl) Transaction(ctx context.Context, f TxFunc, opts ...TxOption) error {
	txOpts := &sql.TxOptions{}
	for _, opt := range opts {
		opt(txOpts)
	}

	db, err := c.getDB()
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, txOpts)
	if err != nil {
		return err
	}

	if err = f(tx); err != nil {
		if e := tx.Rollback(); e != nil {
			return e
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		if e := tx.Rollback(); e != nil {
			return e
		}

		return err
	}

	return nil
}

// Query executes a query that returns rows, typically a SELECT.
//
// The args are for any placeholder parameters in the query.
// Loop executes scan function when returns rows not empty.
func (c *clientProxyImpl) Query(ctx context.Context, f ScanFunc, query string, args ...interface{}) error {
	db, err := c.getDB()
	if err != nil {
		return err
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err = f(rows); err != nil {
			return err
		}
	}

	return rows.Err()
}

// QueryRow executes a query that is expected to return at most one row.
//
// Scan the columns from the matched row into the values pointed at by dest.
// If more than one row matches the query, will uses the first row and discards the rest.
// sql.ErrNoRows is returned if the result set is empty.
func (c *clientProxyImpl) QueryRow(ctx context.Context, dest []interface{}, query string, args ...interface{}) error {
	db, err := c.getDB()
	if err != nil {
		return err
	}

	row := db.QueryRowContext(ctx, query, args...)
	if err := row.Scan(dest...); err != nil {
		return err
	}

	return nil
}

// Select executes a query and storing the matched row into the
// struct slice pointed at by dest.
//
// If you have null fields and use SELECT *, you must use sql.Null* in your struct.
func (c *clientProxyImpl) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	db, err := c.getDB()
	if err != nil {
		return err
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return sqlx.StructScan(rows, dest)
}

// Get executes a query that is expected to return at most one row
// and storing the result set into the struct pointed at by dest.
//
// If you have null fields and use SELECT *, you must use sql.Null* in your struct.
// sql.ErrNoRows is returned if the result set is empty.
func (c *clientProxyImpl) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	db, err := c.getDB()
	if err != nil {
		return err
	}

	return sqlx.NewDb(db, "mysql").GetContext(ctx, dest, query, args...)
}

func (c *clientProxyImpl) getDB() (*sql.DB, error) {
	return getDB(c.name, c.opts...)
}
