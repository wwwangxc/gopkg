package example_test

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	// gopkg/redis will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/redis"
)

func ExampleNewClientProxy() {
	cli := redis.NewClientProxy("client_name",
		redis.WithClientDSN("dsn"),             // set dsn, default use database.client.dsn
		redis.WithClientMaxIdle(20),            // set max idel. default 2048
		redis.WithClientMaxActive(100),         // set max active. default 0
		redis.WithClientIdleTimeout(180000),    // set idle timeout. unit millisecond, default 180000
		redis.WithClientTimeout(1000),          // set command timeout. unit millisecond, default 1000
		redis.WithClientMaxConnLifetime(10000), // set max conn life time, default 0
		redis.WithClientWait(true),             // set wait
	)

	cli.Do(context.Background(), "GET", "foo")
	// do something...
}

func ExampleClientProxy_Conn() {
	c := redis.NewClientProxy("client_name").Conn()
	defer c.Close()
	c.Send("SET", "foo", "bar")
	c.Send("GET", "foo")
	c.Flush()
	c.Receive() // reply from SET
	c.Receive() // reply from GET
}

func ExampleClientProxy_Locker() {
	l := redis.NewClientProxy("client_name").Locker()

	// try lock
	// not block the current goroutine.
	// return uuid when the lock is acquired
	// return error when lock fail or lock not acquired
	// support reentrant unlock
	// support automatically renewal
	uuid, err := l.TryLock(context.Background(), "locker_key",
		redis.WithLockExpire(1000*time.Millisecond),
		redis.WithLockHeartbeat(500*time.Millisecond))

	if err != nil {

		// return ErrLockNotAcquired when lock not acquired
		if redis.IsLockNotAcquired(err) {
			fmt.Printf("lock not acquired\n")
			return
		}

		fmt.Printf("try lock fail. error: %v\n", err)
		return
	}

	defer func() {

		// return ErrLockNotExist if the key does not exist
		// return ErrNotOwnerOfKey if the uuid invalid
		// support reentrant unlock
		if err := l.Unlock(context.Background(), "locker_key", uuid); err != nil {
			fmt.Printf("unlock fail. error: %v\n", err)
		}
	}()

	// reentrant lock when uuid not empty
	// will block the current goroutine until lock is acquired when not reentrant lock
	_, err = l.Lock(context.Background(), "locker_key",
		redis.WithLockUUID(uuid),
		redis.WithLockExpire(1000*time.Millisecond),
		redis.WithLockHeartbeat(500*time.Millisecond))

	if err != nil {
		fmt.Printf("lock fail. error: %v\n", err)
		return
	}
}

func ExampleClientProxy_Fetcher() {
	obj := struct {
		FieldA string `json:"field_a"`
		FieldB int    `json:"field_b"`
	}{}

	callback := func() (interface{}, error) {
		// do something...
		return nil, nil
	}

	f := redis.NewClientProxy("client_name").Fetcher()

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
