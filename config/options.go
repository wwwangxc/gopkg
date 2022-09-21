package config

import (
	"github.com/wwwangxc/gopkg/config/unmarshaler"
)

// LoadOption load option for config
type LoadOption func(*configureImpl)

// WithUnmarshaler assign unmarshaler
func WithUnmarshaler(name string) LoadOption {
	return func(c *configureImpl) {
		c.unmarshaler = unmarshaler.Get(name)
	}
}

// WithWatchCallback with watch callback
func WithWatchCallback(callback func(Configure)) LoadOption {
	return func(c *configureImpl) {
		c.watchCallback = callback
	}
}

func withTest() LoadOption {
	return func(c *configureImpl) {
		c.watcher = nil
	}
}
