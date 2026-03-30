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

func ExampleNewClientProxy() {
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

func ExampleClientProxy_Do() {
	// Request options to be set automatically:
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
	//		- [AllowMethodGetPayloadGetter]           // allows the GET method with payload.
	//		- [AllowMethodDeletePayloadGetter]        // allows the DELETE method with payload.
	//		- [DebugGetter]                           // enable debug mode.
	//		- [TraceGetter]                           // enable trace for current request.
	type MyRequest struct {
		httpx.RequestProtocol // must implement [httpx.RequestProtocol]
	}

	// Create request protocol
	req := &MyRequest{}

	// Create HTTP/HTTPS client with config
	cli := httpx.NewClientProxy("name")

	// Send HTTP/HTTPS request with request protocol
	var rsp map[string]any
	_ = cli.Do(context.Background(), req, &rsp)
}

func ExampleClientProxy_Get() {
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
		httpx.R.WithDebug(),
		httpx.R.WithTrace(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Get(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}

func ExampleClientProxy_Head() {
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
		httpx.R.WithDebug(),
		httpx.R.WithTrace(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Head(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}

func ExampleClientProxy_Post() {
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
		httpx.R.WithDebug(),
		httpx.R.WithTrace(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Post(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}

func ExampleClientProxy_Put() {
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
		httpx.R.WithDebug(),
		httpx.R.WithTrace(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Put(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}

func ExampleClientProxy_Delete() {
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
		httpx.R.WithDebug(),
		httpx.R.WithTrace(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Delete(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}

func ExampleClientProxy_Options() {
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
		httpx.R.WithDebug(),
		httpx.R.WithTrace(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Options(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}

func ExampleClientProxy_Patch() {
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
		httpx.R.WithDebug(),
		httpx.R.WithTrace(),
	}

	// Send request
	var rsp map[string]any
	_ = cli.Patch(context.Background(), "https://httpbin.org/anything", &rsp, requestOpts...)
}
