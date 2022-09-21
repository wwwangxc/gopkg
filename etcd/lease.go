package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// LeaseProxy etcd lease proxy
//go:generate mockgen -source=lease.go -destination=mocketcd/lease_mock.go -package=mocketcd
type LeaseProxy interface {
	// Grant creates a new lease.
	Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error)

	// Revoke revokes the given lease.
	Revoke(ctx context.Context, id clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error)

	// TimeToLive retrieves the lease information of the given lease ID.
	TimeToLive(ctx context.Context,
		id clientv3.LeaseID, opts ...clientv3.LeaseOption) (*clientv3.LeaseTimeToLiveResponse, error)

	// Leases retrieves all leases.
	Leases(ctx context.Context) (*clientv3.LeaseLeasesResponse, error)

	// KeepAlive attempts to keep the given lease alive forever. If the keepalive responses posted
	// to the channel are not consumed promptly the channel may become full. When full, the lease
	// client will continue sending keep alive requests to the etcd server, but will drop responses
	// until there is capacity on the channel to send more responses.
	//
	// If client keep alive loop halts with an unexpected error (e.g. "etcdserver: no leader") or
	// canceled by the caller (e.g. context.Canceled), KeepAlive returns a ErrKeepAliveHalted error
	// containing the error reason.
	//
	// The returned "LeaseKeepAliveResponse" channel closes if underlying keep
	// alive stream is interrupted in some way the client cannot handle itself;
	// given context "ctx" is canceled or timed out.
	KeepAlive(ctx context.Context, id clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error)
}

func newLeaseProxy(name string, opts ...ClientOption) LeaseProxy {
	return &leaseProxyImpl{
		name: name,
		opts: opts,
	}
}

type leaseProxyImpl struct {
	name string
	opts []ClientOption
}

// Grant creates a new lease.
func (l *leaseProxyImpl) Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	cli, err := l.GetCli()
	if err != nil {
		return nil, err
	}

	return cli.Grant(ctx, ttl)
}

// Revoke revokes the given lease.
func (l *leaseProxyImpl) Revoke(ctx context.Context,
	id clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error) {

	cli, err := l.GetCli()
	if err != nil {
		return nil, err
	}

	return cli.Revoke(ctx, id)
}

// TimeToLive retrieves the lease information of the given lease ID.
func (l *leaseProxyImpl) TimeToLive(ctx context.Context,
	id clientv3.LeaseID, opts ...clientv3.LeaseOption) (*clientv3.LeaseTimeToLiveResponse, error) {

	cli, err := l.GetCli()
	if err != nil {
		return nil, err
	}

	return cli.TimeToLive(ctx, id, opts...)
}

// Leases retrieves all leases.
func (l *leaseProxyImpl) Leases(ctx context.Context) (*clientv3.LeaseLeasesResponse, error) {
	cli, err := l.GetCli()
	if err != nil {
		return nil, err
	}

	return cli.Leases(ctx)
}

// KeepAlive attempts to keep the given lease alive forever. If the keepalive responses posted
// to the channel are not consumed promptly the channel may become full. When full, the lease
// client will continue sending keep alive requests to the etcd server, but will drop responses
// until there is capacity on the channel to send more responses.
//
// If client keep alive loop halts with an unexpected error (e.g. "etcdserver: no leader") or
// canceled by the caller (e.g. context.Canceled), KeepAlive returns a ErrKeepAliveHalted error
// containing the error reason.
//
// The returned "LeaseKeepAliveResponse" channel closes if underlying keep
// alive stream is interrupted in some way the client cannot handle itself;
// given context "ctx" is canceled or timed out.
func (l *leaseProxyImpl) KeepAlive(ctx context.Context,
	id clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {

	cli, err := l.GetCli()
	if err != nil {
		return nil, err
	}

	return cli.KeepAlive(ctx, id)
}

func (l *leaseProxyImpl) GetCli() (*clientv3.Client, error) {
	return getETCDClient(l.name, l.opts...)
}
