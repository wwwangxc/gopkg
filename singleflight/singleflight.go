package singleflight

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

// Do helper for singleflight
//
// For details, see: https://pkg.go.dev/golang.org/x/sync/singleflight
func Do(ctx context.Context,
	key string, fn func(context.Context) (interface{}, error), opts ...Option) (interface{}, error) {
	opt := newOptions(opts...)

	var g singleflight.Group
	ch := g.DoChan(key, func() (interface{}, error) {
		go func() {
			time.Sleep(opt.expiresIn)
			g.Forget(key)
		}()

		return fn(ctx)
	})

	select {
	case <-ctx.Done():
		return nil, ErrTimeout
	case ret := <-ch:
		return ret.Val, ret.Err
	}
}
