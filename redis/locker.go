package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

// LockerProxy distributed lock provider
//go:generate mockgen -source=locker.go -destination=mockredis/locker_mock.go -package=mockredis
type LockerProxy interface {

	// LockAndCall try get lock first and call f() when lock acquired. Unlock will be performed
	// regardless of whether the f reports an error or not.
	//
	// Will block the current goroutine when lock not acquired
	// Will reentrant lock when UUID option not empty.
	// If Heartbeat option not empty and not a reentrant lock, will automatically
	// renewal until unlocked.
	LockAndCall(ctx context.Context, key string, f func() error, opts ...LockOption) error

	// TryLock try get lock, if lock acquired will return lock uuid.
	//
	// Not block the current goroutine.
	// Return ErrLockNotAcquired when lock not acquired.
	// Will reentrant lock when UUID option not empty.
	// If Heartbeat option not empty and not a reentrant lock, will automatically
	// renewal until unlocked.
	TryLock(ctx context.Context, key string, opts ...LockOption) (uuid string, err error)

	// Lock try get lock until the context canceled or the lock acquired
	//
	// Will block the current goroutine.
	// Will reentrant lock when UUID option not empty.
	// If Heartbeat option not empty and not a reentrant lock, will automatically
	// renewal until unlocked.
	Lock(ctx context.Context, key string, opts ...LockOption) (uuid string, err error)

	// Unlock
	//
	// Return ErrLockNotExist if the key does not exist.
	// Return ErrNotOwnerOfKey if the uuid invalid.
	// Support reentrant unlock.
	Unlock(ctx context.Context, key, uuid string) error
}

type lockerImpl struct {
	name string
	opts []ClientOption
}

// NewLockerProxy new locker proxy
func NewLockerProxy(name string, opts ...ClientOption) LockerProxy {
	return &lockerImpl{
		name: name,
		opts: opts,
	}
}

// LockAndCall try get lock first and call f() when lock acquired. Unlock will be performed
// regardless of whether the f reports an error or not.
//
// Will block the current goroutine when lock not acquired
// Will reentrant lock when UUID option not empty.
// If Heartbeat option not empty and not a reentrant lock, will automatically
// renewal until unlocked.
func (l *lockerImpl) LockAndCall(ctx context.Context, key string, f func() error, opts ...LockOption) error {
	var err error

	if f == nil {
		return nil
	}

	uuid, err := l.Lock(ctx, key, opts...)
	if err != nil {
		err = fmt.Errorf("lock fail case %v", err)
		return err
	}

	defer func() {
		if err := l.Unlock(ctx, key, uuid); err != nil {
			logErrorf("lock:%s unlock fail: %v", key, err)
		}
	}()

	return f()
}

// TryLock try get lock, if lock acquired will return lock uuid.
//
// Not block the current goroutine.
// Return ErrLockNotAcquired when lock not acquired.
// Will reentrant lock when UUID option not empty.
// If Heartbeat option not empty and not a reentrant lock, will automatically
// renewal until unlocked.
func (l *lockerImpl) TryLock(ctx context.Context, key string, opts ...LockOption) (string, error) {
	k := fmt.Sprintf("%s.lock", strings.TrimSuffix(key, ".lock"))
	options := newLockOptions(opts...)
	script := redigo.NewScript(1, luaScriptLock)

	conn := l.getConn()
	defer func() {
		if err := conn.Close(); err != nil {
			logErrorf("connect close fail. error:%v", err)
		}
	}()

	lockCount, err := Int(script.DoContext(ctx, conn, k, options.UUID, options.Expire.Milliseconds()))
	if err != nil {
		return "", err
	}

	if lockCount == 0 {
		return "", ErrLockNotAcquired
	}

	if lockCount == 1 && options.Heartbeat > 0 {
		go l.sendLockHeartbeat(key, options.Expire, options.Heartbeat)
	}

	return options.UUID, nil
}

// Lock try get lock until the context canceled or the lock acquired
//
// Will block the current goroutine.
// Will reentrant lock when UUID option not empty.
// If Heartbeat option not empty and not a reentrant lock, will automatically
// renewal until unlocked.
func (l *lockerImpl) Lock(ctx context.Context, key string, opts ...LockOption) (string, error) {
	options := newLockOptions(opts...)
	for {
		select {
		case <-ctx.Done():
			return "", ErrTimeout

		default:
			uuid, err := l.TryLock(ctx, key, opts...)
			if err != nil {
				if IsLockNotAcquired(err) {
					time.Sleep(options.Retry)
					continue
				}

				return "", err
			}

			return uuid, nil
		}
	}
}

// Unlock
//
// Return ErrLockNotExist if the key does not exist.
// Return ErrNotOwnerOfKey if the uuid invalid.
// Support reentrant unlock.
func (l *lockerImpl) Unlock(ctx context.Context, key, uuid string) error {
	k := fmt.Sprintf("%s.lock", strings.TrimSuffix(key, ".lock"))
	script := redigo.NewScript(1, luaScriptUnlock)

	conn := l.getConn()
	defer func() {
		if err := conn.Close(); err != nil {
			logErrorf("connect close fail. error:%v", err)
		}
	}()

	ret, err := Int(script.DoContext(ctx, conn, k, uuid))
	if err != nil {
		return err
	}

	switch ret {
	case 0:
		return ErrLockNotExist
	case 1:
		return ErrNotOwnerOfLock
	case 2:
		return errors.New("locker key delete fail")
	case 666:
		return nil
	}

	return errors.New("error unknown")
}

func (l *lockerImpl) sendLockHeartbeat(key string, expire, heartbeatInterval time.Duration) {
	conn := l.getConn()
	defer func() {
		if err := conn.Close(); err != nil {
			logErrorf("connect close fail. error:%v", err)
		}
	}()

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		exist, err := Bool(conn.Do("EXISTS", key))
		if err != nil {
			logErrorf(err.Error())
			return
		}

		if !exist {
			return
		}

		_, err = Bool(conn.Do("PEXPIRE", key, expire.Milliseconds()))
		if err != nil {
			logErrorf(err.Error())
			return
		}
	}
}

func (l *lockerImpl) getConn() redigo.Conn {
	return getRedisPool(l.name, l.opts...).Get()
}
