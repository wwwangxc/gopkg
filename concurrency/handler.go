package concurrency

import "context"

// Handler concurrenty handler
type Handler interface {
	Invoke(ctx context.Context) (interface{}, error)
}
