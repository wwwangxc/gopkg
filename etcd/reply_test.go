package etcd

import (
	"fmt"
	"reflect"
	"testing"

	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestPutResult(t *testing.T) {
	type args struct {
		in0 *clientv3.PutResponse
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "return error",
			wantErr: true,
			args: args{
				err: fmt.Errorf(""),
			},
		},
		{
			name:    "return nil",
			wantErr: false,
			args:    args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PutResult(tt.args.in0, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("PutResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetResult(t *testing.T) {
	type args struct {
		resp *clientv3.GetResponse
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "error",
			wantErr: true,
			want:    nil,
			args: args{
				err: fmt.Errorf(""),
			},
		},
		{
			name:    "normal",
			wantErr: false,
			want: map[string]string{
				"key": "value",
			},
			args: args{
				resp: &clientv3.GetResponse{
					Kvs: []*mvccpb.KeyValue{
						{Key: []byte("key"), Value: []byte("value")},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResult(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResultString(t *testing.T) {
	type args struct {
		resp *clientv3.GetResponse
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "error",
			wantErr: true,
			want:    "",
			args: args{
				err: fmt.Errorf(""),
			},
		},
		{
			name:    "error nil",
			wantErr: true,
			want:    "",
			args: args{
				resp: &clientv3.GetResponse{
					Count: 0,
					Kvs: []*mvccpb.KeyValue{
						{Key: []byte("key"), Value: []byte("value")},
					},
				},
			},
		},
		{
			name:    "normal",
			wantErr: false,
			want:    "value",
			args: args{
				resp: &clientv3.GetResponse{
					Count: 1,
					Kvs: []*mvccpb.KeyValue{
						{Key: []byte("key"), Value: []byte("value")},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResultString(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResultString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetResultString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResultStrings(t *testing.T) {
	type args struct {
		resp *clientv3.GetResponse
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "error",
			wantErr: true,
			want:    nil,
			args: args{
				err: fmt.Errorf(""),
			},
		},
		{
			name:    "error nil",
			wantErr: true,
			want:    nil,
			args: args{
				resp: &clientv3.GetResponse{
					Count: 0,
					Kvs: []*mvccpb.KeyValue{
						{Key: []byte("key"), Value: []byte("value")},
					},
				},
			},
		},
		{
			name:    "normal",
			wantErr: false,
			want:    []string{"value"},
			args: args{
				resp: &clientv3.GetResponse{
					Count: 1,
					Kvs: []*mvccpb.KeyValue{
						{Key: []byte("key"), Value: []byte("value")},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResultStrings(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResultStrings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResultStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResultInt(t *testing.T) {
	type args struct {
		resp *clientv3.GetResponse
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "error",
			wantErr: true,
			want:    0,
			args: args{
				err: fmt.Errorf(""),
			},
		},
		{
			name:    "error nil",
			wantErr: true,
			want:    0,
			args: args{
				resp: &clientv3.GetResponse{
					Count: 0,
					Kvs: []*mvccpb.KeyValue{
						{Key: []byte("key"), Value: []byte("1")},
					},
				},
			},
		},
		{
			name:    "normal",
			wantErr: false,
			want:    1,
			args: args{
				resp: &clientv3.GetResponse{
					Count: 1,
					Kvs: []*mvccpb.KeyValue{
						{Key: []byte("key"), Value: []byte("1")},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResultInt(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResultInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetResultInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResultInts(t *testing.T) {
	type args struct {
		resp *clientv3.GetResponse
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name:    "error",
			wantErr: true,
			want:    nil,
			args: args{
				err: fmt.Errorf(""),
			},
		},
		{
			name:    "error nil",
			wantErr: true,
			want:    nil,
			args: args{
				resp: &clientv3.GetResponse{
					Count: 0,
					Kvs: []*mvccpb.KeyValue{
						{Key: []byte("key"), Value: []byte("1")},
					},
				},
			},
		},
		{
			name:    "normal",
			wantErr: false,
			want:    []int{1},
			args: args{
				resp: &clientv3.GetResponse{
					Count: 1,
					Kvs: []*mvccpb.KeyValue{
						{Key: []byte("key"), Value: []byte("1")},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResultInts(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResultInts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResultInts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResultCount(t *testing.T) {
	type args struct {
		resp *clientv3.GetResponse
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "error",
			wantErr: true,
			want:    0,
			args: args{
				err: fmt.Errorf(""),
			},
		},
		{
			name:    "normal",
			wantErr: false,
			want:    1,
			args: args{
				resp: &clientv3.GetResponse{
					Count: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResultCount(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResultCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetResultCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteResult(t *testing.T) {
	type args struct {
		resp *clientv3.DeleteResponse
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "return error",
			wantErr: true,
			want:    0,
			args: args{
				err: fmt.Errorf(""),
			},
		},
		{
			name:    "normal",
			wantErr: false,
			want:    1,
			args: args{
				resp: &clientv3.DeleteResponse{
					Deleted: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteResult(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTxnResult(t *testing.T) {
	type args struct {
		in0 *clientv3.TxnResponse
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "return error",
			wantErr: true,
			args: args{
				err: fmt.Errorf(""),
			},
		},
		{
			name:    "normal process",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := TxnResult(tt.args.in0, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("TxnResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLeaseGrantResult(t *testing.T) {
	type args struct {
		resp *clientv3.LeaseGrantResponse
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    clientv3.LeaseID
		wantErr bool
	}{
		{
			name:    "return error",
			wantErr: true,
			want:    0,
			args: args{
				err: fmt.Errorf(""),
			},
		},
		{
			name:    "normal",
			wantErr: false,
			want:    1,
			args: args{
				resp: &clientv3.LeaseGrantResponse{
					ID: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LeaseGrantResult(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("LeaseGrantResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LeaseGrantResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLeaseRevokeResult(t *testing.T) {
	type args struct {
		in0 *clientv3.LeaseRevokeResponse
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "return error",
			wantErr: true,
			args: args{
				err: fmt.Errorf(""),
			},
		},
		{
			name:    "normal process",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LeaseRevokeResult(tt.args.in0, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("LeaseRevokeResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLeaseTimeToLiveResult(t *testing.T) {
	type args struct {
		resp *clientv3.LeaseTimeToLiveResponse
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "return error",
			wantErr: true,
			want:    0,
			args: args{
				err: fmt.Errorf(""),
			},
		},
		{
			name:    "normal",
			wantErr: false,
			want:    1,
			args: args{
				resp: &clientv3.LeaseTimeToLiveResponse{
					TTL: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LeaseTimeToLiveResult(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("LeaseTimeToLiveResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LeaseTimeToLiveResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
