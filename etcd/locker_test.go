package etcd

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func Test_newLockerProxy(t *testing.T) {
	type args struct {
		cli    *clientv3.Client
		prefix string
		ttl    int
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		newSessionErr error
	}{
		{
			name:          "new session fail",
			wantErr:       true,
			newSessionErr: fmt.Errorf(""),
		},
		{
			name:    "normal process",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyFunc(concurrency.NewSession,
				func(*clientv3.Client, ...concurrency.SessionOption) (*concurrency.Session, error) {
					return nil, tt.newSessionErr
				})
			defer patches.Reset()

			_, err := newLockerProxy(tt.args.cli, tt.args.prefix, tt.args.ttl)
			if (err != nil) != tt.wantErr {
				t.Errorf("newLockerProxy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_lockerProxyImpl_TryLock(t *testing.T) {
	type fields struct {
		session *concurrency.Session
		mutex   *concurrency.Mutex
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		tryLockErr error
		gotErr     error
	}{
		{
			name:       "lock not acquired",
			wantErr:    true,
			tryLockErr: concurrency.ErrLocked,
			gotErr:     ErrLockNotAcquired,
		},
		{
			name:       "try lock fail",
			wantErr:    true,
			tryLockErr: fmt.Errorf(""),
			gotErr:     fmt.Errorf(""),
		},
		{
			name:    "normal process",
			wantErr: false,
			gotErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mutex *concurrency.Mutex
			patches := gomonkey.ApplyMethod(reflect.TypeOf(mutex), "TryLock",
				func(*concurrency.Mutex, context.Context) error {
					return tt.tryLockErr
				})
			defer patches.Reset()

			l := &lockerProxyImpl{
				session: tt.fields.session,
				mutex:   tt.fields.mutex,
			}
			err := l.TryLock(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("lockerProxyImpl.TryLock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.gotErr, err)
		})
	}
}
