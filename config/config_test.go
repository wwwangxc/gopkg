package config

import (
	"reflect"
	"testing"

	"github.com/wwwangxc/gopkg/config/unmarshaler"
)

func Test_configureImpl_Unmarshal(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	tests := []struct {
		name     string
		wantErr  bool
		instance *configureImpl
	}{
		{
			name:     "unmarshaler not exist",
			wantErr:  true,
			instance: &configureImpl{},
		},
		{
			name:     "normal",
			wantErr:  false,
			instance: c,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := map[string]interface{}{}
			if err := tt.instance.Unmarshal(&out); (err != nil) != tt.wantErr {
				t.Errorf("configureImpl.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configureImpl_IsExist(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "subkey exist",
			args: args{
				k: "subkey.string_value",
			},
			want: true,
		},
		{
			name: "subkey not exist",
			args: args{
				k: "subkey.not exist key",
			},
			want: false,
		},
		{
			name: "exist",
			args: args{
				k: "string_value",
			},
			want: true,
		},
		{
			name: "not exist",
			args: args{
				k: "not exist key",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.IsExist(tt.args.k); got != tt.want {
				t.Errorf("configureImpl.IsExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_Get(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k          string
		defaultVal interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "return default value",
			args: args{
				k:          "not exist key",
				defaultVal: "default value",
			},
			want: "default value",
		},
		{
			name: "normal",
			args: args{
				k:          "string_value",
				defaultVal: "default value",
			},
			want: "string value",
		},
		{
			name: "subkey return default value",
			args: args{
				k:          "subkey.not exist key",
				defaultVal: "default value",
			},
			want: "default value",
		},
		{
			name: "subkey normal",
			args: args{
				k:          "subkey.string_value",
				defaultVal: "default value",
			},
			want: "string value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Get(tt.args.k, tt.args.defaultVal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configureImpl.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_GetString(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k          string
		defaultVal string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "return default value",
			args: args{
				k:          "not exist key",
				defaultVal: "default value",
			},
			want: "default value",
		},
		{
			name: "normal",
			args: args{
				k:          "string_value",
				defaultVal: "default value",
			},
			want: "string value",
		},
		{
			name: "subkey return default value",
			args: args{
				k:          "subkey.not exist key",
				defaultVal: "default value",
			},
			want: "default value",
		},
		{
			name: "subkey normal",
			args: args{
				k:          "subkey.string_value",
				defaultVal: "default value",
			},
			want: "string value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GetString(tt.args.k, tt.args.defaultVal); got != tt.want {
				t.Errorf("configureImpl.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_GetBool(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k          string
		defaultVal bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "return default value",
			args: args{
				k:          "not exist key",
				defaultVal: false,
			},
			want: false,
		},
		{
			name: "normal",
			args: args{
				k:          "bool_value",
				defaultVal: true,
			},
			want: true,
		},
		{
			name: "subkey return default value",
			args: args{
				k:          "subkey.not exist key",
				defaultVal: false,
			},
			want: false,
		},
		{
			name: "subkey normal",
			args: args{
				k:          "subkey.bool_value",
				defaultVal: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GetBool(tt.args.k, tt.args.defaultVal); got != tt.want {
				t.Errorf("configureImpl.GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_GetInt(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k          string
		defaultVal int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "return default value",
			args: args{
				k:          "not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "normal",
			args: args{
				k:          "int_value",
				defaultVal: 999,
			},
			want: -1,
		},
		{
			name: "subkey return default value",
			args: args{
				k:          "subkey.not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "subkey normal",
			args: args{
				k:          "subkey.int_value",
				defaultVal: 999,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GetInt(tt.args.k, tt.args.defaultVal); got != tt.want {
				t.Errorf("configureImpl.GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_GetInt32(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k          string
		defaultVal int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "return default value",
			args: args{
				k:          "not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "normal",
			args: args{
				k:          "int32_value",
				defaultVal: 999,
			},
			want: -2,
		},
		{
			name: "subkey return default value",
			args: args{
				k:          "subkey.not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "subkey normal",
			args: args{
				k:          "subkey.int32_value",
				defaultVal: 999,
			},
			want: -2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GetInt32(tt.args.k, tt.args.defaultVal); got != tt.want {
				t.Errorf("configureImpl.GetInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_GetInt64(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k          string
		defaultVal int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "return default value",
			args: args{
				k:          "not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "normal",
			args: args{
				k:          "int64_value",
				defaultVal: 999,
			},
			want: -3,
		},
		{
			name: "subkey return default value",
			args: args{
				k:          "subkey.not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "subkey normal",
			args: args{
				k:          "subkey.int64_value",
				defaultVal: 999,
			},
			want: -3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GetInt64(tt.args.k, tt.args.defaultVal); got != tt.want {
				t.Errorf("configureImpl.GetInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_GetUint(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k          string
		defaultVal uint
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{
			name: "return default value",
			args: args{
				k:          "not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "normal",
			args: args{
				k:          "uint_value",
				defaultVal: 999,
			},
			want: 1,
		},
		{
			name: "subkey return default value",
			args: args{
				k:          "subkey.not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "subkey normal",
			args: args{
				k:          "subkey.uint_value",
				defaultVal: 999,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GetUint(tt.args.k, tt.args.defaultVal); got != tt.want {
				t.Errorf("configureImpl.GetUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_GetUint32(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k          string
		defaultVal uint32
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "return default value",
			args: args{
				k:          "not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "normal",
			args: args{
				k:          "uint32_value",
				defaultVal: 999,
			},
			want: 2,
		},
		{
			name: "subkey return default value",
			args: args{
				k:          "subkey.not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "subkey normal",
			args: args{
				k:          "subkey.uint32_value",
				defaultVal: 999,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GetUint32(tt.args.k, tt.args.defaultVal); got != tt.want {
				t.Errorf("configureImpl.GetUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_GetUint64(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k          string
		defaultVal uint64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "return default value",
			args: args{
				k:          "not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "normal",
			args: args{
				k:          "uint64_value",
				defaultVal: 999,
			},
			want: 3,
		},
		{
			name: "subkey return default value",
			args: args{
				k:          "subkey.not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "subkey normal",
			args: args{
				k:          "subkey.uint64_value",
				defaultVal: 999,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GetUint64(tt.args.k, tt.args.defaultVal); got != tt.want {
				t.Errorf("configureImpl.GetUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_GetFloat32(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k          string
		defaultVal float32
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{
			name: "return default value",
			args: args{
				k:          "not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "normal",
			args: args{
				k:          "float32_value",
				defaultVal: 999,
			},
			want: 1.11,
		},
		{
			name: "subkey return default value",
			args: args{
				k:          "subkey.not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "subkey normal",
			args: args{
				k:          "subkey.float32_value",
				defaultVal: 999,
			},
			want: 1.11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GetFloat32(tt.args.k, tt.args.defaultVal); got != tt.want {
				t.Errorf("configureImpl.GetFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_GetFloat64(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k          string
		defaultVal float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "return default value",
			args: args{
				k:          "not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "normal",
			args: args{
				k:          "float64_value",
				defaultVal: 999,
			},
			want: 2.22,
		},
		{
			name: "subkey return default value",
			args: args{
				k:          "subkey.not exist key",
				defaultVal: 999,
			},
			want: 999,
		},
		{
			name: "subkey normal",
			args: args{
				k:          "subkey.float64_value",
				defaultVal: 999,
			},
			want: 2.22,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GetFloat64(tt.args.k, tt.args.defaultVal); got != tt.want {
				t.Errorf("configureImpl.GetFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_getWithDefaultVal(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k          string
		defaultVal interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "key not exist",
			args: args{
				k:          "not exist key",
				defaultVal: "default value",
			},
			want: "default value",
		},
		{
			name: "unsupport default value type",
			args: args{
				k:          "string_value",
				defaultVal: map[string]interface{}{},
			},
			want: map[string]interface{}{},
		},
		{
			name: "value convert fail",
			args: args{
				k:          "string_value",
				defaultVal: 1,
			},
			want: 1,
		},
		{
			name: "normal",
			args: args{
				k:          "string_value",
				defaultVal: "default value",
			},
			want: "string value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.getWithDefaultVal(tt.args.k, tt.args.defaultVal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configureImpl.getWithDefaultVal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_get(t *testing.T) {
	c := defaultConfigure("./testdata/config.yaml")
	_ = c.Load()

	type args struct {
		k string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "not exist",
			args: args{
				k: "not exist key",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "normal",
			args: args{
				k: "string_value",
			},
			want:    "string value",
			wantErr: false,
		},
		{
			name: "subkey normal",
			args: args{
				k: "subkey.string_value",
			},
			want:    "string value",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.get(tt.args.k)
			if (err != nil) != tt.wantErr {
				t.Errorf("configureImpl.get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configureImpl.get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configureImpl_load(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  bool
		instance *configureImpl
	}{
		{
			name:     "unmarshaler not exist",
			wantErr:  true,
			instance: &configureImpl{},
		},
		{
			name:    "read file fail",
			wantErr: true,
			instance: &configureImpl{
				path:        "./abc.yaml",
				unmarshaler: &unmarshaler.YAML{},
			},
		},
		{
			name:    "unmarshal fail",
			wantErr: true,
			instance: &configureImpl{
				path:        "./testdata/invalid_config.yaml",
				unmarshaler: &unmarshaler.YAML{},
			},
		},
		{
			name:    "normal",
			wantErr: false,
			instance: &configureImpl{
				path:        "./testdata/config.yaml",
				unmarshaler: &unmarshaler.YAML{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.instance.Load(); (err != nil) != tt.wantErr {
				t.Errorf("configureImpl.load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
