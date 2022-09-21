package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// ClientProxy etcd client proxy
//go:generate mockgen -source=client.go -destination=mocketcd/client_mock.go -package=mocketcd
type ClientProxy interface {

	// Put puts a key-value pair into etcd.
	// Note that key,value can be plain bytes array and string is
	// an immutable representation of that bytes array.
	// To get a string of bytes, do string([]byte{0x10, 0x20}).
	Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error)

	// Get retrieves keys.
	// By default, Get will return the value for "key", if any.
	// When passed WithRange(end), Get will return the keys in the range [key, end).
	// When passed WithFromKey(), Get returns keys greater than or equal to key.
	// When passed WithRev(rev) with rev > 0, Get retrieves keys at the given revision;
	// if the required revision is compacted, the request will fail with ErrCompacted .
	// When passed WithLimit(limit), the number of returned keys is bounded by limit.
	// When passed WithSort(), the keys will be sorted.
	Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)

	// Delete deletes a key, or optionally using WithRange(end), [key, end).
	Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error)

	// Watch watches on a key or prefix. The watched events will be returned
	// through the returned channel. If revisions waiting to be sent over the
	// watch are compacted, then the watch will be canceled by the server, the
	// client will post a compacted error watch response, and the channel will close.
	// If the requested revision is 0 or unspecified, the returned channel will
	// return watch events that happen after the server receives the watch request.
	// If the context "ctx" is canceled or timed out, returned "WatchChan" is closed,
	// and "WatchResponse" from this closed channel has zero events and nil "Err()".
	// The context "ctx" MUST be canceled, as soon as watcher is no longer being used,
	// to release the associated resources.
	//
	// If the context is "context.Background/TODO", returned "WatchChan" will
	// not be closed and block until event is triggered, except when server
	// returns a non-recoverable error (e.g. ErrCompacted).
	// For example, when context passed with "WithRequireLeader" and the
	// connected server has no leader (e.g. due to network partition),
	// error "etcdserver: no leader" (ErrNoLeader) will be returned,
	// and then "WatchChan" is closed with non-nil "Err()".
	// In order to prevent a watch stream being stuck in a partitioned node,
	// make sure to wrap context with "WithRequireLeader".
	//
	// Otherwise, as long as the context has not been canceled or timed out,
	// watch will retry on other recoverable errors forever until reconnected.
	Watch(ctx context.Context, key string, opts ...clientv3.OpOption) (clientv3.WatchChan, error)

	// Txn creates a transaction.
	//
	// Step 1:
	// 	If takes a list of comparison. If all comparisons passed in succeed,
	// 	the operations passed into Then() will be executed. Or the operations
	// 	passed into Else() will be executed.
	//
	// Step 2:
	// 	Then takes a list of operations. The Ops list will be executed, if the
	// 	comparisons passed in If() succeed.
	//
	// 	Else takes a list of operations. The Ops list will be executed, if the
	// 	comparisons passed in If() fail.
	//
	// Step 3:
	// 	Commit tries to commit the transaction.
	Txn(ctx context.Context,
		cmps []clientv3.Cmp, thenOps []clientv3.Op, elseOps []clientv3.Op) (*clientv3.TxnResponse, error)

	// Lease get etcd lease proxy
	Lease() LeaseProxy

	// Locker gets the lock operation proxy for the key prefix.
	//
	// It while create an leased session and keep the lease alive until client error
	// or invork close function.
	Locker(prefix string, ttl int) (LockerProxy, error)
}

// NewClientProxy new etcd client proxy
func NewClientProxy(name string, opts ...ClientOption) ClientProxy {
	return &clientProxyImpl{
		name: name,
		opts: opts,
	}
}

type clientProxyImpl struct {
	name string
	opts []ClientOption
}

// Put puts a key-value pair into etcd.
// Note that key,value can be plain bytes array and string is
// an immutable representation of that bytes array.
// To get a string of bytes, do string([]byte{0x10, 0x20}).
func (c *clientProxyImpl) Put(ctx context.Context,
	key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {

	cli, err := c.getClient()
	if err != nil {
		return nil, err
	}

	return cli.Put(ctx, key, val, opts...)
}

// Get retrieves keys.
// By default, Get will return the value for "key", if any.
// When passed WithRange(end), Get will return the keys in the range [key, end).
// When passed WithFromKey(), Get returns keys greater than or equal to key.
// When passed WithRev(rev) with rev > 0, Get retrieves keys at the given revision;
// if the required revision is compacted, the request will fail with ErrCompacted .
// When passed WithLimit(limit), the number of returned keys is bounded by limit.
// When passed WithSort(), the keys will be sorted.
func (c *clientProxyImpl) Get(ctx context.Context,
	key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {

	cli, err := c.getClient()
	if err != nil {
		return nil, err
	}

	return cli.Get(ctx, key, opts...)
}

// Delete deletes a key, or optionally using WithRange(end), [key, end).
func (c *clientProxyImpl) Delete(ctx context.Context,
	key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {

	cli, err := c.getClient()
	if err != nil {
		return nil, err
	}

	return cli.Delete(ctx, key, opts...)
}

// Watch watches on a key or prefix. The watched events will be returned
// through the returned channel. If revisions waiting to be sent over the
// watch are compacted, then the watch will be canceled by the server, the
// client will post a compacted error watch response, and the channel will close.
// If the requested revision is 0 or unspecified, the returned channel will
// return watch events that happen after the server receives the watch request.
// If the context "ctx" is canceled or timed out, returned "WatchChan" is closed,
// and "WatchResponse" from this closed channel has zero events and nil "Err()".
// The context "ctx" MUST be canceled, as soon as watcher is no longer being used,
// to release the associated resources.
//
// If the context is "context.Background/TODO", returned "WatchChan" will
// not be closed and block until event is triggered, except when server
// returns a non-recoverable error (e.g. ErrCompacted).
// For example, when context passed with "WithRequireLeader" and the
// connected server has no leader (e.g. due to network partition),
// error "etcdserver: no leader" (ErrNoLeader) will be returned,
// and then "WatchChan" is closed with non-nil "Err()".
// In order to prevent a watch stream being stuck in a partitioned node,
// make sure to wrap context with "WithRequireLeader".
//
// Otherwise, as long as the context has not been canceled or timed out,
// watch will retry on other recoverable errors forever until reconnected.
func (c *clientProxyImpl) Watch(ctx context.Context,
	key string, opts ...clientv3.OpOption) (clientv3.WatchChan, error) {

	cli, err := c.getClient()
	if err != nil {
		return nil, err
	}

	return cli.Watch(ctx, key, opts...), nil
}

// Txn creates a transaction.
//
// Step 1:
// 	If takes a list of comparison. If all comparisons passed in succeed,
// 	the operations passed into Then() will be executed. Or the operations
// 	passed into Else() will be executed.
//
// Step 2:
// 	Then takes a list of operations. The Ops list will be executed, if the
// 	comparisons passed in If() succeed.
//
// 	Else takes a list of operations. The Ops list will be executed, if the
// 	comparisons passed in If() fail.
//
// Step 3:
// 	Commit tries to commit the transaction.
func (c *clientProxyImpl) Txn(ctx context.Context,
	cmps []clientv3.Cmp, thenOps []clientv3.Op, elseOps []clientv3.Op) (*clientv3.TxnResponse, error) {

	cli, err := c.getClient()
	if err != nil {
		return nil, err
	}

	return cli.Txn(ctx).If(cmps...).Then(thenOps...).Else(elseOps...).Commit()
}

// Lease get etcd lease proxy
func (c *clientProxyImpl) Lease() LeaseProxy {
	return newLeaseProxy(c.name, c.opts...)
}

// Locker gets the lock operation proxy for the key prefix.
//
// It while create an leased session and keep the lease alive until client error
// or invork close function.
func (c *clientProxyImpl) Locker(prefix string, ttl int) (LockerProxy, error) {
	cli, err := c.getClient()
	if err != nil {
		return nil, err
	}

	return newLockerProxy(cli, prefix, ttl)
}

func (c *clientProxyImpl) getClient() (*clientv3.Client, error) {
	return getETCDClient(c.name, c.opts...)
}
