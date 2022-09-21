package etcd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/wwwangxc/gopkg/etcd/log"
)

// LockerProxy etcd locker proxy
//go:generate mockgen -source=locker.go -destination=mocketcd/locker_mock.go -package=mocketcd
type LockerProxy interface {

	// LockAndCall try get lock first and call f() when lock acquired. Unlock will be performed
	// regardless of whether the f reports an error or not.
	//
	// Will block the current goroutine when lock not acquired
	// Will reentrant lock when UUID option not empty.
	// If Heartbeat option not empty and not a reentrant lock, will automatically
	// renewal until unlocked.
	LockAndCall(ctx context.Context, f func() error) error

	// Lock locks the mutex with a cancelable context. If the context is canceled
	// while trying to acquire the lock, the mutex tries to clean its stale lock entry.
	Lock(ctx context.Context) error

	// TryLock locks the mutex if not already locked by another session.
	// If lock is held by another session, return immediately after attempting necessary cleanup
	// The ctx argument is used for the sending/receiving Txn RPC.
	// Return 'ErrLockNotAcquired' when lock not acquired.
	TryLock(ctx context.Context) error

	// Unlock
	// Return error if the lock key delete fail.
	Unlock(ctx context.Context) error

	// Close orphans the session and revokes the session lease.
	Close() error
}

func newLockerProxy(cli *clientv3.Client, prefix string, ttl int) (LockerProxy, error) {
	s, err := concurrency.NewSession(cli, concurrency.WithTTL(ttl))
	if err != nil {
		return nil, err
	}

	return &lockerProxyImpl{
		session: s,
		mutex:   concurrency.NewMutex(s, strings.TrimSuffix(prefix, "/")),
	}, nil
}

type lockerProxyImpl struct {
	session *concurrency.Session
	mutex   *concurrency.Mutex
}

// LockAndCall try get lock first and call f() when lock acquired. Unlock will be performed
// regardless of whether the f reports an error or not.
//
// Will block the current goroutine when lock not acquired
// Will reentrant lock when UUID option not empty.
// If Heartbeat option not empty and not a reentrant lock, will automatically
// renewal until unlocked.
func (l *lockerProxyImpl) LockAndCall(ctx context.Context, f func() error) error {
	if f == nil {
		return nil
	}

	if err := l.Lock(ctx); err != nil {
		return fmt.Errorf("lock fail case %v", err)
	}

	defer func() {
		if err := l.Unlock(ctx); err != nil {
			log.Errorf("unlock fail case %v", err)
		}
	}()

	return f()
}

// Lock locks the mutex with a cancelable context. If the context is canceled
// while trying to acquire the lock, the mutex tries to clean its stale lock entry.
func (l *lockerProxyImpl) Lock(ctx context.Context) error {
	return l.mutex.Lock(ctx)
}

// TryLock locks the mutex if not already locked by another session.
// If lock is held by another session, return immediately after attempting necessary cleanup
// The ctx argument is used for the sending/receiving Txn RPC.
// Return 'ErrLockNotAcquired' when lock not acquired.
func (l *lockerProxyImpl) TryLock(ctx context.Context) error {
	if err := l.mutex.TryLock(ctx); err != nil {
		if errors.Is(err, concurrency.ErrLocked) {
			return ErrLockNotAcquired
		}

		return err
	}

	return nil
}

// Unlock
// Return error if the lock key delete fail.
func (l *lockerProxyImpl) Unlock(ctx context.Context) error {
	return l.mutex.Unlock(ctx)
}

// Close orphans the session and revokes the session lease.
func (l *lockerProxyImpl) Close() error {
	return l.session.Close()
}
