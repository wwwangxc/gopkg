package mysql

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
)

func Test_getDB(t *testing.T) {
	dbs = map[string]*sql.DB{
		"123": {},
	}
	type args struct {
		name string
		opts []Option
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		newDBErr error
	}{
		{
			name:    "db exist",
			wantErr: false,
			args: args{
				name: "123",
			},
		},
		{
			name:     "new db fail",
			wantErr:  true,
			newDBErr: fmt.Errorf(""),
			args: args{
				name: "test",
			},
		},
		{
			name:    "normal",
			wantErr: false,
			args: args{
				name: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyFunc(newDB,
				func(*serviceConfig) (*sql.DB, error) {
					return &sql.DB{}, tt.newDBErr
				})
			defer patches.Reset()

			_, err := getDB(tt.args.name, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_newDB(t *testing.T) {
	dbs = map[string]*sql.DB{
		"123": {},
	}
	type args struct {
		cfg *serviceConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		openErr error
	}{
		{
			name:    "db exsit",
			wantErr: false,
			args: args{
				cfg: &serviceConfig{
					Name: "123",
				},
			},
		},
		{
			name:    "open fail",
			wantErr: true,
			openErr: fmt.Errorf(""),
			args: args{
				cfg: &serviceConfig{
					Name: "test",
				},
			},
		},
		{
			name:    "normal",
			wantErr: false,
			args: args{
				cfg: &serviceConfig{
					Name: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyFunc(sql.Open,
				func(string, string) (*sql.DB, error) {
					return &sql.DB{}, tt.openErr
				})
			defer patches.Reset()
			_, err := newDB(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("newDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				_, exist := dbs[tt.args.cfg.Name]
				assert.True(t, exist)
			}
		})
	}
}
