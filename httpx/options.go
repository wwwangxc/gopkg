package httpx

import (
	"net/http"
	"time"

	"resty.dev/v3"
)

// C is the client proxy option provider
var C = (*clientOption)(nil)

// ClientOption client proxy option
type ClientOption func(*resty.Client)

type clientOption struct{}

// WithHost set api host
func (*clientOption) WithHost(host ...string) ClientOption {
	return func(c *resty.Client) {
		if c == nil {
			return
		}

		rr, err := resty.NewRoundRobin(host...)
		if err != nil {
			panic(err)
		}

		c.SetLoadBalancer(rr)
	}
}

// WithTimeout set timeout.
//
// It can be overridden at the request level. See [R.WithTimeout] or [TimeoutGetter]
func (*clientOption) WithTimeout(timeout time.Duration) ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.SetTimeout(timeout)
		}
	}
}

// WithHeader set http header
//
// It can be overridden at the request level. See [R.WithHeader] or [HeaderGetter]
func (*clientOption) WithHeader(header map[string]string) ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.SetHeaders(header)
		}
	}
}

// WithTransport set http transport
func (*clientOption) WithTransport(transport http.RoundTripper) ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.SetTransport(transport)
		}
	}
}

// WithRequestMiddlewares method allows Resty users to override the default request
// middlewares sequence
func (*clientOption) WithRequestMiddlewares(middlewares ...resty.RequestMiddleware) ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.SetRequestMiddlewares(middlewares...)
		}
	}
}

// WithResponseMiddlewares method allows Resty users to override the default response
// middlewares sequence
func (*clientOption) WithResponseMiddlewares(middlewares ...resty.ResponseMiddleware) ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.SetResponseMiddlewares(middlewares...)
		}
	}
}

// OnSuccess method adds a callback that will be run whenever a request execution
// succeeds.  This is called after all retries have been attempted (if any).
//
// Out of the [resty.Client.OnSuccess], [resty.Client.OnError], [resty.Client.OnInvalid], [resty.Client.OnPanic]
// callbacks, exactly one set will be invoked for each call to [RequestProtocol.Execute] that completes.
//
// NOTE:
//   - Do not use [C] setter methods within OnSuccess hooks; deadlock will happen.
func (*clientOption) OnSuccess(hook resty.SuccessHook) ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.OnSuccess(hook)
		}
	}
}

// OnError method adds a callback that will be run whenever a request execution fails.
// This is called after all retries have been attempted (if any).
// If there was a response from the server, the error will be wrapped in [ResponseError]
// which has the last response received from the server.
//
// Out of the [resty.Client.OnSuccess], [resty.Client.OnError], [resty.Client.OnInvalid], [resty.Client.OnPanic]
// callbacks, exactly one set will be invoked for each call to [RequestProtocol.Execute] that completes.
//
// NOTE:
//   - Do not use [C] setter methods within OnError hooks; deadlock will happen.
func (*clientOption) OnError(hook resty.ErrorHook) ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.OnError(hook)
		}
	}
}

// OnInvalid method adds a callback that will be run whenever a request execution
// fails before it starts because the request is invalid.
//
// Out of the [resty.Client.OnSuccess], [resty.Client.OnError], [resty.Client.OnInvalid], [resty.Client.OnPanic]
// callbacks, exactly one set will be invoked for each call to [RequestProtocol.Execute] that completes.
//
// NOTE:
//   - Do not use [C] setter methods within OnInvalid hooks; deadlock will happen.
func (*clientOption) OnInvalid(hook resty.ErrorHook) ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.OnInvalid(hook)
		}
	}
}

// OnPanic method adds a callback that will be run whenever a request execution
// panics.
//
// Out of the [resty.Client.OnSuccess], [resty.Client.OnError], [resty.Client.OnInvalid], [resty.Client.OnPanic]
// callbacks, exactly one set will be invoked for each call to [RequestProtocol.Execute] that completes.
//
// If an [Client.OnSuccess], [Client.OnError], or [Client.OnInvalid] callback panics,
// then exactly one rule can be violated.
//
// NOTE:
//   - Do not use [C] setter methods within OnPanic hooks; deadlock will happen.
func (*clientOption) OnPanic(hook resty.ErrorHook) ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.OnPanic(hook)
		}
	}
}

// WithAllowMethodGetPayload method allows the GET method with payload on the Resty client.
// By default, Resty does not allow.
func (*clientOption) WithAllowMethodGetPayload() ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.AllowMethodGetPayload()
		}
	}
}

// WithAllowMethodDeletePayload method allows the DELETE method with payload on the Resty client.
// By default, Resty does not allow.
func (*clientOption) WithAllowMethodDeletePayload() ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.AllowMethodDeletePayload()
		}
	}
}

// WithTrace method enables the Resty client trace for the requests fired from
// the client using [httptrace.ClientTrace] and provides insights.
func (*clientOption) WithTrace() ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.EnableTrace()
		}
	}
}

// WithDebug method enables the debug mode on the Resty client. The client logs details
// of every request and response.
func (*clientOption) WithDebug() ClientOption {
	return func(c *resty.Client) {
		if c != nil {
			c.EnableDebug()
		}
	}
}

// R is the request option provider
var R = (*requestOption)(nil)

// RequestOption request option
type RequestOption func(r *resty.Request)

type requestOption struct{}

// WithBasicAuth set username & password in basic authentication
func (s *requestOption) WithBasicAuth(username, password string) RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetBasicAuth(username, password)
		}
	}
}

// WithBearerTokenAuth set authentication token
func (s *requestOption) WithBearerTokenAuth(token string) RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetAuthToken(token)
		}
	}
}

// WithCookies append the cookies for current request
func (s *requestOption) WithCookies(cookies ...*http.Cookie) RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetCookies(cookies)
		}
	}
}

// WithHeader set the headers for current request
func (s *requestOption) WithHeader(header map[string]string) RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetHeaders(header)
		}
	}
}

// WithPathParams set the parameters for current request path
func (s *requestOption) WithPathParams(params map[string]string) RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetPathParams(params)
		}
	}
}

// WithQueryParams set the query params for current request
func (s *requestOption) WithQueryParams(params map[string]string) RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetQueryParams(params)
		}
	}
}

// WithQueryString set the query string for current request
func (s *requestOption) WithQueryString(query string) RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetQueryString(query)
		}
	}
}

// WithFormParams set the form params for current request
func (s *requestOption) WithFormParams(params map[string]string) RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetFormData(params)
		}
	}
}

// WithBody set the body for current request
func (s *requestOption) WithBody(body any) RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetBody(body)
		}
	}
}

// WithTimeout set the timeout for current request
func (s *requestOption) WithTimeout(timeout time.Duration) RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetTimeout(timeout)
		}
	}
}

// WithRetry set retry strategy for current request
func (s *requestOption) WithRetry(times int, conds ...resty.RetryConditionFunc) RequestOption {
	return func(r *resty.Request) {
		if r == nil {
			return
		}

		if len(conds) == 0 {
			conds = []resty.RetryConditionFunc{
				// Retry if the error is not empty.
				func(r *resty.Response, err error) bool { return err != nil },
			}
		}

		r.SetRetryCount(times).
			SetRetryConditions(conds...).
			SetAllowNonIdempotentRetry(true)
	}
}

// WithRetryWait set the wait time before retry sleep for current request
func (s *requestOption) WithRetryWait(wait time.Duration) RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetRetryWaitTime(wait)
		}
	}
}

// WithRetryHooks set retry hooks for current request
func (s *requestOption) WithRetryHooks(hooks ...resty.RetryHookFunc) RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetRetryHooks(hooks...)
		}
	}
}

// AllowResponseBodyUnlimitedReads enable the response body in memory that provides an ability to do unlimited reads
//
// Unlimited reads are possible in a few scenarios, even without enabling it.
//   - When debug mode is enabled
//
// NOTE: Use with care
//   - Turning on this feature keeps the response body in memory, which might cause additional memory usage.
func (s *requestOption) AllowResponseBodyUnlimitedReads() RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetResponseBodyUnlimitedReads(true)
		}
	}
}

// AllowMethodGetPayload allows the GET method with payload on the Resty client
func (s *requestOption) AllowMethodGetPayload() RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetAllowMethodGetPayload(true)
		}
	}
}

// AllowMethodDeletePayload allows the DELETE method with payload on the Resty client
func (s *requestOption) AllowMethodDeletePayload() RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.SetAllowMethodDeletePayload(true)
		}
	}
}

// WithDebug enables the debug mode on the current request. It logs
// the details current request and response.
func (s *requestOption) WithDebug() RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.EnableDebug()
		}
	}
}

// WithTrace enables trace for the current request
func (s *requestOption) WithTrace() RequestOption {
	return func(r *resty.Request) {
		if r != nil {
			r.EnableTrace()
		}
	}
}
