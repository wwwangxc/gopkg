package config

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cast"

	"github.com/wwwangxc/gopkg/config/unmarshaler"
)

// Load load config
func Load(path string, opts ...LoadOption) (Configure, error) {
	return defaultLoader.Load(path, opts...)
}

// Configure ...
//go:generate mockgen -source=config.go -destination=mockconfig/config_mock.go -package=mockconfig
type Configure interface {

	// Unmarshal unmarshal config raw data
	Unmarshal(interface{}) error

	// IsExist check the key exist
	IsExist(string) bool

	// Get get value by key
	//
	// return defaultVal, when key not exist
	// k support key1.key2.key3
	Get(string, interface{}) interface{}

	// GetString get string value by key
	//
	// return defaultVal, when key not exist
	// k support key1.key2.key3
	GetString(string, string) string

	// GetBool get bool value by key
	//
	// return defaultVal, when key not exist
	// k support key1.key2.key3
	GetBool(string, bool) bool

	// GetInt get int value by key
	//
	// return defaultVal, when key not exist
	// k support key1.key2.key3
	GetInt(string, int) int

	// GetInt32 get int32 value by key
	//
	// return defaultVal, when key not exist
	// k support key1.key2.key3
	GetInt32(string, int32) int32

	// GetInt64 get int64 value by key
	//
	// return defaultVal, when key not exist
	// k support key1.key2.key3
	GetInt64(string, int64) int64

	// GetUint get uint value by key
	//
	// return defaultVal, when key not exist
	// k support key1.key2.key3
	GetUint(string, uint) uint

	// GetUint32 get uint32 value by key
	//
	// return defaultVal, when key not exist
	// k support key1.key2.key3
	GetUint32(string, uint32) uint32

	// GetUint64 get uint64 value by key
	//
	// return defaultVal, when key not exist
	// k support key1.key2.key3
	GetUint64(string, uint64) uint64

	// GetFloat32 get float32 value by key
	//
	// return defaultVal, when key not exist
	// k support key1.key2.key3
	GetFloat32(string, float32) float32

	// GetFloat64 get float64 value by key
	//
	// return defaultVal, when key not exist
	// k support key1.key2.key3
	GetFloat64(string, float64) float64
}

// configureImpl ...
type configureImpl struct {
	path            string
	rawData         []byte
	unmarshaledData map[string]interface{}

	rw            sync.RWMutex
	watchCallback func(Configure)
	unmarshaler   unmarshaler.Unmarshaler
	watcher       *fsnotify.Watcher
}

func defaultConfigure(path string) *configureImpl {
	c := &configureImpl{
		path:        path,
		unmarshaler: &unmarshaler.YAML{},
	}

	var err error
	c.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		logErrorf("%s: new file watcher fail. err:%v\n", packageName, err)
	}

	return c
}

// Unmarshal unmarshal config raw data
func (c *configureImpl) Unmarshal(out interface{}) error {
	if c.unmarshaler == nil {
		return ErrUnmarshalerNotExist
	}

	c.rw.RLock()
	defer c.rw.RUnlock()

	return c.unmarshaler.Unmarshal(c.rawData, out)
}

// IsExist check the key exist
func (c *configureImpl) IsExist(k string) bool {
	_, err := c.get(k)
	return err == nil
}

// Get get value by key
//
// return defaultVal, when key not exist
// k support key1.key2.key3
func (c *configureImpl) Get(k string, defaultVal interface{}) interface{} {
	val, err := c.get(k)
	if err != nil {
		return defaultVal
	}

	return val
}

// GetString get string value by key
//
// return defaultVal, when key not exist
// k support key1.key2.key3
func (c *configureImpl) GetString(k string, defaultVal string) string {
	return cast.ToString(c.getWithDefaultVal(k, defaultVal))
}

// GetBool get bool value by key
//
// return defaultVal, when key not exist
// k support key1.key2.key3
func (c *configureImpl) GetBool(k string, defaultVal bool) bool {
	return cast.ToBool(c.getWithDefaultVal(k, defaultVal))
}

// GetInt get int value by key
//
// return defaultVal, when key not exist
// k support key1.key2.key3
func (c *configureImpl) GetInt(k string, defaultVal int) int {
	return cast.ToInt(c.getWithDefaultVal(k, defaultVal))
}

// GetInt32 get int32 value by key
//
// return defaultVal, when key not exist
// k support key1.key2.key3
func (c *configureImpl) GetInt32(k string, defaultVal int32) int32 {
	return cast.ToInt32(c.getWithDefaultVal(k, defaultVal))
}

// GetInt64 get int64 value by key
//
// return defaultVal, when key not exist
// k support key1.key2.key3
func (c *configureImpl) GetInt64(k string, defaultVal int64) int64 {
	return cast.ToInt64(c.getWithDefaultVal(k, defaultVal))
}

// GetUint get uint value by key
//
// return defaultVal, when key not exist
// k support key1.key2.key3
func (c *configureImpl) GetUint(k string, defaultVal uint) uint {
	return cast.ToUint(c.getWithDefaultVal(k, defaultVal))
}

// GetUint32 get uint32 value by key
//
// return defaultVal, when key not exist
// k support key1.key2.key3
func (c *configureImpl) GetUint32(k string, defaultVal uint32) uint32 {
	return cast.ToUint32(c.getWithDefaultVal(k, defaultVal))
}

// GetUint64 get uint64 value by key
//
// return defaultVal, when key not exist
// k support key1.key2.key3
func (c *configureImpl) GetUint64(k string, defaultVal uint64) uint64 {
	return cast.ToUint64(c.getWithDefaultVal(k, defaultVal))
}

// GetFloat32 get float32 value by key
//
// return defaultVal, when key not exist
// k support key1.key2.key3
func (c *configureImpl) GetFloat32(k string, defaultVal float32) float32 {
	return cast.ToFloat32(c.getWithDefaultVal(k, defaultVal))
}

// GetFloat64 get float64 value by key
//
// return defaultVal, when key not exist
// k support key1.key2.key3
func (c *configureImpl) GetFloat64(k string, defaultVal float64) float64 {
	return cast.ToFloat64(c.getWithDefaultVal(k, defaultVal))
}

func (c *configureImpl) getWithDefaultVal(k string, defaultVal interface{}) interface{} {
	data, err := c.get(k)
	if err != nil {
		return defaultVal
	}

	switch defaultVal.(type) {
	case string:
		_, err = cast.ToStringE(data)
	case bool:
		_, err = cast.ToBoolE(data)
	case int:
		_, err = cast.ToIntE(data)
	case int32:
		_, err = cast.ToInt32E(data)
	case int64:
		_, err = cast.ToInt64E(data)
	case uint:
		_, err = cast.ToUintE(data)
	case uint32:
		_, err = cast.ToUint32E(data)
	case uint64:
		_, err = cast.ToUint64E(data)
	case float64:
		_, err = cast.ToFloat64E(data)
	case float32:
		_, err = cast.ToFloat32E(data)
	default:
		return defaultVal
	}

	if err != nil {
		return defaultVal
	}

	return data
}

func (c *configureImpl) get(k string) (interface{}, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	subkeys := strings.Split(k, ".")
	val, exist := fetchFromMap(c.unmarshaledData, subkeys)
	if !exist {
		return nil, ErrConfigNotExist
	}

	return val, nil
}

// Load ...
func (c *configureImpl) Load() error {
	if c.unmarshaler == nil {
		return ErrUnmarshalerNotExist
	}

	data, err := ioutil.ReadFile(c.path)
	if err != nil {
		return fmt.Errorf("read file fail. err:%w", err)
	}

	unmarshaledData := map[string]interface{}{}
	if err = c.unmarshaler.Unmarshal(data, &unmarshaledData); err != nil {
		return fmt.Errorf("unmarshal fail. err:%w", err)
	}

	c.rw.Lock()
	defer c.rw.Unlock()

	c.rawData = data
	c.unmarshaledData = unmarshaledData

	return nil
}

// Reload ...
func (c *configureImpl) Reload() {
	if c.unmarshaler == nil {
		return
	}

	data, err := ioutil.ReadFile(c.path)
	if err != nil {
		logErrorf("%s: reload file fail. err:%v\n", packageName, err)
		return
	}

	unmarshalData := map[string]interface{}{}
	if err = c.unmarshaler.Unmarshal(data, &unmarshalData); err != nil {
		logErrorf("%s: reload file unmarshal fail. err:%v\n", packageName, err)
		return
	}

	c.rw.Lock()
	defer c.rw.Unlock()

	c.rawData = data
	c.unmarshaledData = unmarshalData
}

func (c *configureImpl) watch(callback func(*configureImpl)) {
	if c.watcher == nil {
		return
	}

	go func() {
		for event := range c.watcher.Events {
			if event.Op&fsnotify.Write != fsnotify.Write {
				logInfo("%s: ignore file event:%s. file:%s\n", packageName, event.Name, c.path)
				continue
			}

			c.Reload()

			if callback != nil {
				callback(c)
			}

			if c.watchCallback != nil {
				go c.watchCallback(c)
			}
		}

		logInfo("%s: break file watch. file:%s\n", packageName, c.path)
	}()
}

func fetchFromMap(m map[string]interface{}, subkeys []string) (interface{}, bool) {
	if len(subkeys) == 0 {
		return nil, false
	}

	data, exist := m[subkeys[0]]
	if !exist {
		return nil, false
	}

	if len(subkeys) == 1 {
		return data, true
	}

	switch val := data.(type) {
	case map[interface{}]interface{}:
		return fetchFromMap(cast.ToStringMap(val), subkeys[1:])
	case map[string]interface{}:
		return fetchFromMap(val, subkeys[1:])
	default:
		return nil, false
	}
}
