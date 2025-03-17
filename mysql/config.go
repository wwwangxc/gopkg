package mysql

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
		MySQLConfig mysqlConfig     `yaml:"mysql"`
		Service     []serviceConfig `yaml:"service"`
	} `yaml:"client"`
}

func (a *appConfig) getServiceConfigs() []serviceConfig {
	if a == nil {
		return []serviceConfig{}
	}

	serviceConfigs := make([]serviceConfig, 0, len(a.Client.Service))
	for _, v := range a.Client.Service {
		if v.MaxIdle == 0 {
			v.MaxIdle = a.Client.MySQLConfig.MaxIdle
		}

		if v.MaxOpen == 0 {
			v.MaxOpen = a.Client.MySQLConfig.MaxOpen
		}

		if v.MaxIdleTime == 0 {
			v.MaxIdleTime = a.Client.MySQLConfig.MaxIdleTime
		}

		serviceConfigs = append(serviceConfigs, v)
	}

	return serviceConfigs
}

type mysqlConfig struct {
	MaxIdle     int `yaml:"max_idle"`
	MaxOpen     int `yaml:"max_open"`
	MaxIdleTime int `yaml:"max_idle_time"`
}

type serviceConfig struct {
	Name string `yaml:"name"`
	DSN  string `yaml:"dsn"`

	mysqlConfig `yaml:",inline"`
}

func initAppConfig(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

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
		}
		serviceConfigMap[name] = c
	}

	return c
}
