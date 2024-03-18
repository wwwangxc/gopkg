package singleflight

import "errors"

var (
	// ErrTimeout timeout error
	ErrTimeout = errors.New("timeout")
)

// IsTimeout will return true when the error is timeout
func IsTimeout(err error) bool {
	return errors.Is(err, ErrTimeout)
}
