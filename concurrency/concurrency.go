package concurrency

import (
	"context"
	"fmt"
	"sync"
)

// Start run concurrenty
//
// Result is collection of all handlers results and error
func Start(ctx context.Context, handlers []Handler, concurrenty uint8) *Result {
	result := &Result{}
	limitCh := make(chan struct{}, concurrenty)
	var wg sync.WaitGroup
	wg.Add(len(handlers))

	go func() {
		for _, v := range handlers {
			limitCh <- struct{}{}
			go func(handler Handler) {
				defer func() {
					<-limitCh
					if e := recover(); e != nil {
						result.append(nil, fmt.Errorf("[PANIC]%v", e))
					}
					wg.Done()
				}()

				result.append(handler.Invoke(ctx))
			}(v)
		}
	}()

	wg.Wait()
	return result
}
