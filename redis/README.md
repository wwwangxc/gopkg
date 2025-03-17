# gopkg Redis Package

[![Go Report Card](https://goreportcard.com/badge/github.com/wwwangxc/gopkg/redis)](https://goreportcard.com/report/github.com/wwwangxc/gopkg/redis)
[![GoDoc](https://pkg.go.dev/badge/github.com/wwwangxc/gopkg/redis?status.svg)](https://pkg.go.dev/github.com/wwwangxc/gopkg/redis)
[![OSCS Status](https://www.oscs1024.com/platform/badge/wwwangxc/gopkg.svg?size=small)](https://www.murphysec.com/dr/c1TuOdJ62DzT0agLwg)

gopkg/redis is an componentized redis package.

It provides:

- An easy way to configre and manage redis client.
- Lock handler.
- Object fetcher.

Based on [gomodule/redigo](https://github.com/gomodule/redigo).

## Install

```sh
go get github.com/wwwangxc/gopkg/redis
```

## Quick Start

### Client Proxy

```go
package main

import (
        "context"
        "fmt"

        // gopkg/redis will automatically read configuration
        // files (./app.yaml) when package loaded
        "github.com/wwwangxc/gopkg/redis"
)

func main() {
        // default config is read from ./app.yaml
        // you can also set config by code
        // redis.LoadConfig("./app.yaml")

        cli := redis.NewClientProxy("client_name",
                redis.WithDSN("dsn"),             // set dsn, default use database.client.dsn
                redis.WithMaxIdle(20),            // set max idel. default 2048
                redis.WithMaxActive(100),         // set max active. default 0
                redis.WithIdleTimeout(180000),    // set idle timeout. unit millisecond, default 180000
                redis.WithTimeout(1000),          // set command timeout. unit millisecond, default 1000
                redis.WithMaxConnLifetime(10000), // set max conn life time, default 0
                redis.WithWait(true),             // set wait
        )

        // Exec GET command
        cli.Do(context.Background(), "GET", "foo")

        // Pipeline
        // get a redis connection
        c := cli.Conn()
        defer c.Close()

        c.Send("SET", "foo", "bar")
        c.Send("GET", "foo")
        c.Flush()
        c.Receive()           // reply from SET
        v, err := c.Receive() // reply from GET
        fmt.Sprintf("reply: %s", v)
        fmt.Sprintf("error: %v", err)
}
```

### Locker Proxy

```go
package main

import (
        "context"
        "fmt"

        // gopkg/redis will automatically read configuration
        // files (./app.yaml) when package loaded
        "github.com/wwwangxc/gopkg/redis"
)

func main() {
        // cli := redis.NewClientProxy("client_name").Locker()
        cli := redis.NewLockerProxy("client_name",
                redis.WithDSN("dsn"),             // set dsn, default use database.client.dsn
                redis.WithMaxIdle(20),            // set max idel. default 2048
                redis.WithMaxActive(100),         // set max active. default 0
                redis.WithIdleTimeout(180000),    // set idle timeout. unit millisecond, default 180000
                redis.WithTimeout(1000),          // set command timeout. unit millisecond, default 1000
                redis.WithMaxConnLifetime(10000), // set max conn life time, default 0
                redis.WithWait(true),             // set wait
        )

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

        f := func() error {
                // do something...
                return nil
        }

        // try get lock first and call f() when lock acquired. Unlock will be performed
        // regardless of whether the f reports an error or not.
        if err := l.LockAndCall(context.Background(), "locker_key", f); err != nil {
                fmt.Printf("lock and call fail. error: %v\n", err)
                return
        }
}
```

### Fetcher Proxy

```go
package main

import (
        "context"
        "fmt"

        // gopkg/redis will automatically read configuration
        // files (./app.yaml) when package loaded
        "github.com/wwwangxc/gopkg/redis"
)

func main() {
        // f := redis.NewClientProxy("client_name").Fetcher()
        f := redis.NewFetcherProxy("client_name",
                redis.WithDSN("dsn"),             // set dsn, default use database.client.dsn
                redis.WithMaxIdle(20),            // set max idel. default 2048
                redis.WithMaxActive(100),         // set max active. default 0
                redis.WithIdleTimeout(180000),    // set idle timeout. unit millisecond, default 180000
                redis.WithTimeout(1000),          // set command timeout. unit millisecond, default 1000
                redis.WithMaxConnLifetime(10000), // set max conn life time, default 0
                redis.WithWait(true),             // set wait
        )

        obj := struct {
                FieldA string `json:"field_a"`
                FieldB int    `json:"field_b"`
        }{}
        
        callback := func() (interface{}, error) {
                // do something...
                return nil, nil
        }
        
        // fetch object
        //
        // The callback function will be called if the key does not exist.
        // Will cache the callback results into the key and set timeout.
        // Default do nothing.
        //
        // The marshal function will be called before cache.
        //
        // Default callback do nothing, use json.Marshal and json.Unmarshal
        err := f.Fetch(context.Background(), "fetcher_key", &obj,
                redis.WithFetchCallback(callback, 1000*time.Millisecond),
                redis.WithFetchUnmarshal(json.Unmarshal),
                redis.WithFetchMarshal(json.Marshal))
        
        if err != nil {
                fmt.Printf("fetch fail. error: %v\n", err)
                return
        }
}
```

### Config

```yaml
client:
  redis:
    max_idle: 20
    max_active: 100
    max_conn_lifetime: 1000
    idle_timeout: 180000
    timeout: 1000
    wait: true
  service:
    - name: redis_1
      dsn: redis://username:password@127.0.0.1:6379/1?timeout=1000ms
    - name: redis_2
      dsn: redis://username:password@127.0.0.1:6379/2?timeout=1000ms
      max_idle: 22
      max_active: 111
      max_conn_lifetime: 2000
      idle_timeout: 200000
      timeout: 2000

```

## How To Mock

### Client Proxy

```go
package tests

import (
    "testing"
    
    "github.com/agiledragon/gomonkey"
    "github.com/golang/mock/gomock"

    "github.com/wwwangxc/gopkg/redis"
    "github.com/wwwangxc/gopkg/redis/mockredis"
)

func TestMockClientProxy(t *testing.T){
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Mock redis client
    mockConn := mockredis.NewMockConn(ctrl)
    mockConn.EXPECT().Close().Return(nil).AnyTimes()
    mockConn.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
    mockConn.EXPECT().Flush().Return(nil).AnyTimes()
    mockConn.EXPECT().Receive().Return(nil, nil).AnyTimes()

    // Mock locker
    mockLocker := mockredis.NewMockLocker(ctrl)
    mockLocker.EXPECT().TryLock(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
    mockLocker.EXPECT().Lock(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
    mockLocker.EXPECT().Unlock(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

    // Mock fetcher
    mockFetcher := mockredis.NewMockFetcher(ctrl)
    mockFetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

    // Mock client proxy
    mockCli := mockredis.NewMockClientProxy(ctrl)
    mockCli.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any()).Return("reply", nil).AnyTimes()   // Do
    mockCli.EXPECT().Conn().Return(mockConn).AnyTimes()        // Conn
    mockCli.EXPECT().Locker().Return(mockLocker).AnyTimes()    // Locker
    mockCli.EXPECT().Fetcher().Return(mockFetcher).AnyTimes()  // Fetcher
    
    patches := gomonkey.ApplyFunc(redis.NewClientProxy,
        func(string, ...redis.ClientOption) redis.ClientProxy {
            return mockCli
        })
    defer patches.Reset()

    // do something...
}
```

### Locker Proxy

```go
package tests

import (
    "testing"
    
    "github.com/agiledragon/gomonkey"
    "github.com/golang/mock/gomock"

    "github.com/wwwangxc/gopkg/redis"
    "github.com/wwwangxc/gopkg/redis/mockredis"
)

func TestMockLockerProxy(t *testing.T){
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Mock locker
    mockLocker := mockredis.NewMockLocker(ctrl)
    mockLocker.EXPECT().TryLock(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
    mockLocker.EXPECT().Lock(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
    mockLocker.EXPECT().Unlock(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
    mockLocker.EXPECT().LockAndCall(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
    
    patches := gomonkey.ApplyFunc(redis.NewLockerProxy,
        func(string, ...redis.ClientOption) redis.Locker {
            return mockLocker
        })
    defer patches.Reset()

    // do something...
}
```

### Fetcher Proxy

```go
package tests

import (
    "testing"
    
    "github.com/agiledragon/gomonkey"
    "github.com/golang/mock/gomock"

    "github.com/wwwangxc/gopkg/redis"
    "github.com/wwwangxc/gopkg/redis/mockredis"
)

func TestMockFetcherProxy(t *testing.T){
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Mock fetcher
    mockFetcher := mockredis.NewMockFetcher(ctrl)
    mockFetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
    
    patches := gomonkey.ApplyFunc(redis.NewFetcherProxy,
        func(string, ...redis.ClientOption) redis.Fetcher {
            return mockFetcher
        })
    defer patches.Reset()

    // do something...
}
```
