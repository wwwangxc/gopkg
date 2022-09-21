package etcd

import "errors"

var (
	// ErrLockNotAcquired lock not acquired
	ErrLockNotAcquired = errors.New("lock not acquired")
)

// IsErrLockNotAcquired is lock not acquired error
func IsErrLockNotAcquired(err error) bool {
	return errors.Is(err, ErrLockNotAcquired)
}
