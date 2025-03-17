package redis

import (
	"fmt"
	"os"
	"sync"

	"github.com/wwwangxc/gopkg/config"
)

var (
	serviceConfigMap = map[string]serviceConfig{}
	serviceConfigMu  sync.Mutex
)

func init() {
	_ = initAppConfig("./app.yaml")
}

// LoadConfig load config from file
func LoadConfig(path string) error {
	return initAppConfig(path)
}

type appConfig struct {
	Client struct {
		RedisCfg redisConfig     `yaml:"redis"`
		Service  []serviceConfig `yaml:"service"`
	} `yaml:"client"`
}

func (a *appConfig) getServiceConfigs() []serviceConfig {
	if a == nil {
		return []serviceConfig{}
	}

	clientConfigs := make([]serviceConfig, 0, len(a.Client.Service))
	for _, v := range a.Client.Service {
		v.Wait = a.Client.RedisCfg.Wait

		if v.MaxIdle == 0 {
			v.MaxIdle = a.Client.RedisCfg.MaxIdle
		}

		if v.MaxActive == 0 {
			v.MaxActive = a.Client.RedisCfg.MaxActive
		}

		if v.IdleTimeout == 0 {
			v.IdleTimeout = a.Client.RedisCfg.IdleTimeout
		}

		if v.MaxConnLifetime == 0 {
			v.MaxConnLifetime = a.Client.RedisCfg.MaxConnLifetime
		}

		if v.Timeout == 0 {
			v.Timeout = a.Client.RedisCfg.Timeout
		}

		clientConfigs = append(clientConfigs, v)
	}

	return clientConfigs
}

type redisConfig struct {
	MaxIdle         int  `yaml:"max_idle"`
	MaxActive       int  `yaml:"max_active"`
	MaxConnLifetime int  `yaml:"max_conn_lifetime"`
	IdleTimeout     int  `yaml:"idle_timeout"`
	Timeout         int  `yaml:"timeout"`
	Wait            bool `yaml:"wait"`
}

type serviceConfig struct {
	Name string `yaml:"name"`
	DSN  string `yaml:"dsn"`

	redisConfig `yaml:",inline"`
}

func initAppConfig(path string) error {
	c, err := loadAppConfig(path)
	if err != nil {
		return fmt.Errorf("config load fail. error:%v", err)
	}

	for _, v := range c.getServiceConfigs() {
		registerServiceConfig(v)
	}

	return nil
}

func loadAppConfig(path string) (*appConfig, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	configure, err := config.Load(path)
	if err != nil {
		return &appConfig{}, nil
	}

	c := &appConfig{}
	if err = configure.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("config unmarshal fail. error:%v", err)
	}

	return c, nil
}

func registerServiceConfig(c serviceConfig) {
	serviceConfigMu.Lock()
	defer serviceConfigMu.Unlock()
	serviceConfigMap[c.Name] = c
}

func getServiceConfig(name string) serviceConfig {
	serviceConfigMu.Lock()
	defer serviceConfigMu.Unlock()

	c, exist := serviceConfigMap[name]
	if !exist {
		c = serviceConfig{
			Name: name,
			redisConfig: redisConfig{
				MaxIdle:         2048,
				MaxActive:       0,
				IdleTimeout:     180000,
				MaxConnLifetime: 0,
				Timeout:         1000,
				Wait:            false,
			},
		}
		serviceConfigMap[name] = c
	}

	return c
}
