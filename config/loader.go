package config

import (
	"fmt"
	"sync"
)

var (
	defaultLoader = newLoader()
)

type loader struct {
	m  map[string]Configure
	rw sync.RWMutex
}

func newLoader() *loader {
	return &loader{
		m: map[string]Configure{},
	}
}

// Load load and cache config
func (l *loader) Load(path string, opts ...LoadOption) (Configure, error) {
	c := defaultConfigure(path)
	for _, opt := range opts {
		opt(c)
	}

	if c.unmarshaler == nil {
		return nil, ErrUnmarshalerNotExist
	}

	key := fmt.Sprintf("%s:%s", path, c.unmarshaler.Name())
	l.rw.RLock()
	tmp, exist := l.m[key]
	l.rw.RUnlock()
	if exist {
		return tmp, nil
	}

	if err := c.Load(); err != nil {
		return nil, fmt.Errorf("%s: config load fail. err:%w", packageName, err)
	}

	l.rw.Lock()
	l.m[key] = c
	l.rw.Unlock()

	c.watch(func(c *configureImpl) {
		l.rw.Lock()
		defer l.rw.Unlock()
		l.m[key] = c
	})

	return c, nil
}
