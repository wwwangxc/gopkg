package httpx

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/wwwangxc/gopkg/config"
)

var (
	clientConfigRW  sync.RWMutex
	clientConfigMap = map[string]clientConfig{}
)

func init() {
	_ = initAppConfig("./app.yaml")
}

// LoadConfig load config from file
func LoadConfig(path string) error {
	return initAppConfig(path)
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

type appConfig struct {
	Client struct {
		HTTPCfg httpConfig     `yaml:"http"`
		Service []clientConfig `yaml:"service"`
	} `yaml:"client"`
}

func (a *appConfig) registerClientConfig() {
	httpConfigDefault := defaultHTTPConfig()

	for _, v := range a.Client.Service {
		if len(v.httpConfig.Header) == 0 {
			v.httpConfig.Header = a.Client.HTTPCfg.Header
		}

		if v.httpConfig.TransportCfg == nil {
			v.httpConfig.TransportCfg = a.Client.HTTPCfg.TransportCfg
		}

		if v.httpConfig.TransportCfg == nil {
			v.httpConfig.TransportCfg = httpConfigDefault.TransportCfg
		}

		registerClientConfig(v)
	}
}

type clientConfig struct {
	Name    string `yaml:"name"`
	DSN     string `yaml:"dsn"`
	Timeout int64  `yaml:"timeout"`

	httpConfig `yaml:",inline"`
}

func (s *clientConfig) toHTTPClient() *http.Client {
	return &http.Client{
		Transport: s.getTransport(),
	}
}

func defaultClientConfig(name string) clientConfig {
	return clientConfig{
		Name:       name,
		DSN:        "",
		Timeout:    3000,
		httpConfig: defaultHTTPConfig(),
	}
}

type httpConfig struct {
	Header       map[string]string    `yaml:"header"`
	TransportCfg *httpTransportConfig `yaml:"transport"`
}

func defaultHTTPConfig() httpConfig {
	transportCfg := defaultHTTPTransportConfig()
	return httpConfig{
		Header:       map[string]string{},
		TransportCfg: &transportCfg,
	}
}

func (s *httpConfig) getTransport() *http.Transport {
	if s == nil || s.TransportCfg == nil {
		return nil
	}

	dialer := &net.Dialer{
		Timeout:   s.TransportCfg.Dial.Timeout,
		KeepAlive: s.TransportCfg.Dial.KeepAlive,
	}

	return &http.Transport{
		MaxIdleConns:          s.TransportCfg.MaxIdleConns,
		MaxIdleConnsPerHost:   s.TransportCfg.MaxIdleConnsPerHost,
		MaxConnsPerHost:       s.TransportCfg.MaxConnsPerHost,
		IdleConnTimeout:       s.TransportCfg.IdleConnTimeout,
		TLSHandshakeTimeout:   s.TransportCfg.TLSHandshakeTimeout,
		ExpectContinueTimeout: s.TransportCfg.ExpectContinueTimeout,
		ResponseHeaderTimeout: s.TransportCfg.ResponseHeaderTimeout,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
	}
}

type httpTransportConfig struct {
	MaxIdleConns        int           `yaml:"max_idle_conns"`
	MaxIdleConnsPerHost int           `yaml:"max_idle_conns_per_host"`
	MaxConnsPerHost     int           `yaml:"max_conns_per_host"`
	IdleConnTimeout     time.Duration `yaml:"idle_conn_timeout"`

	TLSHandshakeTimeout   time.Duration `yaml:"tls_handshake_timeout"`
	ExpectContinueTimeout time.Duration `yaml:"expect_continue_timeout"`
	ResponseHeaderTimeout time.Duration `yaml:"response_header_timeout"`

	Dial struct {
		Timeout   time.Duration `yaml:"timeout"`
		KeepAlive time.Duration `yaml:"keep_alive"`
	} `yaml:"dial"`
}

func defaultHTTPTransportConfig() httpTransportConfig {
	return httpTransportConfig{
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   50,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		Dial: struct {
			Timeout   time.Duration `yaml:"timeout"`
			KeepAlive time.Duration `yaml:"keep_alive"`
		}{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		},
	}
}

func registerClientConfig(c clientConfig) {
	clientConfigRW.Lock()
	defer clientConfigRW.Unlock()
	clientConfigMap[c.Name] = c
}

func getClientConfig(name string) clientConfig {
	clientConfigRW.RLock()
	if c, exist := clientConfigMap[name]; exist {
		clientConfigRW.RUnlock()
		return c
	}
	clientConfigRW.RUnlock()

	clientConfigRW.Lock()
	defer clientConfigRW.Unlock()
	if c, exist := clientConfigMap[name]; exist {
		return c
	}

	c := defaultClientConfig(name)
	clientConfigMap[name] = c
	return c
}
