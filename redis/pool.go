package redis

import (
	"context"
	"net"
	"sync"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

var (
	pools   = map[string]*redigo.Pool{}
	poolsRW sync.RWMutex
)

func getRedisPool(name string, opts ...ClientOption) *redigo.Pool {
	poolsRW.RLock()
	pool, ok := pools[name]
	poolsRW.RUnlock()
	if ok {
		return pool
	}

	cfg := getServiceConfig(name)
	for _, opt := range opts {
		opt(&cfg)
	}

	return newRedisPool(&cfg)
}

func newRedisPool(cfg *serviceConfig) *redigo.Pool {
	poolsRW.Lock()
	defer poolsRW.Unlock()

	pool, ok := pools[cfg.Name]
	if ok {
		return pool
	}

	timeout := time.Duration(cfg.Timeout) * time.Millisecond
	pool = &redigo.Pool{
		MaxIdle:         cfg.MaxIdle,
		MaxActive:       cfg.MaxActive,
		IdleTimeout:     time.Duration(cfg.IdleTimeout) * time.Millisecond,
		MaxConnLifetime: time.Duration(cfg.MaxConnLifetime) * time.Millisecond,
		Dial: func() (redigo.Conn, error) {
			dialOpts := []redigo.DialOption{
				redigo.DialWriteTimeout(timeout),
				redigo.DialReadTimeout(timeout),
				redigo.DialConnectTimeout(timeout),
				redigo.DialContextFunc(func(ctx context.Context, network, addr string) (net.Conn, error) {
					dialer := &net.Dialer{
						Timeout: timeout,
					}
					return dialer.DialContext(ctx, network, addr)
				}),
			}

			c, err := redigo.DialURL(cfg.DSN, dialOpts...)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		Wait: cfg.Wait,
	}

	pools[cfg.Name] = pool
	return pool
}
