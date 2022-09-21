package example

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	// gopkg/redis will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/redis"
)

func ExampleNewFetcherProxy() {
	obj := struct {
		FieldA string `json:"field_a"`
		FieldB int    `json:"field_b"`
	}{}

	callback := func() (interface{}, error) {
		// do something...
		return nil, nil
	}

	f := redis.NewFetcherProxy("client_name",
		redis.WithClientDSN("dsn"),             // set dsn, default use database.client.dsn
		redis.WithClientMaxIdle(20),            // set max idel. default 2048
		redis.WithClientMaxActive(100),         // set max active. default 0
		redis.WithClientIdleTimeout(180000),    // set idle timeout. unit millisecond, default 180000
		redis.WithClientTimeout(1000),          // set command timeout. unit millisecond, default 1000
		redis.WithClientMaxConnLifetime(10000), // set max conn life time, default 0
		redis.WithClientWait(true),             // set wait
	)

	// fetch object
	err := f.Fetch(context.Background(), "fetcher_key", &obj,
		redis.WithFetchCallback(callback, 1000*time.Millisecond),
		redis.WithFetchUnmarshal(json.Unmarshal),
		redis.WithFetchMarshal(json.Marshal))

	if err != nil {
		fmt.Printf("fetch fail. error: %v\n", err)
		return
	}
}
