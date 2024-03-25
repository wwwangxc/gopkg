package redis

import (
	"context"
	"errors"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/wwwangxc/gopkg/singleflight"
)

// FetcherProxy object fetcher
//
//go:generate mockgen -source=fetcher.go -destination=mockredis/fetcher_mock.go -package=mockredis
type FetcherProxy interface {

	// Fetch data and storing the result into the struct pointed at by dest.
	//
	// Use json decode
	Fetch(ctx context.Context, key string, dest interface{}, opts ...FetchOption) error
}

type fetcherImpl struct {
	name string
	opts []ClientOption
}

// NewFetcherProxy new object fetcher proxy
func NewFetcherProxy(name string, opts ...ClientOption) FetcherProxy {
	return &fetcherImpl{
		name: name,
		opts: opts,
	}
}

// Fetch data and storing the result into the struct pointed at by dest.
//
// Use json decode
func (f *fetcherImpl) Fetch(ctx context.Context, key string, dest interface{}, opts ...FetchOption) error {
	options := newFetchOptions(opts...)
	if options.ExpireSingleflight == 0 {
		return f.fetch(ctx, key, dest, options)
	}

	_, err := singleflight.Do(ctx, key, func(ctx context.Context) (interface{}, error) {
		return nil, f.fetch(ctx, key, dest, options)
	}, singleflight.WithExpiresIn(options.ExpireSingleflight))

	return err
}

func (f *fetcherImpl) fetch(ctx context.Context, key string, dest interface{}, options *FetchOptions) error {
	conn := f.getConn()
	defer func() {
		if err := conn.Close(); err != nil {
			logErrorf("conn close fail. error:%v", err)
		}
	}()

	data, err := Bytes(redigo.DoContext(conn, ctx, "GET", key))
	if err != nil && !errors.Is(redigo.ErrNil, err) {
		return err
	}

	if errors.Is(redigo.ErrNil, err) {

		if options.Callback == nil {
			return ErrKeyNotExist
		}

		val, err := options.Callback()
		if err != nil {
			return err
		}

		data, err = options.Marshal(val)
		if err != nil {
			return err
		}

		_, err = redigo.DoContext(conn, ctx, "PSETEX", key, options.Expire.Milliseconds(), data)
		if err != nil {
			return err
		}
	}

	return options.Unmarshal(data, dest)

}
func (f *fetcherImpl) getConn() redigo.Conn {
	return getRedisPool(f.name, f.opts...).Get()
}
