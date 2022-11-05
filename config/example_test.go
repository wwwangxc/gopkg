package config_test

import (
	"fmt"

	"github.com/wwwangxc/gopkg/config"
)

func Example() {
	// {
	//   "machine_id": 1,
	//   "local": {
	//     "env": "test",
	//     "login_name": "${LOGNAME}"
	//   }
	// }
	configure, err := config.Load("config.json", config.WithUnmarshaler("json"), config.WithWatchCallback(watch))
	if err != nil {
		fmt.Printf("load config fail. err:%v", err)
		return
	}
	_ = configure.GetUint32("machine_id", 999)      // return 1
	_ = configure.GetString("local.env", "default") // return "test"
	_ = configure.GetString("app", "default")       // return "default"

	// machine_id: 1
	// local:
	//   env: test
	//   login_name: ${LOGNAME}
	configure, err = config.Load("config.yaml", config.WithWatchCallback(watch))
	// or
	// configure, err = config.Load("config.yaml", config.WithUnmarshaler("yaml"), config.WithWatchCallback(watch))
	if err != nil {
		fmt.Printf("load config fail. err:%v", err)
		return
	}
	_ = configure.GetUint32("machine_id", 999)      // return 1
	_ = configure.GetString("local.env", "default") // return "test"
	_ = configure.GetString("app", "default")       // return "default"

	// machine_id=1
	//
	// [local]
	// env_name="test"
	// login_name="${LOGNAME}"
	configure, err = config.Load("config.toml", config.WithUnmarshaler("toml"), config.WithWatchCallback(watch))
	if err != nil {
		fmt.Printf("load config fail. err:%v", err)
		return
	}
	_ = configure.GetUint32("machine_id", 999)      // return 1
	_ = configure.GetString("local.env", "default") // return "test"
	_ = configure.GetString("app", "default")       // return "default"

}

func watch(configure config.Configure) {
	// callback when file changed
}
