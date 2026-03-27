package httpx

import (
	"testing"
	"time"

	c "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	c.Convey("Init app config", t, func() {
		c.Convey("Default config", func() {
			c.Convey("With global config", func() {
				c, ok := clientConfigMap["http1"]
				assert.True(t, ok)
				assert.Equal(t, "http1", c.Name)
				assert.Equal(t, "https://httpbin.org", c.DSN)
				assert.Equal(t, int64(3000), c.Timeout)

				userAgent, ok := c.Header["User-Agent"]
				assert.True(t, ok)
				assert.Equal(t, "gopkg/httpx", userAgent)

				contentType, ok := c.Header["Content-Type"]
				assert.True(t, ok)
				assert.Equal(t, "application/json;charset=UTF-8", contentType)

				transport := c.TransportCfg
				assert.NotNil(t, transport)
				assert.Equal(t, 200, transport.MaxIdleConns)
				assert.Equal(t, 50, transport.MaxIdleConnsPerHost)
				assert.Equal(t, 100, transport.MaxConnsPerHost)
				assert.Equal(t, 90*time.Second, transport.IdleConnTimeout)
				assert.Equal(t, 3*time.Second, transport.TLSHandshakeTimeout)
				assert.Equal(t, time.Second, transport.ExpectContinueTimeout)
				assert.Equal(t, 5*time.Second, transport.ResponseHeaderTimeout)
				assert.NotNil(t, 3*time.Second, transport.Dial.Timeout)
				assert.NotNil(t, 30*time.Second, transport.Dial.KeepAlive)
			})

			c.Convey("With custom config", func() {
				c, ok := clientConfigMap["http2"]
				assert.True(t, ok)
				assert.Equal(t, "http2", c.Name)
				assert.Equal(t, "https://httpbin1.org,https://httpbin2.org", c.DSN)
				assert.Equal(t, int64(1000), c.Timeout)

				userAgent, ok := c.Header["User-Agent"]
				assert.True(t, ok)
				assert.Equal(t, "custom_agent", userAgent)

				contentType, ok := c.Header["Content-Type"]
				assert.True(t, ok)
				assert.Equal(t, "custom_content_type", contentType)

				transport := c.TransportCfg
				assert.NotNil(t, transport)
				assert.Equal(t, 1, transport.MaxIdleConns)
				assert.Equal(t, 2, transport.MaxIdleConnsPerHost)
				assert.Equal(t, 3, transport.MaxConnsPerHost)
				assert.Equal(t, 1*time.Second, transport.IdleConnTimeout)
				assert.Equal(t, 2*time.Second, transport.TLSHandshakeTimeout)
				assert.Equal(t, 3*time.Second, transport.ExpectContinueTimeout)
				assert.Equal(t, 4*time.Second, transport.ResponseHeaderTimeout)
				assert.NotNil(t, 5*time.Second, transport.Dial.Timeout)
				assert.NotNil(t, 6*time.Second, transport.Dial.KeepAlive)
			})
		})

		c.Convey("Custom config", func() {
			c.Convey("Load custom config", func() {
				assert.Nil(t, LoadConfig("./custom.yaml"))

				c.Convey("With global config", func() {
					c, ok := clientConfigMap["httpA"]
					assert.True(t, ok)
					assert.Equal(t, "httpA", c.Name)
					assert.Equal(t, "https://httpbin.org", c.DSN)
					assert.Equal(t, int64(3000), c.Timeout)

					userAgent, ok := c.Header["User-Agent"]
					assert.True(t, ok)
					assert.Equal(t, "agent_global", userAgent)

					contentType, ok := c.Header["Content-Type"]
					assert.True(t, ok)
					assert.Equal(t, "content_type_global", contentType)

					transport := c.TransportCfg
					assert.NotNil(t, transport)
					assert.Equal(t, 100, transport.MaxIdleConns)
					assert.Equal(t, 200, transport.MaxIdleConnsPerHost)
					assert.Equal(t, 300, transport.MaxConnsPerHost)
					assert.Equal(t, 1*time.Second, transport.IdleConnTimeout)
					assert.Equal(t, 2*time.Second, transport.TLSHandshakeTimeout)
					assert.Equal(t, 3*time.Second, transport.ExpectContinueTimeout)
					assert.Equal(t, 4*time.Second, transport.ResponseHeaderTimeout)
					assert.NotNil(t, 5*time.Second, transport.Dial.Timeout)
					assert.NotNil(t, 6*time.Second, transport.Dial.KeepAlive)
				})

				c.Convey("With custom config", func() {
					c, ok := clientConfigMap["httpB"]
					assert.True(t, ok)
					assert.Equal(t, "httpB", c.Name)
					assert.Equal(t, "https://httpbin1.org,https://httpbin2.org", c.DSN)
					assert.Equal(t, int64(1000), c.Timeout)

					userAgent, ok := c.Header["User-Agent"]
					assert.True(t, ok)
					assert.Equal(t, "agent_custom", userAgent)

					contentType, ok := c.Header["Content-Type"]
					assert.True(t, ok)
					assert.Equal(t, "content_type_custom", contentType)

					transport := c.TransportCfg
					assert.NotNil(t, transport)
					assert.Equal(t, 300, transport.MaxIdleConns)
					assert.Equal(t, 200, transport.MaxIdleConnsPerHost)
					assert.Equal(t, 100, transport.MaxConnsPerHost)
					assert.Equal(t, 6*time.Second, transport.IdleConnTimeout)
					assert.Equal(t, 5*time.Second, transport.TLSHandshakeTimeout)
					assert.Equal(t, 4*time.Second, transport.ExpectContinueTimeout)
					assert.Equal(t, 3*time.Second, transport.ResponseHeaderTimeout)
					assert.NotNil(t, 2*time.Second, transport.Dial.Timeout)
					assert.NotNil(t, 1*time.Second, transport.Dial.KeepAlive)
				})
			})
		})
	})
}
