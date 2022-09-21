package redis

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ClientOption redis client proxy option
type ClientOption func(*serviceConfig)

// WithClientDSN set dsn
func WithClientDSN(dsn string) ClientOption {
	return func(b *serviceConfig) {
		b.DSN = dsn
	}
}

// WithClientMaxIdle set max idle
//
// Maximum number of connections in the idle connection pool.
// Default 2048
func WithClientMaxIdle(maxIdle int) ClientOption {
	return func(b *serviceConfig) {
		b.MaxIdle = maxIdle
	}
}

// WithClientMaxActive set max active
//
// Maximum number of connections allocated by the pool at a given time.
// When zero, there is no limit on the number of connections in the pool.
// Default 0
func WithClientMaxActive(maxActive int) ClientOption {
	return func(b *serviceConfig) {
		b.MaxActive = maxActive
	}
}

// WithClientIdleTimeout set idle timeout
//
// Close connections after remaining idle for this duration. If the value
// is zero, then idle connections are not closed. Applications should set
// the timeout to a value less than the server's timeout.
// Unit millisecond, default 180000
func WithClientIdleTimeout(idleTimeout int) ClientOption {
	return func(b *serviceConfig) {
		b.IdleTimeout = idleTimeout
	}
}

// WithClientMaxConnLifetime set max conn lifetime
//
// Close connections older than this duration. If the value is zero, then
// the pool does not close connections based on age.
// Unit millisecond, default 0
func WithClientMaxConnLifetime(maxConnLifetime int) ClientOption {
	return func(b *serviceConfig) {
		b.MaxConnLifetime = maxConnLifetime
	}
}

// WithClientTimeout set timeout
//
// Write, read and connect timeout
// Unit millisecond, default 1000
func WithClientTimeout(timeout int) ClientOption {
	return func(b *serviceConfig) {
		b.Timeout = timeout
	}
}

// WithClientWait set wait
//
// If Wait is true and the pool is at the MaxActive limit, then Get() waits
// for a connection to be returned to the pool before returning.
func WithClientWait(wait bool) ClientOption {
	return func(b *serviceConfig) {
		b.Wait = wait
	}
}

// LockOptions distributed lock options
type LockOptions struct {
	// UUID of the lock
	// A non-null UUID indicates a reentrant lock.
	UUID string

	// Expire of the lock
	// Default 1000 millisecond
	Expire time.Duration

	// Heartbeat indicates the time interval for automatically renewal.
	// Heartbeat = 0 means the lock will not automatically renewal when it expires.
	// Default 0
	Heartbeat time.Duration

	// Retry indicates the time interval for retrying the acquire lock.
	// Default 1000 millisecond
	Retry time.Duration
}

func newLockOptions(opts ...LockOption) *LockOptions {
	options := defaultLockOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func defaultLockOptions() *LockOptions {
	return &LockOptions{
		UUID:      uuid.NewString(),
		Expire:    1000 * time.Millisecond,
		Heartbeat: 0,
		Retry:     1000 * time.Millisecond,
	}
}

// LockOption distributed lock option
type LockOption func(*LockOptions)

// WithLockUUID set uuid of the distributed lock
//
// A non-null UUID indicates a reentrant lock.
func WithLockUUID(uuid string) LockOption {
	return func(options *LockOptions) {
		options.UUID = uuid
	}
}

// WithLockExpire set lock expire
//
// default 1000 millisecond
func WithLockExpire(expire time.Duration) LockOption {
	return func(options *LockOptions) {
		options.Expire = expire
	}
}

// WithLockHeartbeat set heartbeat
//
// Heartbeat indicates the time interval for automatically renewal.
// Heartbeat = 0 means the lock will not automatically renewal when it expires.
// default 0
func WithLockHeartbeat(heartbeat time.Duration) LockOption {
	return func(options *LockOptions) {
		options.Heartbeat = heartbeat
	}
}

// WithLockRetry set retry
// Retry indicates the time interval for retrying the acquire lock.
// Default 1000 millisecond
func WithLockRetry(retry time.Duration) LockOption {
	return func(options *LockOptions) {
		options.Retry = retry
	}
}

// FetchOptions fetch options
type FetchOptions struct {
	Expire    time.Duration
	Callback  func() (interface{}, error)
	Marshal   func(v interface{}) ([]byte, error)
	Unmarshal func(data []byte, dest interface{}) error
}

func newFetchOptions(opts ...FetchOption) *FetchOptions {
	options := defaultFetchOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func defaultFetchOptions() *FetchOptions {
	return &FetchOptions{
		Expire:    1000 * time.Millisecond,
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
		Callback:  nil,
	}
}

// FetchOption fetch option
type FetchOption func(*FetchOptions)

// WithFetchCallback set fetch callback & expire option
//
// The callback function will be called if the key does not exist.
// Will cache the callback results into the key and set timeout.
// Default do nothing.
func WithFetchCallback(callback func() (interface{}, error), expire time.Duration) FetchOption {
	return func(options *FetchOptions) {
		options.Expire = expire
		options.Callback = callback
	}
}

// WithFetchMarshal set mashal function to fetcher
//
// The marshal function will be called before cache.
// Default use json.Marshal.
func WithFetchMarshal(marshal func(v interface{}) ([]byte, error)) FetchOption {
	return func(options *FetchOptions) {
		options.Marshal = marshal
	}
}

// WithFetchUnmarshal set unmarshal function to fetcher
//
// Default use json.Unmarshal.
func WithFetchUnmarshal(unmarshal func(data []byte, dest interface{}) error) FetchOption {
	return func(options *FetchOptions) {
		options.Unmarshal = unmarshal
	}
}
