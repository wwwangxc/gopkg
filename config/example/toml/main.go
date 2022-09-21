package main

import (
	"fmt"

	"github.com/wwwangxc/gopkg/config"
)

func main() {
	// load config and serialize with toml
	configure, err := config.Load("config.toml", config.WithUnmarshaler("toml"), config.WithWatchCallback(watch))
	if err != nil {
		fmt.Printf("load config fail. err:%v", err)
		return
	}

	// read uint32 value
	fmt.Printf("machine_id: %d\n", configure.GetUint32("machine_id", 999))

	// read string value
	fmt.Printf("env_name: %s\n", configure.GetString("local.env_name", "env name"))

	// read string value
	// return default value when key not exist
	fmt.Printf("key_not_exist: %s\n", configure.GetString("key_not_exist", "this is default value"))

	c := &Config{}
	if err = configure.Unmarshal(c); err != nil {
		fmt.Printf("configure unmarshal fail. err:%v", err)
		return
	}

	fmt.Println("configure unmarshal success")
	fmt.Printf("Config: %+v", c)
}

func watch(configure config.Configure) {
	// callback when file changed
}

type Config struct {
	MachineID uint32 `toml:"machine_id"`
	Local     struct {
		EnvName   string `toml:"env_name"`
		LoginName string `toml:"login_name"`
	} `toml:"local"`
}
