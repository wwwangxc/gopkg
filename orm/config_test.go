package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	cfg, exist := serviceConfigMap["db_mysql"]
	assert.True(t, exist, "db_mysql should exist")
	assert.Equal(t, "db_mysql", cfg.Name)
	assert.Equal(t, "root:root@tcp(127.0.0.1:3306)/db1?charset=utf8&parseTime=True", cfg.DSN)
	assert.Equal(t, "mysql", cfg.Driver)
	assert.Equal(t, 111, cfg.MaxIdle)
	assert.Equal(t, 222, cfg.MaxOpen)
	assert.Equal(t, 333, cfg.MaxIdleTime)

	cfg, exist = serviceConfigMap["db_postgresql"]
	assert.True(t, exist, "db_postgresql should exist")
	assert.Equal(t, "db_postgresql", cfg.Name)
	assert.Equal(t, "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai", cfg.DSN)
	assert.Equal(t, "postgresql", cfg.Driver)
	assert.Equal(t, 11, cfg.MaxIdle)
	assert.Equal(t, 22, cfg.MaxOpen)
	assert.Equal(t, 33, cfg.MaxIdleTime)

	cfg, exist = serviceConfigMap["db_sqlite"]
	assert.True(t, exist, "db_sqlite should exist")
	assert.Equal(t, "db_sqlite", cfg.Name)
	assert.Equal(t, "gopkg.db", cfg.DSN)
	assert.Equal(t, "sqlite", cfg.Driver)
	assert.Equal(t, 11, cfg.MaxIdle)
	assert.Equal(t, 22, cfg.MaxOpen)
	assert.Equal(t, 33, cfg.MaxIdleTime)

	cfg, exist = serviceConfigMap["db_sqlserver"]
	assert.True(t, exist, "db_sqlserver should exist")
	assert.Equal(t, "db_sqlserver", cfg.Name)
	assert.Equal(t, "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm", cfg.DSN)
	assert.Equal(t, "sqlserver", cfg.Driver)
	assert.Equal(t, 11, cfg.MaxIdle)
	assert.Equal(t, 22, cfg.MaxOpen)
	assert.Equal(t, 33, cfg.MaxIdleTime)

	cfg, exist = serviceConfigMap["db_clickhouse"]
	assert.True(t, exist, "db_clickhouse should exist")
	assert.Equal(t, "db_clickhouse", cfg.Name)
	assert.Equal(t, "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20", cfg.DSN)
	assert.Equal(t, "clickhouse", cfg.Driver)
	assert.Equal(t, 11, cfg.MaxIdle)
	assert.Equal(t, 22, cfg.MaxOpen)
	assert.Equal(t, 33, cfg.MaxIdleTime)
}
