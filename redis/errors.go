package redis

import "errors"

var (
	// ErrTimeout
	ErrTimeout = errors.New("timeout")

	// ErrLockNotAcquired lock not acquired
	ErrLockNotAcquired = errors.New("lock not acquired")

	// ErrLockNotExist lock dose not exist
	ErrLockNotExist = errors.New("lock does not exist")

	// ErrNotOwnerOfLock not the owner of the key
	ErrNotOwnerOfLock = errors.New("not the owner of the lock")

	// ErrKeyNotExist key not exist
	ErrKeyNotExist = errors.New("key not exist")
)

// IsTimeout is timeout error
func IsTimeout(err error) bool {
	return errors.Is(err, ErrTimeout)
}

// IsLockNotAcquired is lock not acquired error
func IsLockNotAcquired(err error) bool {
	return errors.Is(err, ErrLockNotAcquired)
}

// IsLockNotExist is lock not exist error
func IsLockNotExist(err error) bool {
	return errors.Is(err, ErrLockNotExist)
}

// IsErrNotOwnerOfLock is not owner of lock
func IsErrNotOwnerOfLock(err error) bool {
	return errors.Is(err, ErrNotOwnerOfLock)
}

// IsKeyNotExist is key not exist error
func IsKeyNotExist(err error) bool {
	return errors.Is(err, ErrKeyNotExist)
}
