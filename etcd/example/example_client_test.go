package example_test

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"

	// gopkg/etcd will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/etcd"
)

func ExampleNewClientProxy() {
	_ = etcd.NewClientProxy("etcd1",
		etcd.WithEndpoints([]string{"127.0.0.1:2379", "127.0.0.1:2380"}), // set endpoints
		etcd.WithTimeout(3000),                                 // set timeout, unit millisecond, default 1000.
		etcd.WithAuth("username", "password"),                  // set username and password for authentication
		etcd.WithTLSKeyPath("/usr/local/etcd_conf/key.pem"),    // set tls key file path.
		etcd.WithTLSCertPath("/usr/local/etcd_conf/cert.pem"),  // set tls cert file path.
		etcd.WithCACertPath("/usr/local/etcd_conf/cacert.pem"), // set ca cert file path.
	)
}

func ExampleClientProxy_Put() {
	// do etcd put operation
	_, err := etcd.NewClientProxy("etcd1").Put(context.Background(), "key", "val")
	if err != nil {
		fmt.Printf("put operation fail. error:%v", err)
		return
	}

	// or

	// do etcd put operation and convert result to an error
	if err = etcd.PutResult(etcd.NewClientProxy("etcd").Put(context.Background(), "key", "val")); err != nil {
		fmt.Printf("put operation fail. error:%v", err)
		return
	}
}

func ExampleClientProxy_PutWithLease() {
	cli := etcd.NewClientProxy("etcd1")
	lease := cli.Lease()

	// create a lease
	id, err := etcd.LeaseGrantResult(lease.Grant(context.Background(), 10))
	if err != nil {
		fmt.Printf("lease grant fail. error:%v", err)
		return
	}

	// put with lease
	err = etcd.PutResult(cli.Put(context.Background(), "key", "val", clientv3.WithLease(id)))
	if err != nil {
		fmt.Printf("put operation fail. error:%v", err)
	}
}

func ExampleClientProxy_Get() {
	// do etcd get operation
	resp, err := etcd.NewClientProxy("etcd1").Get(context.Background(), "key")
	if err != nil {
		fmt.Printf("get operation fail. error:%v", err)
		return
	}

	for k, v := range resp.Kvs {
		fmt.Printf("key index: %d\n", k)
		fmt.Printf("val: %s\n", v)
	}

	// or

	// do etcd get operation and convert result to map[string]string and an error
	m, err := etcd.GetResult(etcd.NewClientProxy("etcd").Get(context.Background(), "key"))
	if err != nil {
		fmt.Printf("get operation fail. error:%v", err)
		return
	}

	for k, v := range m {
		fmt.Printf("key: %s\n", k)
		fmt.Printf("val: %s\n", v)
	}
}

func ExampleClientProxy_Delete() {
	// do etcd delete operation
	resp, err := etcd.NewClientProxy("etcd1").Delete(context.Background(), "key")
	if err != nil {
		fmt.Printf("delete operation fail. error:%v", err)
		return
	}

	fmt.Printf("number of keys deleted: %d\n", resp.Deleted)

	// or

	// do etcd delte operation and convert result to delete number and an error
	deletedNumber, err := etcd.DeleteResult(etcd.NewClientProxy("etcd").Delete(context.Background(), "key"))
	if err != nil {
		fmt.Printf("delete operation fail. error:%v", err)
		return
	}

	fmt.Printf("number of keys deleted: %d\n", deletedNumber)
}

func ExampleClientProxy_Watch() {
	// do etcd watch operation
	watchChan, err := etcd.NewClientProxy("etcd1").Watch(context.Background(), "key")
	if err != nil {
		fmt.Printf("watch operation fail. error:%v", err)
		return
	}

	for _ = range watchChan {
		// do something...
	}
}

func ExampleClientProxy_Txn() {
	_, err := etcd.NewClientProxy("etcd1").Txn(context.Background(),
		[]clientv3.Cmp{clientv3.Compare(clientv3.Value("key"), "=", "val")}, // if key's value == val
		[]clientv3.Op{clientv3.OpPut("key", "val1")},                        // then put key's value = val1
		[]clientv3.Op{clientv3.OpPut("key", "val")})                         // else put key's value = val
	if err != nil {
		fmt.Printf("txn fail. error:%v", err)
		return
	}

	// or

	err = etcd.TxnResult(etcd.NewClientProxy("etcd1").Txn(context.Background(),
		[]clientv3.Cmp{clientv3.Compare(clientv3.Value("key"), "=", "val")}, // if key's value == val
		[]clientv3.Op{clientv3.OpPut("key", "val1")},                        // then put key's value = val1
		[]clientv3.Op{clientv3.OpPut("key", "val")}))                        // else put key's value = val
	if err != nil {
		fmt.Printf("txn fail. error:%v", err)
		return
	}

}
