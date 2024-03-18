# gopkg/singleflight

[![Go Report Card](https://goreportcard.com/badge/github.com/wwwangxc/gopkg/singleflight)](https://goreportcard.com/report/github.com/wwwangxc/gopkg/singleflight)
[![GoDoc](https://pkg.go.dev/badge/github.com/wwwangxc/gopkg/singleflight?status.svg)](https://pkg.go.dev/github.com/wwwangxc/gopkg/singleflight)
[![OSCS Status](https://www.oscs1024.com/platform/badge/wwwangxc/gopkg.svg?size=small)](https://www.murphysec.com/dr/c1TuOdJ62DzT0agLwg)

`gopkg/singleflight` is an helper for [golang.org/x/sync/singleflight](https://pkg.go.dev/golang.org/x/sync/singleflight).

## Install

```sh
go get github.com/wwwangxc/gopkg/singleflight
```

## Quick Start
```go
package main

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// the function for key `function_key` will expire in 1 minute. default 1 second
	ret, err := singleflight.Do(ctx, "function_key",
		func(ctx context.Context) (interface{}, error) {
			// dosomething...
			time.Sleep(2 * time.Second)
			return "Successfully", nil
		}, singleflight.WithExpiresIn(time.Minute))

	if err != nil {
		if singleflight.IsTimeout(err) {
			fmt.Println("timeout")
			return
		}

		fmt.Println(err)
		return
	}

	// Output: Successfully
	fmt.Println(ret.(string))
}
```
