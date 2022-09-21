package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	r1, exist := serviceConfigMap["redis_1"]
	assert.True(t, exist, "redis_1 should exist")
	assert.Equal(t, "redis_1", r1.Name)
	assert.Equal(t, "redis://username:password@127.0.0.1:6379/1?timeout=1000ms", r1.DSN)
	assert.Equal(t, 20, r1.MaxIdle)
	assert.Equal(t, 100, r1.MaxActive)
	assert.Equal(t, 1000, r1.MaxConnLifetime)
	assert.Equal(t, 180000, r1.IdleTimeout)
	assert.Equal(t, 1000, r1.Timeout)

	r2, exist := serviceConfigMap["redis_2"]
	assert.True(t, exist, "redis_2 should exist")
	assert.Equal(t, "redis_2", r2.Name)
	assert.Equal(t, "redis://username:password@127.0.0.1:6379/2?timeout=1000ms", r2.DSN)
	assert.Equal(t, 22, r2.MaxIdle)
	assert.Equal(t, 111, r2.MaxActive)
	assert.Equal(t, 2000, r2.MaxConnLifetime)
	assert.Equal(t, 200000, r2.IdleTimeout)
	assert.Equal(t, 2000, r2.Timeout)
}
