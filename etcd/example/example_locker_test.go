package example_test

import (
	// gopkg/etcd will automatically read configuration
	// files (./app.yaml) when package loaded
	"context"
	"fmt"

	"github.com/wwwangxc/gopkg/etcd"
)

func ExampleLockerProxy_Lock() {
	lockKeyPrefix := "lock/example/lock"

	// gets the lock operation proxy for the key prefix.
	// It while create an leased session and keep the lease alive until client error
	// or invork close function.
	locker, err := etcd.NewClientProxy("etcd1").Locker(lockKeyPrefix, 3)
	if err != nil {
		fmt.Printf("get locker proxy fail:%v\n", err)
		return
	}

	defer func() {

		// Close orphans the session and revokes the session lease.
		if err := locker.Close(); err != nil {
			fmt.Printf("locker close fail:%v", err)
			return
		}
	}()

	// Will block the current goroutine until locked.
	// If the context is canceled while trying to acquire the lock, the mutex tries to clean its stale lock entry.
	if err := locker.Lock(context.Background()); err != nil {
		fmt.Printf("lock fail:%v\n", err)
		return
	}

	// lock success

	defer func() {
		if err := locker.Unlock(context.Background()); err != nil {
			fmt.Printf("unlock fail:%v", err)
		}
	}()

	// do something...
}

func ExampleLockerProxy_TryLock() {
	lockKeyPrefix := "lock/example/try_lock"

	// gets the lock operation proxy for the key prefix.
	// It while create an leased session and keep the lease alive until client error
	// or invork close function.
	locker, err := etcd.NewClientProxy("etcd1").Locker(lockKeyPrefix, 3)
	if err != nil {
		fmt.Printf("get locker proxy fail:%v\n", err)
		return
	}

	defer func() {

		// Close orphans the session and revokes the session lease.
		if err := locker.Close(); err != nil {
			fmt.Printf("locker close fail:%v", err)
			return
		}
	}()

	if err = locker.TryLock(context.Background()); err != nil {

		// return 'ErrLockNotAcquired' when lock not acquired.
		if etcd.IsErrLockNotAcquired(err) {
			fmt.Printf("lock not acquired\n")
			return
		}

		fmt.Printf("try lock fail:%v\n", err)
		return
	}

	// lock success

	defer func() {
		if err := locker.Unlock(context.Background()); err != nil {
			fmt.Printf("unlock fail:%v", err)
		}
	}()

	// do something...
}

func ExampleLockerProxy_LockAndCall() {
	lockKeyPrefix := "lock/example/try_lock"

	// gets the lock operation proxy for the key prefix.
	// It while create an leased session and keep the lease alive until client error
	// or invork close function.
	locker, err := etcd.NewClientProxy("etcd1").Locker(lockKeyPrefix, 3)
	if err != nil {
		fmt.Printf("get locker proxy fail:%v\n", err)
		return
	}

	defer func() {

		// Close orphans the session and revokes the session lease.
		if err := locker.Close(); err != nil {
			fmt.Printf("locker close fail:%v", err)
			return
		}
	}()

	f := func() error {
		// do something...
		return nil
	}

	if err := locker.LockAndCall(context.Background(), f); err != nil {
		fmt.Printf("lock and call fail. error: %v\n", err)
		return
	}
}
