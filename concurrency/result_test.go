package concurrency

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult_MergedError(t *testing.T) {
	type fields struct {
		resultSet []*singleResult
		failed    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name:    "return nil err",
			wantErr: nil,
			fields: fields{
				failed: false,
			},
		},
		{
			name:    "normal",
			wantErr: fmt.Errorf("2 errors occurreed:\n    * error message\n    * error message"),
			fields: fields{
				resultSet: []*singleResult{
					{err: fmt.Errorf("error message")},
					{err: nil},
					{err: fmt.Errorf("error message")},
				},
				failed: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{
				resultSet: tt.fields.resultSet,
				failed:    tt.fields.failed,
			}
			err := r.MergedError()
			assert.True(t, tt.wantErr == err || !errors.Is(tt.wantErr, err) || err.Error() != tt.wantErr.Error())
		})
	}
}

func TestResult_Errors(t *testing.T) {
	type fields struct {
		resultSet []*singleResult
		failed    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []error
	}{
		{
			name: "no failed result",
			want: nil,
			fields: fields{
				failed: false,
			},
		},
		{
			name: "normal",
			want: []error{
				fmt.Errorf("error message"),
				fmt.Errorf("error message"),
			},
			fields: fields{
				resultSet: []*singleResult{
					{err: fmt.Errorf("error message")},
					{err: nil},
					{err: fmt.Errorf("error message")},
				},
				failed: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{
				resultSet: tt.fields.resultSet,
				failed:    tt.fields.failed,
			}
			got := r.Errors()
			assert.Equal(t, len(tt.want), len(got))
			for i, v := range got {
				assert.True(t, errors.Is(tt.want[i], v) || tt.want[i].Error() == v.Error())
			}
		})
	}
}
