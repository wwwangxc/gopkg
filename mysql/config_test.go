package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	cli1, exist := serviceConfigMap["client1"]
	assert.True(t, exist, "client1 should exist")
	assert.Equal(t, "client1", cli1.Name)
	assert.Equal(t, "root:root@tcp(127.0.0.1:3306)/db1?charset=utf8&parseTime=True", cli1.DSN)
	assert.Equal(t, 11, cli1.MaxIdle)
	assert.Equal(t, 22, cli1.MaxOpen)
	assert.Equal(t, 33, cli1.MaxIdleTime)

	cli2, exist := serviceConfigMap["client2"]
	assert.True(t, exist, "client2 should exist")
	assert.Equal(t, "client2", cli2.Name)
	assert.Equal(t, "root:root@tcp(127.0.0.1:3306)/db2?charset=utf8&parseTime=True", cli2.DSN)
	assert.Equal(t, 111, cli2.MaxIdle)
	assert.Equal(t, 222, cli2.MaxOpen)
	assert.Equal(t, 333, cli2.MaxIdleTime)
}
