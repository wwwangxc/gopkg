package etcd

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"crypto/tls"

	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/wwwangxc/gopkg/config"
)

var (
	clientConfigMap = map[string]clientConfig{}
	clientConfigRW  sync.RWMutex
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
		ETCDCfg etcdConfig     `yaml:"etcd"`
		Service []clientConfig `yaml:"service"`
	} `yaml:"client"`
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

	c.registerClientConfig()
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

func (a *appConfig) registerClientConfig() {
	defaultTimeout := defaultClientConfig("").Timeout

	for _, v := range a.Client.Service {
		if v.Timeout < 1 {
			v.Timeout = a.Client.ETCDCfg.Timeout
		}

		if v.Timeout < 1 {
			v.Timeout = defaultTimeout
		}

		if v.TLSKeyPath == "" {
			v.TLSKeyPath = a.Client.ETCDCfg.TLSKeyPath
		}

		if v.TLSCertPath == "" {
			v.TLSCertPath = a.Client.ETCDCfg.TLSCertPath
		}

		if v.CACertPath == "" {
			v.CACertPath = a.Client.ETCDCfg.CACertPath
		}

		v.Endpoints = v.getEndpoints()
		v.Username = v.getUsername()
		v.Password = v.getPassword()

		registerClientConfig(v)
	}
}

type etcdConfig struct {
	Timeout     int    `yaml:"timeout"`
	TLSKeyPath  string `yaml:"tls_key"`
	TLSCertPath string `yaml:"tls_cert"`
	CACertPath  string `yaml:"ca_cert"`
}

type clientConfig struct {
	Name string `yaml:"name"`
	DSN  string `yaml:"dsn"`

	Username  string   `yaml:"-"`
	Password  string   `yaml:"_"`
	Endpoints []string `yaml:"-"`

	etcdConfig `yaml:",inline"`
}

func defaultClientConfig(name string) clientConfig {
	return clientConfig{
		Name: name,
		DSN:  "",
		etcdConfig: etcdConfig{
			Timeout: 1000,
		},
	}
}

func (s *clientConfig) getUsername() string {
	tmp := strings.Split(s.DSN, "@")
	if len(tmp) < 2 {
		return ""
	}

	auth := strings.Split(tmp[0], ":")
	if len(auth) < 2 {
		return ""
	}

	return auth[0]
}

func (s *clientConfig) getPassword() string {
	tmp := strings.Split(s.DSN, "@")
	if len(tmp) < 2 {
		return ""
	}

	auth := strings.Split(tmp[0], ":")
	if len(auth) < 2 {
		return ""
	}

	return auth[1]
}

func (s *clientConfig) getEndpoints() []string {
	tmp := strings.Split(s.DSN, "@")

	if len(tmp) < 1 {
		return []string{}
	}

	if len(tmp) == 1 {
		return strings.Split(tmp[0], ",")
	}

	var endpoints []string
	tmp = strings.Split(tmp[1], ",")
	for _, v := range tmp {
		endpoints = append(endpoints, fmt.Sprintf("http://%s", strings.TrimPrefix(v, "http://")))
	}

	return endpoints
}

func (s *clientConfig) clientConfig() (*clientv3.Config, error) {
	tlsConfig, err := s.tlsConfig()
	if err != nil {
		return nil, err
	}

	return &clientv3.Config{
		Endpoints:   s.Endpoints,
		DialTimeout: time.Millisecond * time.Duration(s.Timeout),
		Username:    s.Username,
		Password:    s.Password,
		TLS:         tlsConfig,
	}, nil
}

func (s *clientConfig) tlsConfig() (*tls.Config, error) {
	if s.TLSKeyPath == "" || s.TLSCertPath == "" || s.CACertPath == "" {
		return nil, nil
	}

	tlsInfo := &transport.TLSInfo{
		TrustedCAFile: s.CACertPath,
		CertFile:      s.TLSCertPath,
		KeyFile:       s.TLSKeyPath,
	}

	return tlsInfo.ClientConfig()
}

func registerClientConfig(c clientConfig) {
	clientConfigRW.Lock()
	defer clientConfigRW.Unlock()
	clientConfigMap[c.Name] = c
}

func getClientConfig(name string) clientConfig {
	clientConfigRW.RLock()
	defer clientConfigRW.RUnlock()

	c, exist := clientConfigMap[name]
	if !exist {
		clientConfigMap[name] = defaultClientConfig(name)
	}

	return c
}
