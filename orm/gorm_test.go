package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"gorm.io/gorm"

	"github.com/wwwangxc/gopkg/orm/driver"
)

func Test_getGORMDB(t *testing.T) {
	gormDBMap = map[string]*gorm.DB{
		"123": {},
	}
	type args struct {
		name string
		opts []GORMOption
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		newGORMDBErr error
	}{
		{
			name:    "db exist",
			wantErr: false,
			args: args{
				name: "123",
			},
		},
		{
			name:         "new gorm db fail",
			wantErr:      true,
			newGORMDBErr: fmt.Errorf(""),
			args: args{
				name: "abc",
			},
		},
		{
			name:    "normal",
			wantErr: false,
			args: args{
				name: "abc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyFunc(newGORMDB,
				func(*serviceConfig) (*gorm.DB, error) {
					return &gorm.DB{}, tt.newGORMDBErr
				})
			defer patches.Reset()

			_, err := getGORMDB(tt.args.name, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("getGORMDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_newGORMDB(t *testing.T) {
	type args struct {
		cfg *serviceConfig
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		openErr  error
		getDBErr error
	}{
		{
			name:    "invalid driver",
			wantErr: true,
			args: args{
				cfg: &serviceConfig{
					Driver: "test driver",
				},
			},
		},
		{
			name:    "open fail",
			wantErr: true,
			openErr: fmt.Errorf(""),
			args: args{
				cfg: &serviceConfig{
					Driver: driver.NameMySQL,
				},
			},
		},
		{
			name:     "get db fail",
			wantErr:  true,
			getDBErr: fmt.Errorf(""),
			args: args{
				cfg: &serviceConfig{
					Driver: driver.NameMySQL,
				},
			},
		},
		{
			name:    "normal",
			wantErr: false,
			args: args{
				cfg: &serviceConfig{
					Name:   "test",
					Driver: driver.NameMySQL,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyFunc(gorm.Open,
				func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error) {
					return &gorm.DB{}, tt.openErr
				})
			defer patches.Reset()

			var db *gorm.DB
			patches.ApplyMethod(reflect.TypeOf(db), "DB",
				func(*gorm.DB) (*sql.DB, error) {
					return &sql.DB{}, tt.getDBErr
				})
			_, err := newGORMDB(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("newGORMDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
