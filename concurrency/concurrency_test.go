package concurrency

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	type args struct {
		ctx         context.Context
		handlers    []Handler
		concurrenty uint8
	}
	tests := []struct {
		name string
		args args
		want *Result
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				handlers: []Handler{
					&handlerNormal{},
					&handlerNormal{},
					&handlerNormal{},
					&handlerNormal{},
					&handlerNormal{},
				},
				concurrenty: 3,
			},
			want: &Result{
				failed: false,
				resultSet: []*singleResult{
					{result: "success", err: nil},
					{result: "success", err: nil},
					{result: "success", err: nil},
					{result: "success", err: nil},
					{result: "success", err: nil},
				},
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				handlers: []Handler{
					&handlerError{},
					&handlerError{},
					&handlerError{},
					&handlerError{},
					&handlerError{},
				},
				concurrenty: 3,
			},
			want: &Result{
				failed: true,
				resultSet: []*singleResult{
					{result: nil, err: fmt.Errorf("error message")},
					{result: nil, err: fmt.Errorf("error message")},
					{result: nil, err: fmt.Errorf("error message")},
					{result: nil, err: fmt.Errorf("error message")},
					{result: nil, err: fmt.Errorf("error message")},
				},
			},
		},
		{
			name: "panic",
			args: args{
				ctx: context.Background(),
				handlers: []Handler{
					&handlerPanic{},
					&handlerPanic{},
					&handlerPanic{},
					&handlerPanic{},
					&handlerPanic{},
				},
				concurrenty: 3,
			},
			want: &Result{
				failed: true,
				resultSet: []*singleResult{
					{result: nil, err: fmt.Errorf("[PANIC]panic message")},
					{result: nil, err: fmt.Errorf("[PANIC]panic message")},
					{result: nil, err: fmt.Errorf("[PANIC]panic message")},
					{result: nil, err: fmt.Errorf("[PANIC]panic message")},
					{result: nil, err: fmt.Errorf("[PANIC]panic message")},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Start(tt.args.ctx, tt.args.handlers, tt.args.concurrenty)
			assert.Equal(t, tt.want.failed, got.failed)
			assert.Equal(t, len(tt.want.resultSet), len(got.resultSet))
			for i, v := range got.resultSet {
				wantRet := tt.want.resultSet[i]
				assert.Equal(t, wantRet.result, v.result)
				assert.True(t, errors.Is(wantRet.err, v.err) || wantRet.err.Error() == v.err.Error())
			}
		})
	}
}

type handlerPanic struct{}

func (e *handlerPanic) Invoke(ctx context.Context) (interface{}, error) {
	panic("panic message")
}

type handlerError struct{}

func (e *handlerError) Invoke(ctx context.Context) (interface{}, error) {
	return nil, fmt.Errorf("error message")
}

type handlerNormal struct{}

func (e *handlerNormal) Invoke(ctx context.Context) (interface{}, error) {
	return "success", nil
}
