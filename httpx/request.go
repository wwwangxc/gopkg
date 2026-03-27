package httpx

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"resty.dev/v3"
)

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

func checkRequestProtocol(proto RequestProtocol) error {
	switch {
	case proto == nil:
		return fmt.Errorf("request protocol is nil")
	case proto.Host() == "":
		return fmt.Errorf("request host is empty")

	default:
	}

	return nil
}

func buildRequest(ctx context.Context, cli *resty.Client, opts ...RequestOption) (*resty.Request, error) {
	if cli == nil {
		return nil, fmt.Errorf("client proxy is nil")
	}

	r := cli.R().
		SetContext(ctx).
		SetExpectResponseContentType("application/json;charset=utf-8")

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

func protocolToOptions(ctx context.Context, proto RequestProtocol) []RequestOption {
	if proto == nil {
		return nil
	}

	var opts []RequestOption
	opts = append(opts, protocolToAuthOptions(ctx, proto)...)
	opts = append(opts, protocolToCookieOptions(ctx, proto)...)
	opts = append(opts, protocolToHeadersOptions(ctx, proto)...)
	opts = append(opts, protocolToParamsOptions(ctx, proto)...)
	opts = append(opts, protocolToTimeoutOptions(ctx, proto)...)
	opts = append(opts, protocolToRetryOptions(ctx, proto)...)
	opts = append(opts, protocolToSecurityOptions(ctx, proto)...)

	return opts
}

func protocolToAuthOptions(_ context.Context, proto RequestProtocol) []RequestOption {
	if proto == nil {
		return nil
	}

	var opts []RequestOption

	// Try set basic authentication
	if getter, ok := proto.(BasicAuthGetter); ok {
		opts = append(opts, R.WithBasicAuth(getter.BasicAuth()))
	}

	// Try set auth token
	if getter, ok := proto.(BearerTokenAuthGetter); ok {
		opts = append(opts, R.WithBearerTokenAuth(getter.BearerTokenAuth()))
	}

	return opts
}

func protocolToCookieOptions(_ context.Context, proto RequestProtocol) []RequestOption {
	if proto == nil {
		return nil
	}

	var opts []RequestOption

	// Try set cookies
	if getter, ok := proto.(CookiesGetter); ok {
		opts = append(opts, R.WithCookies(getter.Cookies()...))
	}

	return opts
}

func protocolToHeadersOptions(_ context.Context, proto RequestProtocol) []RequestOption {
	if proto == nil {
		return nil
	}

	var opts []RequestOption

	// Try set Headers
	if getter, ok := proto.(HeaderGetter); ok {
		opts = append(opts, R.WithHeader(getter.Header()))
	}

	return opts
}

func protocolToParamsOptions(_ context.Context, proto RequestProtocol) []RequestOption {
	if proto == nil {
		return nil
	}

	var opts []RequestOption

	// Try set path params
	if getter, ok := proto.(PathParamsGetter); ok {
		opts = append(opts, R.WithPathParams(getter.PathParams()))
	}

	// Try set query params
	if getter, ok := proto.(QueryParamsGetter); ok {
		opts = append(opts, R.WithQueryParams(getter.QueryParams()))
	}

	// Try set query string
	if getter, ok := proto.(QueryStringGetter); ok {
		opts = append(opts, R.WithQueryString(getter.QueryString()))
	}

	// Try set form params
	if getter, ok := proto.(FormParamsGetter); ok {
		opts = append(opts, R.WithFormParams(getter.FormParams()))
	}

	// Try set body
	if getter, ok := proto.(BodyGetter); ok {
		opts = append(opts, R.WithBody(getter.Body()))
	}

	return opts
}

func protocolToTimeoutOptions(_ context.Context, proto RequestProtocol) []RequestOption {
	if proto == nil {
		return nil
	}

	var opts []RequestOption

	// Try set timeout
	if getter, ok := proto.(TimeoutGetter); ok {
		opts = append(opts, R.WithTimeout(getter.Timeout()))
	}

	return opts
}

func protocolToRetryOptions(ctx context.Context, proto RequestProtocol) []RequestOption {
	if proto == nil {
		return nil
	}

	var opts []RequestOption

	// Try set retry
	if getter, ok := proto.(RetryGetter); ok {
		times, conds := getter.Retry()
		opts = append(opts, R.WithRetry(times, conds...))
	}

	// Try set retry wait
	if getter, ok := proto.(RetryWaitGetter); ok {
		opts = append(opts, R.WithRetryWait(getter.RetryWait()))
	}

	// Try set retry hooks
	if getter, ok := proto.(RetryHooksGetter); ok {
		opts = append(opts, R.WithRetryHooks(getter.RetryHooks(ctx)...))
	}

	return opts
}

func protocolToSecurityOptions(_ context.Context, proto RequestProtocol) []RequestOption {
	if proto == nil {
		return nil
	}

	var opts []RequestOption

	if _, ok := proto.(AllowResponseBodyUnlimitedReadsGetter); ok {
		opts = append(opts, R.AllowResponseBodyUnlimitedReads())
	}

	if _, ok := proto.(AllowMethodGetPayloadGetter); ok {
		opts = append(opts, R.AllowMethodGetPayload())
	}

	if _, ok := proto.(AllowMethodDeletePayloadGetter); ok {
		opts = append(opts, R.AllowMethodDeletePayload())
	}

	return opts
}
