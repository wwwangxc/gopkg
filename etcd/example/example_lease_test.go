package example_test

import (
	"context"
	"fmt"
	"time"

	// gopkg/etcd will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/etcd"
)

func ExampleLeaseProxy_Grant() {
	lease := etcd.NewClientProxy("etcd1").Lease()

	// create a lease
	id, err := etcd.LeaseGrantResult(lease.Grant(context.Background(), 10))
	if err != nil {
		fmt.Printf("lease grant fail. error:%v", err)
		return
	}

	fmt.Printf("lease:0x%x\n", id)
}

func ExampleLeaseProxy_RevokeResult() {
	lease := etcd.NewClientProxy("etcd1").Lease()

	// create a lease
	id, err := etcd.LeaseGrantResult(lease.Grant(context.Background(), 10))
	if err != nil {
		fmt.Printf("lease grant fail. error:%v", err)
		return
	}

	// revoke a lease
	err = etcd.LeaseRevokeResult(lease.Revoke(context.Background(), id))
	if err != nil {
		fmt.Printf("lease revoke fail. error:%v", err)
		return
	}
}

func ExampleLeaseProxy_LeaseTimeToLiveResult() {
	lease := etcd.NewClientProxy("etcd1").Lease()

	// create a lease
	id, err := etcd.LeaseGrantResult(lease.Grant(context.Background(), 10))
	if err != nil {
		fmt.Printf("lease grant fail. error:%v", err)
		return
	}

	for {
		// get lease ttl
		ttl, err := etcd.LeaseTimeToLiveResult(lease.TimeToLive(context.Background(), id))
		if err != nil {
			fmt.Printf("get lease ttl fail. error:%v", err)
			return
		}

		if ttl == -1 {
			fmt.Printf("lease:0x%x expired\n", id)
			break
		}

		time.Sleep(time.Second)
	}
}

func ExampleLeaseProxy_KeepAlive() {
	lease := etcd.NewClientProxy("etcd1").Lease()

	// create a lease
	id, err := etcd.LeaseGrantResult(lease.Grant(context.Background(), 10))
	if err != nil {
		fmt.Printf("lease grant fail. error:%v", err)
		return
	}

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
