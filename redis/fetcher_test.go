package redis

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey"
	redigo "github.com/gomodule/redigo/redis"
)

func Test_fetcherImpl_fetch(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		dest interface{}
		opts []FetchOption
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		bytesErr error
		doErr    error
	}{
		{
			name:     "bytes fail",
			wantErr:  true,
			bytesErr: fmt.Errorf(""),
		},
		{
			name:     "callback fail",
			wantErr:  true,
			bytesErr: redigo.ErrNil,
			args: args{
				opts: []FetchOption{
					WithFetchCallback(func() (interface{}, error) { return nil, fmt.Errorf("") }, time.Microsecond),
				},
			},
		},
		{
			name:     "marshal fail",
			wantErr:  true,
			bytesErr: redigo.ErrNil,
			args: args{
				opts: []FetchOption{
					WithFetchCallback(func() (interface{}, error) { return nil, nil }, time.Microsecond),
					WithFetchMarshal(func(v interface{}) ([]byte, error) { return []byte{}, fmt.Errorf("") }),
				},
			},
		},
		{
			name:     "do fail",
			wantErr:  true,
			bytesErr: redigo.ErrNil,
			doErr:    fmt.Errorf(""),
			args: args{
				opts: []FetchOption{
					WithFetchCallback(func() (interface{}, error) { return nil, nil }, time.Microsecond),
				},
			},
		},
		{
			name:    "normal process",
			wantErr: false,
			args: args{
				dest: &map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyFunc(Bytes,
				func(interface{}, error) ([]byte, error) {
					return []byte("{}"), tt.bytesErr
				})
			defer patches.Reset()

			var cli *clientProxyImpl
			patches.ApplyMethod(reflect.TypeOf(cli), "Do",
				func(*clientProxyImpl, context.Context, string, ...interface{}) (interface{}, error) {
					return nil, tt.doErr
				})

			f := &fetcherImpl{
				name: "client_name",
			}
			if err := f.fetch(tt.args.ctx, tt.args.key, tt.args.dest, newFetchOptions(tt.args.opts...)); (err != nil) != tt.wantErr {
				t.Errorf("fetcherImpl.Fetch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
