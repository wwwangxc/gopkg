package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
)

func Test_loader_Load(t *testing.T) {
	defaultConfig := defaultConfigure("path")
	defaultConfig.watcher = nil

	type args struct {
		path string
		opts []LoadOption
	}
	tests := []struct {
		name    string
		l       *loader
		args    args
		want    Configure
		wantErr bool
		loadErr error
	}{
		{
			name: "unmarshaler not exist",
			l:    newLoader(),
			args: args{
				opts: []LoadOption{
					WithUnmarshaler("not exist unmarshaler"),
					withTest(),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "hit cache",
			l: &loader{
				m: map[string]Configure{
					"path:yaml": defaultConfig,
				},
			},
			args: args{
				path: "path",
				opts: []LoadOption{
					withTest(),
				},
			},
			want:    defaultConfig,
			wantErr: false,
		},
		{
			name: "load fail",
			l:    newLoader(),
			args: args{
				path: "path",
				opts: []LoadOption{
					withTest(),
				},
			},
			want:    nil,
			wantErr: true,
			loadErr: fmt.Errorf(""),
		},
		{
			name: "normal",
			l:    newLoader(),
			args: args{
				path: "path",
				opts: []LoadOption{
					withTest(),
				},
			},
			want:    defaultConfig,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *configureImpl
			patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "Load",
				func(*configureImpl) error {
					return tt.loadErr
				})
			defer patches.Reset()

			got, err := tt.l.Load(tt.args.path, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("loader.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loader.Load() = %v, want %v", got, tt.want)
			}
			if !tt.wantErr {
				key := fmt.Sprintf("%s:%s", tt.args.path, defaultConfigure(tt.args.path).unmarshaler.Name())
				if !reflect.DeepEqual(tt.l.m[key], got) {
					t.Errorf("loader.m = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
