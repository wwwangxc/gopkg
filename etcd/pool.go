package etcd

import (
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	cliPool   = map[string]*clientv3.Client{}
	cliPoolRW sync.RWMutex
)

func getETCDClient(name string, opts ...ClientOption) (*clientv3.Client, error) {
	cliPoolRW.RLock()
	cli, ok := cliPool[name]
	cliPoolRW.RUnlock()
	if ok {
		return cli, nil
	}

	cfg := getClientConfig(name)
	for _, opt := range opts {
		opt(&cfg)
	}

	return newETCDClient(cfg)
}

func newETCDClient(cfg clientConfig) (*clientv3.Client, error) {
	cliPoolRW.Lock()
	defer cliPoolRW.Unlock()

	cli, ok := cliPool[cfg.Name]
	if ok {
		return cli, nil
	}

	clientConfig, err := cfg.clientConfig()
	if err != nil {
		return nil, err
	}

	cli, err = clientv3.New(*clientConfig)
	if err != nil {
		return nil, err
	}

	cliPool[cfg.Name] = cli
	return cli, nil
}
