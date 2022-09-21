package mysql

import "database/sql"

// TxOption transaction option
type TxOption func(*sql.TxOptions)

// WithIsolation set transaction isolation level
//
// Isolation is the transaction isolation level.
// If zero, the driver or database's default level is used.
func WithIsolation(isolation sql.IsolationLevel) TxOption {
	return func(options *sql.TxOptions) {
		options.Isolation = isolation
	}
}

// WithReadOnly set transaction readonly
func WithReadOnly(readOnly bool) TxOption {
	return func(options *sql.TxOptions) {
		options.ReadOnly = readOnly
	}
}

// Option mysql client proxy option
type Option func(*serviceConfig)

// WithDSN set dsn
func WithDSN(dsn string) Option {
	return func(cfg *serviceConfig) {
		cfg.DSN = dsn
	}
}

// WithMaxIdle sets the maximum number of connections in the idle
// connection pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
func WithMaxIdle(maxIdel int) Option {
	return func(cfg *serviceConfig) {
		cfg.MaxIdle = maxIdel
	}
}

// WithMaxOpen sets the maximum number of open connections to the database.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit.
func WithMaxOpen(maxOpen int) Option {
	return func(cfg *serviceConfig) {
		cfg.MaxOpen = maxOpen
	}
}

// WithMaxIdleTime sets the maximum amount of time a connection may be reused.
//
// Expired connections may be closed lazily before reuse.
// Uint: milliseconds
func WithMaxIdleTime(maxIdelTime int) Option {
	return func(cfg *serviceConfig) {
		cfg.MaxIdleTime = maxIdelTime
	}
}
