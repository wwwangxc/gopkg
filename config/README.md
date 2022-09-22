# gopkg/config

[![Go Report Card](https://goreportcard.com/badge/github.com/wwwangxc/gopkg/config)](https://goreportcard.com/report/github.com/wwwangxc/gopkg/config)
[![GoDoc](https://pkg.go.dev/badge/github.com/wwwangxc/gopkg/config?status.svg)](https://pkg.go.dev/github.com/wwwangxc/gopkg/config)
[![OSCS Status](https://www.oscs1024.com/platform/badge/wwwangxc/gopkg.svg?size=small)](https://www.murphysec.com/dr/c1TuOdJ62DzT0agLwg)

gopkg/config is an componentized config plugin.

It provides an easy way to load configuration files.

## Install

```sh
go get github.com/wwwangxc/gopkg/config
```

## Quick Start

**main.go**

```go
package main

import (
        "github.com/wwwangxc/gopkg/config"
)

func main() {
        // load local config file.
        // config.WithUnmarshaler("yaml"): serialize with yaml.
        // config.WithWatchCallback(watch): watch config file, callback watch funcation when config file changed.
        configure, err := config.Load("./config.yaml", config.WithUnmarshaler("yaml"), config.WithWatchCallback(watch))

        // the default unmarshaler is yaml
        configure, err = config.Load("./config.yaml", config.WithWatchCallback(watch))

        // serialize config file with toml
        configure, err = config.Load("./config.toml", config.WithUnmarshaler("toml"))

        // serialize config file with json
        configure, err = config.Load("./config.json", config.WithUnmarshaler("json"))

        // read string value
        configure.GetString("app.env_name", "default")

        // read bool value
        configure.GetBool("app.debug", false)

        // read uint32 value
        configure.GetUint32("machine_id", 1)

        // unmarshal raw data to Config struct
        c := &Config{}
        err = configure.Unmarshal(c)
}

func watch(configure config.Configure) {
        // do something ...
}

type Config struct {
        MachineID uint32 `yaml:"machine_id" toml:"machine_id" json:"machine_id"`
        APP       struct{
                EnvName string `yaml:"env_name" toml:"env_name" json:"env_name"`
                Debug   bool `yaml:"debug" toml:"debug" json:"debug"`
        } `yaml:"app" toml:"app" json:"app"`
}
```

**config.yaml**

```yaml
machine_id: 1
app:
  env_name: ${ENV}
  debug: true
```

**config.toml**

```toml
machine_id = 1

[app]
env_name = "${ENV}"
debug = true
```

**config.json**

```json
{
  "machine_id": 1,
  "app": {
    "env_name": "${ENV}",
    "debug": true
  }
}
```

## How To Mock

```go
package tests

import (
        "testing"

        "github.com/agiledragon/gomonkey"
        "github.com/golang/mock/gomock"

        "github.com/wwwangxc/gopkg/config/mockconfig"
)

func TestLoad(t *testing.T) {
        ctrl := gomock.NewController(t)
        defer ctrl.Finish()
        
        mockConfigure := mockconfig.NewMockConfigure(ctrl)
        mockConfigure.EXPECT().GetString(gomock.Any(), gomock.Any()).
            Return("mocked value").AnyTimes()
        
        dispatch := gomonkey.ApplyFunc(Load,
       	        func(string, ...LoadOption) (Configure, error) {
       	                return mockConfigure, nil
       	        })
        defer dispatch.Reset()

        // test cases
}
```
