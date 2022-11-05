package concurrency_test

import (
	"context"
	"fmt"

	"github.com/wwwangxc/gopkg/concurrency"
)

func Example() {
	handlers := []concurrency.Handler{
		&handlerImpl{arg1: "123", arg2: "456"},
		&handlerImpl{arg1: "123", arg2: "456"},
		&handlerImpl{arg1: "123", arg2: "456"},
		&handlerImpl{arg1: "123", arg2: "456"},
		&handlerImpl{arg1: "123", arg2: "456"},
	}

	result := concurrency.Start(context.Background(), handlers, 2)

	// return true when no handler invoke failed
	result.Succeed()

	// return true when there are handler invoke failed
	result.Failed()

	// return merged error
	//
	// Format like:
	// 2 errors occurred:
	//     * error message ...
	//     * [PANIC]panic message ...
	result.MergedError()

	// return collection of all errors
	result.Errors()

	// return collection of all results
	result.Results()
}

type handlerImpl struct {
	arg1 string
	arg2 string
}

func (e *handlerImpl) Invoke(ctx context.Context) (interface{}, error) {
	fmt.Println(e.arg1)
	fmt.Println(e.arg2)
	return nil, nil
}
