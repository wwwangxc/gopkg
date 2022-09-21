package etcd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	e1, exist := clientConfigMap["etcd1"]
	assert.True(t, exist, "etcd1 should exist")
	assert.Equal(t, "etcd1", e1.Name)
	assert.Equal(t, 1000, e1.Timeout)
	assert.Equal(t, "username", e1.Username)
	assert.Equal(t, "password", e1.Password)
	assert.Equal(t, "/usr/local/etcd_conf/key.pem", e1.TLSKeyPath)
	assert.Equal(t, "/usr/local/etcd_conf/cert.pem", e1.TLSCertPath)
	assert.Equal(t, "/usr/local/etcd_conf/cacert.pem", e1.CACertPath)
}
