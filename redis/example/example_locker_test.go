package example

import (
	"context"
	"fmt"
	"time"

	// gopkg/redis will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/redis"
)

func ExampleNewLockerProxy() {
	l := redis.NewLockerProxy("client_name",
		redis.WithClientDSN("dsn"),             // set dsn, default use database.client.dsn
		redis.WithClientMaxIdle(20),            // set max idel. default 2048
		redis.WithClientMaxActive(100),         // set max active. default 0
		redis.WithClientIdleTimeout(180000),    // set idle timeout. unit millisecond, default 180000
		redis.WithClientTimeout(1000),          // set command timeout. unit millisecond, default 1000
		redis.WithClientMaxConnLifetime(10000), // set max conn life time, default 0
		redis.WithClientWait(true),             // set wait
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
