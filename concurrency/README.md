# gopkg/concurrency

[![Go Report Card](https://goreportcard.com/badge/github.com/wwwangxc/gopkg/concurrency)](https://goreportcard.com/report/github.com/wwwangxc/gopkg/concurrency)
[![GoDoc](https://pkg.go.dev/badge/github.com/wwwangxc/gopkg/concurrency?status.svg)](https://pkg.go.dev/github.com/wwwangxc/gopkg/concurrency)
[![OSCS Status](https://www.oscs1024.com/platform/badge/wwwangxc/gopkg.svg?size=small)](https://www.murphysec.com/dr/c1TuOdJ62DzT0agLwg)

gopkg/concurrency is an concurrency helper.

## Install

```sh
go get github.com/wwwangxc/gopkg/concurrenty
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/wwwangxc/gopkg/concurrency"
)

func main() {
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
```
