package httpx

import (
	"strings"
	"sync"
	"time"

	"resty.dev/v3"
)

var (
	cliPoolRW sync.RWMutex
	cliPool   = map[string]*resty.Client{}
)

func getOrCreateClient(name string, opts ...ClientOption) *resty.Client {
	cliPoolRW.RLock()
	cli, ok := cliPool[name]
	cliPoolRW.RUnlock()
	if ok {
		return cli
	}

	cliPoolRW.Lock()
	defer cliPoolRW.Unlock()

	cli, ok = cliPool[name]
	if ok {
		return cli
	}

	cli = newClient(name, opts...)
	if cli == nil {
		return nil
	}

	cliPool[name] = cli
	return cli
}

func newClient(name string, opts ...ClientOption) *resty.Client {
	config := getClientConfig(name)

	cli := resty.NewWithClient(config.toHTTPClient()).
		SetHeaders(config.Header).
		SetTimeout(time.Duration(config.Timeout) * time.Millisecond)

	if config.DSN != "" {
		rr, err := resty.NewRoundRobin(strings.Split(config.DSN, ",")...)
		if err != nil {
			panic(err)
		}
		cli = cli.SetLoadBalancer(rr)
	}

	for _, opt := range opts {
		opt(cli)
	}

	return cli
}
