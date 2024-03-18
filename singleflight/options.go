package singleflight

import "time"

type options struct {
	expiresIn time.Duration
}

func newOptions(opts ...Option) *options {
	opt := defaultOptions()
	for _, v := range opts {
		v(opt)
	}

	return opt
}

func defaultOptions() *options {
	return &options{
		expiresIn: time.Second,
	}
}

// Option of method do
type Option func(*options)

// WithExpiresIn set expiration time
//
// Default `time.Second`
func WithExpiresIn(expiresIn time.Duration) Option {
	return func(o *options) {
		o.expiresIn = expiresIn
	}
}
