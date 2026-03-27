package httpx

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/wwwangxc/wheel/reflectx"
)

// ClientProxy http client proxy
//
//go:generate mockgen -source=client.go -destination=mockhttpx/client_mock.go -package=mockhttpx
type ClientProxy interface {
	// Do http request by protocol
	Do(ctx context.Context, req RequestProtocol, dest any) error

	// Get will do http GET request
	//
	// URL format such as:
	//	          Normal: http://127.0.0.1:8080/v1/users/details
	//	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
	Get(ctx context.Context, url string, dest any, opts ...RequestOption) error

	// Head will do http HEAD request
	//
	// URL format such as:
	// 	          Normal: http://127.0.0.1:8080/v1/users/details
	// 	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
	Head(ctx context.Context, url string, dest any, opts ...RequestOption) error

	// Post will do http POST request
	//
	// URL format such as:
	// 	          Normal: http://127.0.0.1:8080/v1/users/details
	// 	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
	Post(ctx context.Context, url string, dest any, opts ...RequestOption) error

	// Put will do http PUT request
	//
	// URL format such as:
	// 	          Normal: http://127.0.0.1:8080/v1/users/details
	// 	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
	Put(ctx context.Context, url string, dest any, opts ...RequestOption) error

	// Delete will do http DELETE request
	//
	// URL format such as:
	// 	          Normal: http://127.0.0.1:8080/v1/users/details
	// 	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
	Delete(ctx context.Context, url string, dest any, opts ...RequestOption) error

	// Options will do http OPTIONS request
	//
	// URL format such as:
	// 	          Normal: http://127.0.0.1:8080/v1/users/details
	// 	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
	Options(ctx context.Context, url string, dest any, opts ...RequestOption) error

	// Patch will do http PATCH request
	//
	// URL format such as:
	// 	          Normal: http://127.0.0.1:8080/v1/users/details
	// 	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
	Patch(ctx context.Context, url string, dest any, opts ...RequestOption) error
}

type clientProxyImpl struct {
	name string
	opts []ClientOption
}

// NewClientProxy new http client proxy
func NewClientProxy(name string, opt ...ClientOption) ClientProxy {
	return &clientProxyImpl{
		name: name,
		opts: opt,
	}
}

// Do http request by protocol
func (s *clientProxyImpl) Do(ctx context.Context, req RequestProtocol, dest any) error {
	if s == nil {
		return fmt.Errorf("invalid client proxy")
	}

	if err := checkRequestProtocol(req); err != nil {
		return err
	}

	method := http.MethodGet
	if getter, ok := req.(MethodGetter); ok {
		method = getter.Method()
	}

	host := strings.TrimSuffix(req.Host(), "/")
	path := strings.TrimPrefix(req.Path(), "/")
	url := fmt.Sprintf("%s/%s", host, path)
	opts := protocolToOptions(ctx, req)

	return s.do(ctx, method, url, dest, opts...)
}

// Get will do http GET request
//
// URL format such as:
//
//	          Normal: http://127.0.0.1:8080/v1/users/details
//	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
func (s *clientProxyImpl) Get(ctx context.Context, url string, dest any, opts ...RequestOption) error {
	return s.do(ctx, http.MethodGet, url, dest, opts...)
}

// Head will do http HEAD request
//
// URL format such as:
//
//	          Normal: http://127.0.0.1:8080/v1/users/details
//	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
func (s *clientProxyImpl) Head(ctx context.Context, url string, dest any, opts ...RequestOption) error {
	return s.do(ctx, http.MethodHead, url, dest, opts...)
}

// Post will do http POST request
//
// URL format such as:
//
//	          Normal: http://127.0.0.1:8080/v1/users/details
//	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
func (s *clientProxyImpl) Post(ctx context.Context, url string, dest any, opts ...RequestOption) error {
	return s.do(ctx, http.MethodPost, url, dest, opts...)
}

// Put will do http PUT request
//
// URL format such as:
//
//	          Normal: http://127.0.0.1:8080/v1/users/details
//	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
func (s *clientProxyImpl) Put(ctx context.Context, url string, dest any, opts ...RequestOption) error {
	return s.do(ctx, http.MethodPut, url, dest, opts...)
}

// Delete will do http DELETE request
//
// URL format such as:
//
//	          Normal: http://127.0.0.1:8080/v1/users/details
//	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
func (s *clientProxyImpl) Delete(ctx context.Context, url string, dest any, opts ...RequestOption) error {
	return s.do(ctx, http.MethodDelete, url, dest, opts...)
}

// Options will do http OPTIONS request
//
// URL format such as:
//
//	          Normal: http://127.0.0.1:8080/v1/users/details
//	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
func (s *clientProxyImpl) Options(ctx context.Context, url string, dest any, opts ...RequestOption) error {
	return s.do(ctx, http.MethodOptions, url, dest, opts...)
}

// Patch will do http PATCH request
//
// URL format such as:
//
//	          Normal: http://127.0.0.1:8080/v1/users/details
//	With path params: http://127.0.0.1:8080/v1/users/{user_id}/repo/{repo_id}/star
func (s *clientProxyImpl) Patch(ctx context.Context, url string, dest any, opts ...RequestOption) error {
	return s.do(ctx, http.MethodPatch, url, dest, opts...)
}

func (s *clientProxyImpl) do(ctx context.Context, method string, url string, dest any, opts ...RequestOption) error {
	if s == nil {
		return fmt.Errorf("invalid client proxy")
	}

	if !reflectx.IsNil(dest) {
		tp := reflect.TypeOf(dest)
		if tp.Kind() != reflect.Pointer {
			return fmt.Errorf("httpx.Request(non-pointer %s)", tp.Kind().String())
		}
	}

	cli := getOrCreateClient(s.name, s.opts...)
	r, err := buildRequest(ctx, cli, opts...)
	if err != nil {
		return err
	}

	if _, err = r.SetResult(dest).Execute(method, url); err != nil {
		return err
	}

	return nil
}
