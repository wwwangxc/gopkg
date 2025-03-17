package orm

import (
	"fmt"
	"os"
	"sync"

	"gorm.io/gorm"

	"github.com/wwwangxc/gopkg/config"
	"github.com/wwwangxc/gopkg/orm/driver"
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
		MySQL      dbConfig        `yaml:"mysql"`
		PostgreSQL dbConfig        `yaml:"postgresql"`
		SQLite     dbConfig        `yaml:"sqlite"`
		SQLServer  dbConfig        `yaml:"sqlserver"`
		Clickhouse dbConfig        `yaml:"clickhouse"`
		Service    []serviceConfig `yaml:"service"`
	} `yaml:"client"`
}

func (a *appConfig) getServiceConfigs() []serviceConfig {
	if a == nil {
		return []serviceConfig{}
	}

	serviceConfigs := make([]serviceConfig, 0, len(a.Client.Service))
	for _, v := range a.Client.Service {
		dbCfg := a.getDBConfig(v.Driver)
		if v.MaxIdle == 0 {
			v.MaxIdle = dbCfg.MaxIdle
		}

		if v.MaxOpen == 0 {
			v.MaxOpen = dbCfg.MaxOpen
		}

		if v.MaxIdleTime == 0 {
			v.MaxIdleTime = dbCfg.MaxIdleTime
		}

		serviceConfigs = append(serviceConfigs, v)
	}

	return serviceConfigs
}

func (a *appConfig) getDBConfig(driverName string) dbConfig {
	switch driverName {
	case driver.NameMySQL:
		return a.Client.MySQL
	case driver.NamePostgreSQL, driver.NamePostgreSQLSimple:
		return a.Client.PostgreSQL
	case driver.NameSQLite:
		return a.Client.SQLite
	case driver.NameSQLServer:
		return a.Client.SQLServer
	case driver.NameClickhouse:
		return a.Client.Clickhouse
	default:
		return dbConfig{}
	}
}

type dbConfig struct {
	MaxIdle     int `yaml:"max_idle"`
	MaxOpen     int `yaml:"max_open"`
	MaxIdleTime int `yaml:"max_idle_time"`
}

type serviceConfig struct {
	Name     string `yaml:"name"`
	DSN      string `yaml:"dsn"`
	Driver   string `yaml:"driver"`
	dbConfig `yaml:",inline"`

	gormConfig *gorm.Config `yaml:"-"`
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
			Name:   name,
			Driver: "mysql",
		}
		serviceConfigMap[name] = c
	}

	return c
}
