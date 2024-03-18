package singleflight_test

import (
	"context"
	"fmt"
	"time"

	"github.com/wwwangxc/gopkg/singleflight"
)

func Example() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ret, err := singleflight.Do(ctx, "function_key",
		func(ctx context.Context) (interface{}, error) {
			// dosomething...
			time.Sleep(2 * time.Second)
			return "Successfully", nil
		}) // the function for key `function_key` will expire in 1 second

	switch {
	case err != nil:
		fmt.Println(err.Error())
	case ret != nil:
		fmt.Println(ret.(string))
	default:
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ret, err = singleflight.Do(ctx, "function_key",
		func(ctx context.Context) (interface{}, error) {
			// dosomething...
			return "Successfully", nil
		}, singleflight.WithExpiresIn(time.Minute)) // the function for key `function_key` will expire in 1 minute

	switch {
	case err != nil:
		fmt.Println(err.Error())
	case ret != nil:
		fmt.Println(ret.(string))
	default:
	}

	// Output:
	// timeout
	// Successfully
}
