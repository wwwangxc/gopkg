package etcd_test

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/wwwangxc/gopkg/etcd"
)

func ExampleClientProxy() {
	// new client proxy
	cli := etcd.NewClientProxy("etcd1",
		etcd.WithEndpoints([]string{"127.0.0.1:2379", "127.0.0.1:2380"}), // set endpoints
		etcd.WithTimeout(3000),                                 // set timeout, unit millisecond, default 1000.
		etcd.WithAuth("username", "password"),                  // set username and password for authentication
		etcd.WithTLSKeyPath("/usr/local/etcd_conf/key.pem"),    // set tls key file path.
		etcd.WithTLSCertPath("/usr/local/etcd_conf/cert.pem"),  // set tls cert file path.
		etcd.WithCACertPath("/usr/local/etcd_conf/cacert.pem"), // set ca cert file path.
	)

	// do etcd put operation and convert result to an error
	if err := etcd.PutResult(cli.Put(context.Background(), "key", "val")); err != nil {
		fmt.Printf("put operation fail. error:%v", err)
		return
	}

	// create a lease
	id, err := etcd.LeaseGrantResult(cli.Lease().Grant(context.Background(), 10))
	if err != nil {
		fmt.Printf("lease grant fail. error:%v", err)
		return
	}

	// put with lease
	err = etcd.PutResult(cli.Put(context.Background(), "key", "val", clientv3.WithLease(id)))
	if err != nil {
		fmt.Printf("put operation fail. error:%v", err)
	}

	// do etcd get operation and convert result to map[string]string and an error
	m, err := etcd.GetResult(cli.Get(context.Background(), "key"))
	if err != nil {
		fmt.Printf("get operation fail. error:%v", err)
		return
	}

	for k, v := range m {
		fmt.Printf("key: %s\n", k)
		fmt.Printf("val: %s\n", v)
	}

	// do etcd delte operation and convert result to delete number and an error
	delNum, err := etcd.DeleteResult(cli.Delete(context.Background(), "key"))
	if err != nil {
		fmt.Printf("delete operation fail. error:%v", err)
		return
	}
	fmt.Printf("delete number: %d\n", delNum)

	// transaction
	err = etcd.TxnResult(cli.Txn(context.Background(),
		[]clientv3.Cmp{clientv3.Compare(clientv3.Value("key"), "=", "val")}, // if key's value == val
		[]clientv3.Op{clientv3.OpPut("key", "val1")},                        // then put key's value = val1
		[]clientv3.Op{clientv3.OpPut("key", "val")}))                        // else put key's value = val
	if err != nil {
		fmt.Printf("txn fail. error:%v", err)
		return
	}
}

func ExampleLeaseProxy() {
	// new lease proxy
	lease := etcd.NewClientProxy("etcd1").Lease()

	// create a lease
	id, err := etcd.LeaseGrantResult(lease.Grant(context.Background(), 10))
	if err != nil {
		fmt.Printf("lease grant fail. error:%v", err)
		return
	}
	fmt.Printf("lease:0x%x\n", id)

	// revoke a lease
	err = etcd.LeaseRevokeResult(lease.Revoke(context.Background(), id))
	if err != nil {
		fmt.Printf("lease revoke fail. error:%v", err)
		return
	}

	go func(id clientv3.LeaseID) {
		for {
			// get lease ttl
			ttl, err := etcd.LeaseTimeToLiveResult(lease.TimeToLive(context.Background(), id))
			if err != nil {
				fmt.Printf("get lease ttl fail. error:%v", err)
				return
			}

			if ttl == -1 {
				break
			}

			time.Sleep(time.Second)
		}

		fmt.Printf("lease:0x%x expired\n", id)
	}(id)

	// keep alive
	ch, err := lease.KeepAlive(context.Background(), id)
	if err != nil {
		fmt.Printf("lease keep alive fail. error:%v", err)
		return
	}

	for {
		ka := <-ch
		if ka == nil {
			fmt.Println("lease timeout")
			return
		}
		fmt.Println("ttl:", ka.TTL)
	}
}

func ExampleLockerProxy() {
	lockKeyPrefix := "lock/example/lock"

	// Gets the lock operation proxy for the key prefix.
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

	// Unlock
	if err := locker.Unlock(context.Background()); err != nil {
		fmt.Printf("unlock fail:%v", err)
		return
	}

	// Try lock
	// Locks the mutex if not already locked by another session.
	if err = locker.TryLock(context.Background()); err != nil {

		// return 'ErrLockNotAcquired' when lock not acquired.
		if etcd.IsErrLockNotAcquired(err) {
			fmt.Printf("lock not acquired\n")
			return
		}

		fmt.Printf("try lock fail:%v\n", err)
		return
	}

	// Unlock
	if err := locker.Unlock(context.Background()); err != nil {
		fmt.Printf("unlock fail:%v", err)
		return
	}

	f := func() error {
		// do something...
		return nil
	}

	// Try get lock first and call f() when lock acquired. Unlock will be performed
	// regardless of whether the f reports an error or not.
	if err := locker.LockAndCall(context.Background(), f); err != nil {
		fmt.Printf("lock and call fail. error: %v\n", err)
		return
	}
}
