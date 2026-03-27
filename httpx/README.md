# gopkg/httpx

gopkg/httpx is a componentized http plugin.

It providels:

- An easy way to configre and manage http client.
- Custom request protocol.

Based on [resty.dev/v3](https://resty.dev)

Required go1.23

## Contents

- [gopkg/httpx](#gopkg%2Fhttpx)
    - [Install](#install)
    - [Quick Start](#quick-start)
      - [Config](#config)
      - [ClientProxy](#clientproxy)
          - [Do Request With Protocol](#do-request-with-protocol)
          - [Get](#get)
          - [Head](#head)
          - [Post](#post)
          - [Put](#put)
          - [Delete](#delete)
          - [Options](#options)
          - [Patch](#patch)
      - [RequestProtocol](#requestprotocol)
          - [MethodGetter](#methodgetter)
          - [BasicAuthGetter](#basicauthgetter)
          - [BearerTokenAuthGetter](#bearertokenauthgetter)
          - [CookiesGetter](#cookiesgetter)
          - [HeaderGetter](#headergetter)
          - [PathParamsGetter](#pathparamsgetter)
          - [QueryParamsGetter](#queryparamsgetter)
          - [QueryStringGetter](#querystringgetter)
          - [FormParamsGetter](#formparamsgetter)
          - [BodyGetter](#bodygetter)
          - [TimeoutGetter](#timeoutgetter)
          - [RetryGetter](#retrygetter)
          - [RetryWaitGetter](#retrywaitgetter)
          - [RetryHooksGetter](#retryhooksgetter)
          - [AllowResponseBodyUnlimitedReadsGetter](#allowresponsebodyunlimitedreadsgetter)
          - [AllowMethodGetPayloadGetter](#allowmethodgetpayloadgetter)
          - [AllowMethodDeletePayloadGetter](#allowmethoddeletepayloadgetter)
    - [How To Mock](#how-to-mock)

## Install

```sh
go get github.com/wwwangxc/gopkg/httpx
```

**[⬆ back to top](#contents)**

## Quick Start

```go
package httpx_test

import (
	"context"

	// gopkg/httpx will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/httpx"
)

func Example() {
	// Load config (optional)
	// gopkg/httpx will automatically read configuration
	// files (./app.yaml) when package loaded
	httpx.LoadConfig("./custom_config.yaml")

	// Create HTTP/HTTPS client with config
	cli := httpx.NewClientProxy("name")

	// Send HTTP/HTTPX GET request
	var rsp map[string]any
	_ = cli.Get(context.Background(), "https://httpbin.org/anything", &rsp)
}
```

**[⬆ back to top](#contents)**

### Config

```yaml
client:
  http: 
    header:
      User-Agent: gopkg/httpx
      Content-Type: application/json;charset=UTF-8
    transport:
      max_idle_conns: 200
      max_idle_conns_per_host: 50
      max_conns_per_host: 100
      idle_conn_timeout: 90s
      tls_handshake_timeout: 3s
      expect_continue_timeout: 1s
      response_header_timeout: 5s
      dial:
        timeout: 3s
        keep_alive: 30s
  service:
    - name: http1
      dsn: https://httpbin.org
      timeout: 3000
    - name: http2
      dsn: https://httpbin1.org,https://httpbin2.org
      timeout: 1000
      header:
        User-Agent: custom_agent
        Content-Type: custom_content_type
      transport:
        max_idle_conns: 1
        max_idle_conns_per_host: 2
        max_conns_per_host: 3
        idle_conn_timeout: 1s
        tls_handshake_timeout: 2s
        expect_continue_timeout: 3s
        response_header_timeout: 4s
        dial:
          timeout: 5s
          keep_alive: 6s
```

**[⬆ back to top](#contents)**

### ClientProxy

```go
package httpx_test

import (
	"net/http"
	"time"

	"resty.dev/v3"

	// gopkg/httpx will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/httpx"
)

func Example() {
	// Create HTTP/HTTPS client proxy with config
	_ = httpx.NewClientProxy("name")

	// Create HTTP/HTTPX client proxy with options
	_ = httpx.NewClientProxy("name",
		httpx.C.WithHost("httpbin1.org", "httpbin2.org"),
		httpx.C.WithTimeout(3*time.Second),
		httpx.C.WithHeader(map[string]string{}),
		httpx.C.WithTransport(&http.Transport{}),
		httpx.C.WithRequestMiddlewares([]resty.RequestMiddleware{}...),
		httpx.C.WithResponseMiddlewares([]resty.ResponseMiddleware{}...),
		httpx.C.OnSuccess(func(c *resty.Client, r *resty.Response) {}),
		httpx.C.OnError(func(r *resty.Request, err error) {}),
		httpx.C.OnInvalid(func(r *resty.Request, err error) {}),
		httpx.C.OnPanic(func(r *resty.Request, err error) {}),
		httpx.C.WithAllowMethodGetPayload(),
		httpx.C.WithAllowMethodDeletePayload(),
		httpx.C.WithTrace(),
		httpx.C.WithDebug(),
	)
}
```

**[⬆ back to top](#contents)**

#### Do Request With Protocol

```go
package httpx_test

import (
	"context"
	"fmt"
	"net/http"

	"resty.dev/v3"

	// gopkg/httpx will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/httpx"
)

func Example() {
	// Create request protocol
	req := &MyRequest{}

	// Create HTTP/HTTPS client with config
	cli := httpx.NewClientProxy("name")

	// Send request with protocol
	var rsp map[string]any
	_ = cli.Do(context.Background(), req, &rsp)
}

// MyRequest my request protocol
//
// By implementing the appropriate Getter interface for the request protocol,
// various options in the request can be set automatically.
//
//	Support getter:
//
//		- [MethodGetter]                          // sets the http method
//		- [BasicAuthGetter]                       // sets the basic authentication header
//		- [BearerTokenAuthGetter]                 // sets the auth token header
//		- [CookiesGetter]                         // append cookies
//		- [HeaderGetter]                          // set multiple header fields and their values
//		- [PathParamsGetter]                      // set multiple URL path key-value pair
//		- [QueryParamsGetter]                     // set multiple parameter
//		- [QueryStringGetter]                     // set the query string
//		- [FormParamsGetter]                      // set form parameters and their values
//		- [BodyGetter]                            // set the request body
//		- [TimeoutGetter]                         // set the request timeout
//		- [RetryGetter]                           // set the retry times and conditions when the request failed
//		- [RetryWaitGetter]                       // set the default wait time for sleep before retrying
//		- [RetryHooksGetter]                      // set retry hooks
//		- [AllowResponseBodyUnlimitedReadsGetter] // enable the response body in memory that provides an ability to do unlimited reads.
//		- [AllowMethodGetPayloadGetter]           // allows the GET method with payload on the Resty client.
//		- [AllowMethodDeletePayloadGetter]        // allows the DELETE method with payload on the Resty client.
type MyRequest struct{}

func (s *MyRequest) Host() string {
	return "https://httpbin.org"
}

func (s *MyRequest) Path() string {
	return "/anything"
}

func (s *MyRequest) AllowResponseBodyUnlimitedReads() {}

func (s *MyRequest) AllowMethodGetPayload() {}

func (s *MyRequest) AllowMethodDeletePayload() {}

func (s *MyRequest) Method() string {
	return http.MethodGet
}

func (s *MyRequest) BasicAuth() (username string, password string) {
	return "username", "password"
}

func (s *MyRequest) BearerTokenAuth() string {
	return "token"
}

func (s *MyRequest) Cookies() []*http.Cookie {
	return []*http.Cookie{
		{Name: "cookie_1"},
		{Name: "cookie_2"},
	}
}

func (s *MyRequest) Header() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}

func (s *MyRequest) PathParams() map[string]string {
	return map[string]string{
		"key": "value",
	}
}

func (s *MyRequest) QueryParams() map[string]string {
	return map[string]string{
		"key": "value",
	}
}

func (s *MyRequest) QueryString() string {
	return "key_1=value_1&key_2=value_2"
}

func (s *MyRequest) FormParams() map[string]string {
	return map[string]string{
		"key": "value",
	}
}

func (s *MyRequest) Body() any {
	return map[string]any{
		"key": "value",
	}
}

func (s *MyRequest) Timeout() time.Duration {
	return 3 * time.Second
}

func (s *MyRequest) Retry() (retryTimes int, retryConds []resty.RetryConditionFunc) {
	return 3, []resty.RetryConditionFunc{
		func(r *resty.Response, err error) bool { return err != nil },
		func(r *resty.Response, err error) bool { return r.StatusCode() != http.StatusOK },
	}
}

func (s *MyRequest) RetryWait() time.Duration {
	return 100 * time.Millisecond
}

func (s *MyRequest) RetryHooks(ctx context.Context) []resty.RetryHookFunc {
	return []resty.RetryHookFunc{
		func(r *resty.Response, err error) { fmt.Println("retry once") },
	}
}
```

**[⬆ back to top](#contents)**

#### Get

```go
package httpx_test

import (
	"context"
	"net/http"
	"time"

	"resty.dev/v3"

	// gopkg/httpx will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/httpx"
)

func Example() {
	// Create HTTP/HTTPS client with config
	cli := httpx.NewClientProxy("name")

	// Request supported options
	requestOpts := []httpx.RequestOption{
		httpx.R.WithBasicAuth("username", "password"),
		httpx.R.WithBearerTokenAuth("token"),
		httpx.R.WithCookies([]*http.Cookie{}...),
		httpx.R.WithHeader(map[string]string{}),
		httpx.R.WithPathParams(map[string]string{}),
		httpx.R.WithQueryParams(map[string]string{}),
		httpx.R.WithQueryString(""),
		httpx.R.WithFormParams(map[string]string{}),
		httpx.R.WithBody(map[string]any{}),
		httpx.R.WithTimeout(3 * time.Second),
		httpx.R.WithRetry(3, []resty.RetryConditionFunc{}...),
		httpx.R.WithRetryWait(100 * time.Millisecond),
		httpx.R.WithRetryHooks([]resty.RetryHookFunc{}...),
		httpx.R.AllowResponseBodyUnlimitedReads(),
		httpx.R.AllowMethodGetPayload(),
		httpx.R.AllowMethodDeletePayload(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Get(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}
```

**[⬆ back to top](#contents)**

#### Head

```go
package httpx_test

import (
	"context"
	"net/http"
	"time"

	"resty.dev/v3"

	// gopkg/httpx will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/httpx"
)

func Example() {
	// Create HTTP/HTTPS client with config
	cli := httpx.NewClientProxy("name")

	// Request supported options
	requestOpts := []httpx.RequestOption{
		httpx.R.WithBasicAuth("username", "password"),
		httpx.R.WithBearerTokenAuth("token"),
		httpx.R.WithCookies([]*http.Cookie{}...),
		httpx.R.WithHeader(map[string]string{}),
		httpx.R.WithPathParams(map[string]string{}),
		httpx.R.WithQueryParams(map[string]string{}),
		httpx.R.WithQueryString(""),
		httpx.R.WithFormParams(map[string]string{}),
		httpx.R.WithBody(map[string]any{}),
		httpx.R.WithTimeout(3 * time.Second),
		httpx.R.WithRetry(3, []resty.RetryConditionFunc{}...),
		httpx.R.WithRetryWait(100 * time.Millisecond),
		httpx.R.WithRetryHooks([]resty.RetryHookFunc{}...),
		httpx.R.AllowResponseBodyUnlimitedReads(),
		httpx.R.AllowMethodGetPayload(),
		httpx.R.AllowMethodDeletePayload(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Head(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}
```

**[⬆ back to top](#contents)**

#### Post

```go
package httpx_test

import (
	"context"
	"net/http"
	"time"

	"resty.dev/v3"

	// gopkg/httpx will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/httpx"
)

func Example() {
	// Create HTTP/HTTPS client with config
	cli := httpx.NewClientProxy("name")

	// Request supported options
	requestOpts := []httpx.RequestOption{
		httpx.R.WithBasicAuth("username", "password"),
		httpx.R.WithBearerTokenAuth("token"),
		httpx.R.WithCookies([]*http.Cookie{}...),
		httpx.R.WithHeader(map[string]string{}),
		httpx.R.WithPathParams(map[string]string{}),
		httpx.R.WithQueryParams(map[string]string{}),
		httpx.R.WithQueryString(""),
		httpx.R.WithFormParams(map[string]string{}),
		httpx.R.WithBody(map[string]any{}),
		httpx.R.WithTimeout(3 * time.Second),
		httpx.R.WithRetry(3, []resty.RetryConditionFunc{}...),
		httpx.R.WithRetryWait(100 * time.Millisecond),
		httpx.R.WithRetryHooks([]resty.RetryHookFunc{}...),
		httpx.R.AllowResponseBodyUnlimitedReads(),
		httpx.R.AllowMethodGetPayload(),
		httpx.R.AllowMethodDeletePayload(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Post(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}
```

**[⬆ back to top](#contents)**

#### Put

```go
package httpx_test

import (
	"context"
	"net/http"
	"time"

	"resty.dev/v3"

	// gopkg/httpx will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/httpx"
)

func Example() {
	// Create HTTP/HTTPS client with config
	cli := httpx.NewClientProxy("name")

	// Request supported options
	requestOpts := []httpx.RequestOption{
		httpx.R.WithBasicAuth("username", "password"),
		httpx.R.WithBearerTokenAuth("token"),
		httpx.R.WithCookies([]*http.Cookie{}...),
		httpx.R.WithHeader(map[string]string{}),
		httpx.R.WithPathParams(map[string]string{}),
		httpx.R.WithQueryParams(map[string]string{}),
		httpx.R.WithQueryString(""),
		httpx.R.WithFormParams(map[string]string{}),
		httpx.R.WithBody(map[string]any{}),
		httpx.R.WithTimeout(3 * time.Second),
		httpx.R.WithRetry(3, []resty.RetryConditionFunc{}...),
		httpx.R.WithRetryWait(100 * time.Millisecond),
		httpx.R.WithRetryHooks([]resty.RetryHookFunc{}...),
		httpx.R.AllowResponseBodyUnlimitedReads(),
		httpx.R.AllowMethodGetPayload(),
		httpx.R.AllowMethodDeletePayload(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Put(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}
```

**[⬆ back to top](#contents)**

#### Delete

```go
package httpx_test

import (
	"context"
	"net/http"
	"time"

	"resty.dev/v3"

	// gopkg/httpx will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/httpx"
)

func Example() {
	// Create HTTP/HTTPS client with config
	cli := httpx.NewClientProxy("name")

	// Request supported options
	requestOpts := []httpx.RequestOption{
		httpx.R.WithBasicAuth("username", "password"),
		httpx.R.WithBearerTokenAuth("token"),
		httpx.R.WithCookies([]*http.Cookie{}...),
		httpx.R.WithHeader(map[string]string{}),
		httpx.R.WithPathParams(map[string]string{}),
		httpx.R.WithQueryParams(map[string]string{}),
		httpx.R.WithQueryString(""),
		httpx.R.WithFormParams(map[string]string{}),
		httpx.R.WithBody(map[string]any{}),
		httpx.R.WithTimeout(3 * time.Second),
		httpx.R.WithRetry(3, []resty.RetryConditionFunc{}...),
		httpx.R.WithRetryWait(100 * time.Millisecond),
		httpx.R.WithRetryHooks([]resty.RetryHookFunc{}...),
		httpx.R.AllowResponseBodyUnlimitedReads(),
		httpx.R.AllowMethodGetPayload(),
		httpx.R.AllowMethodDeletePayload(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Delete(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}
```

**[⬆ back to top](#contents)**

#### Options

```go
package httpx_test

import (
	"context"
	"net/http"
	"time"

	"resty.dev/v3"

	// gopkg/httpx will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/httpx"
)

func Example() {
	// Create HTTP/HTTPS client with config
	cli := httpx.NewClientProxy("name")

	// Request supported options
	requestOpts := []httpx.RequestOption{
		httpx.R.WithBasicAuth("username", "password"),
		httpx.R.WithBearerTokenAuth("token"),
		httpx.R.WithCookies([]*http.Cookie{}...),
		httpx.R.WithHeader(map[string]string{}),
		httpx.R.WithPathParams(map[string]string{}),
		httpx.R.WithQueryParams(map[string]string{}),
		httpx.R.WithQueryString(""),
		httpx.R.WithFormParams(map[string]string{}),
		httpx.R.WithBody(map[string]any{}),
		httpx.R.WithTimeout(3 * time.Second),
		httpx.R.WithRetry(3, []resty.RetryConditionFunc{}...),
		httpx.R.WithRetryWait(100 * time.Millisecond),
		httpx.R.WithRetryHooks([]resty.RetryHookFunc{}...),
		httpx.R.AllowResponseBodyUnlimitedReads(),
		httpx.R.AllowMethodGetPayload(),
		httpx.R.AllowMethodDeletePayload(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Options(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}
```

**[⬆ back to top](#contents)**

#### Patch

```go
package httpx_test

import (
	"context"
	"net/http"
	"time"

	"resty.dev/v3"

	// gopkg/httpx will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/gopkg/httpx"
)

func Example() {
	// Create HTTP/HTTPS client with config
	cli := httpx.NewClientProxy("name")

	// Request supported options
	requestOpts := []httpx.RequestOption{
		httpx.R.WithBasicAuth("username", "password"),
		httpx.R.WithBearerTokenAuth("token"),
		httpx.R.WithCookies([]*http.Cookie{}...),
		httpx.R.WithHeader(map[string]string{}),
		httpx.R.WithPathParams(map[string]string{}),
		httpx.R.WithQueryParams(map[string]string{}),
		httpx.R.WithQueryString(""),
		httpx.R.WithFormParams(map[string]string{}),
		httpx.R.WithBody(map[string]any{}),
		httpx.R.WithTimeout(3 * time.Second),
		httpx.R.WithRetry(3, []resty.RetryConditionFunc{}...),
		httpx.R.WithRetryWait(100 * time.Millisecond),
		httpx.R.WithRetryHooks([]resty.RetryHookFunc{}...),
		httpx.R.AllowResponseBodyUnlimitedReads(),
		httpx.R.AllowMethodGetPayload(),
		httpx.R.AllowMethodDeletePayload(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Patch(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}
```

**[⬆ back to top](#contents)**

### Request Protocol

```go
// RequestProtocol custom request protocol
//
// By implementing the appropriate Getter interface for the request protocol,
// various options in the request can be set automatically.
//
//	Support getter:
//
//		- [MethodGetter]                          // sets the http method
//		- [BasicAuthGetter]                       // sets the basic authentication header
//		- [BearerTokenAuthGetter]                 // sets the auth token header
//		- [CookiesGetter]                         // append cookies
//		- [HeaderGetter]                          // set multiple header fields and their values
//		- [PathParamsGetter]                      // set multiple URL path key-value pair
//		- [QueryParamsGetter]                     // set multiple parameter
//		- [QueryStringGetter]                     // set the query string
//		- [FormParamsGetter]                      // set form parameters and their values
//		- [BodyGetter]                            // set the request body
//		- [TimeoutGetter]                         // set the request timeout
//		- [RetryGetter]                           // set the retry times and conditions when the request failed
//		- [RetryWaitGetter]                       // set the default wait time for sleep before retrying
//		- [RetryHooksGetter]                      // set retry hooks
//		- [AllowResponseBodyUnlimitedReadsGetter] // enable the response body in memory that provides an ability to do unlimited reads.
//		- [AllowMethodGetPayloadGetter]           // allows the GET method with payload on the Resty client.
//		- [AllowMethodDeletePayloadGetter]        // allows the DELETE method with payload on the Resty client.
type RequestProtocol interface {
	Host() string
	Path() string
}
```

**[⬆ back to top](#contents)**

#### MethodGetter

```go
// MethodGetter returns http method
//
// Implement this interface to set the method for current current HTTP request.
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) Method() { return http.MethodGet }
type MethodGetter interface {
	Method() string
}
```

**[⬆ back to top](#contents)**

#### BasicAuthGetter

```go
// BasicAuth returns username & password in basic authentication
//
// Implement this interface to automatically set the Basic Authentication header in the current HTTP request.
//
// Header Format:
//
//	Authorization: Basic <base64-encoded-value>
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) BasicAuth() (username, password string) {
//		return "username", "password"
//	}
type BasicAuthGetter interface {
	BasicAuth() (username, password string)
}
```

**[⬆ back to top](#contents)**

#### BearerTokenAuthGetter

```go
// BearerTokenAuthGetter returns authentication token
//
// Implement this interface to automatically set the Bearer Authentication header in the current HTTP request.
//
// Header Format:
//
//	Authorization: Bearer <token>
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) BearerTokenAuth() string {
//		return "token"
//	}
type BearerTokenAuthGetter interface{ BearerTokenAuth() string }
```

**[⬆ back to top](#contents)**

#### CookieGetters

```go
// CookieGetters returns the cookies of the current request
//
// Implement this interface to automatically append the cookies in the current HTTP request.
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) Cookies() []*http.Cookie {
//		return []*http.Cookie{}
//	}
type CookiesGetter interface{ Cookies() []*http.Cookie }
```

**[⬆ back to top](#contents)**

#### HeaderGetter

```go
// HeaderGetter returns the headers of the current request
//
// Implement this interface to automatically set headers in the current HTTP request.
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) Header() map[string]string {
//		return map[string]string{
//			"Content-Type": "application/json;charset=UTF-8",
//		}
//	}
type HeaderGetter interface{ Header() map[string]string }
```

**[⬆ back to top](#contents)**

#### PathParamsGetter

```go
// PathParamsGetter returns the parameters in the current request path
//
// Implement this interface to automatically set the parameters in the current HTTP request path.
//
// URL Format:
//
//	   Raw: /api/path/{params_name}
//	Cooked: /api/path/666
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) PathParams() map[string]string {
//		return map[string]string{
//			"params_name": "666",
//		}
//	}
type PathParamsGetter interface{ PathParams() map[string]string }
```

**[⬆ back to top](#contents)**

#### QueryParamsGetter

```go
// QueryParamsGetter returns the query parameters in the current request
//
// Implement this interface to automatically set the query parameters in the current HTTP request.
//
// URL Format:
//
//	/api/path?params_1=value_1&params_2=value_2
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) QueryParams() map[string]string {
//		return map[string]string{
//			"params_1": "value_1",
//			"params_2": "value_2",
//		}
//	}
type QueryParamsGetter interface{ QueryParams() map[string]string }
```

**[⬆ back to top](#contents)**

#### QueryStringGetter

```go
// QueryStringGetter returns the query parameters string in the current request
//
// Implement this interface to automatically set the query parameters string in the current HTTP request.
//
// URL Format:
//
//	/api/path?params_1=value_1&params_2=value_2
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) QueryString() string {
//		return "params_1=value_1&params_2=value_2"
//	}
type QueryStringGetter interface{ QueryString() string }
```

**[⬆ back to top](#contents)**

#### FormParamsGetter

```go
// FormParamsGetter returns the form parameters in the current request
//
// Implement this interface to automatically set the form parameters in the current HTTP request.
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) FormParams() map[string]string {
//		return map[string]string{
//			"params_name": "666",
//		}
//	}
type FormParamsGetter interface{ FormParams() map[string]string }
```

**[⬆ back to top](#contents)**

#### BodyGetter

```go
// BodyGetter returns the body in the current request
//
// Implement this interface to automatically set the body in the current HTTP request.
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) Body() any {
//		// support string
//		return `{"params": "value"}`
//
//		// support []byte
//		return []byte("This is my raw request")
//
//		// support map
//		return map[string]any {
//			"params": "value",
//		}
//
//		// support struct
//		return Request {
//			Params: "value",
//		}
//	}
type BodyGetter interface{ Body() any }
```

**[⬆ back to top](#contents)**

#### TimeoutGetter

```go
// TimeoutGetter returns timeout for current request
//
// Implement this interface to automatically set the timeout for current HTTP request.
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) Timeout() any {
//		return 3 * time.Second
//	}
type TimeoutGetter interface{ Timeout() time.Duration }
```

**[⬆ back to top](#contents)**

#### RetryGetter

```go
// RetryGetter returns retry strategy for current request
//
// Implement this interface to automatically set the retry strategy for current HTTP request.
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) Retry() (retryTimes int, retryConds []resty.RetryConditionFunc) {
//		// retry 3 times when response code is http.StatusTooManyRequests or http.StatusBadRequest
//		return 3, []resty.RetryConditionFunc {
//			httpx.RetryWithStatusCodes(http.StatusTooManyRequests, http.StatusBadRequest),
//		}
//
//		// retry 3 times when returns error
//		return 3, nil
//	}
type RetryGetter interface {
	Retry() (retryTimes int, retryConds []resty.RetryConditionFunc)
}
```

**[⬆ back to top](#contents)**

#### RetryWaitGetter

```go
// RetryWaitGetter returns the wait time before retry sleep for current request
//
// Implement this interface to automatically set the wait time before retry sleep for current HTTP request.
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) RetryWait() time.Duration {
//		return time.Second
//	}
type RetryWaitGetter interface{ RetryWait() time.Duration }
```

**[⬆ back to top](#contents)**

#### RetryHooksGetter

```go
// RetryHooksGetter returns retry hooks for current request
//
// Implement this interface to automatically set the retry hooks for current HTTP request.
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) RetryHooks(ctx context.Context) []resty.RetryHookFunc {
//		return []resty.RetryHookFunc{
//			func(r *resty.Response, err error) {
//				if err != nil {
//					fmt.Printf("retry once, because an error occurred, error:%+v", err)
//					return
//				}
//
//				if r != nil && r.Result() != nil {
//					fmt.Printf("retry once, response:%+v", r.Result())
//				}
//			},
//		}
//	}
type RetryHooksGetter interface {
	RetryHooks(ctx context.Context) []resty.RetryHookFunc
}
```

**[⬆ back to top](#contents)**

#### AllowResponseBodyUnlimitedReadsGetter

```go
// AllowResponseBodyUnlimitedReadsGetter return nothing
//
// Implement this interface for enable the response body in memory that provides an ability to do unlimited reads.
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) AllowResponseBodyUnlimitedReadsGetter() {}
//
// Unlimited reads are possible in a few scenarios, even without enabling it.
//   - When debug mode is enabled
//
// NOTE: Use with care
//   - Turning on this feature keeps the response body in memory, which might cause additional memory usage.
type AllowResponseBodyUnlimitedReadsGetter interface{ AllowResponseBodyUnlimitedReads() }
```

**[⬆ back to top](#contents)**

#### AllowMethodGetPayloadGetter

```go
// AllowMethodGetPayloadGetter return nothing
//
// Implement this interface will allows the GET method with payload on the Resty client.
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) AllowMethodGetPayloadGetter() {}
type AllowMethodGetPayloadGetter interface{ AllowMethodGetPayload() }
```

**[⬆ back to top](#contents)**

#### AllowMethodDeletePayloadGetter

```go
// AllowMethodDeletePayloadGetter return nothing
//
// Implement this interface will allows the DELETE method with payload on the Resty client.
//
// Usage Example:
//
//	type MyRequest struct{}
//
//	func (s *MyRequest) AllowMethodDeletePayloadGetter() {}
type AllowMethodDeletePayloadGetter interface{ AllowMethodDeletePayload() }
```

**[⬆ back to top](#contents)**

## How To Mock

```go
package httpx_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"go.uber.org/mock/gomock"

	"github.com/wwwangxc/gopkg/httpx"
	"github.com/wwwangxc/gopkg/httpx/mockhttpx"
)

func TestNewClientProxy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client proxy
	mockedCli := mockhttpx.NewMockClientProxy(ctrl)
	mockedCli.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Head(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Options(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Mock function httpx.NewClientProxy
	patches := gomonkey.ApplyFunc(httpx.NewClientProxy,
		func(string, ...httpx.ClientOption) httpx.ClientProxy {
			return mockedCli
		})
	defer patches.Reset()

	// dosomething...
}
```

**[⬆ back to top](#contents)**
