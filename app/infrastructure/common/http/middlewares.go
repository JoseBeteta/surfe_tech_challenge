package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var excludeMetrics = map[string]struct{}{
	homePath: {},
}

// MetricsAgent is the metrics agent interface
type MetricsAgent interface {
	Timing(name string, duration time.Duration, tags []string)
}

// Consume ensures client is sending data in specific content type / version
func Consume(version Version) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength == 0 {
			c.Next()
			return
		}

		if c.ContentType() == version {
			c.Next()
			return
		}

		errorMessage := fmt.Sprintf(
			"Unsupported content type %q; expected client to send %q",
			c.ContentType(),
			version,
		)

		errorResponseJson(c, http.StatusUnsupportedMediaType, errorMessage)
	}
}

// Produce ensures client is able to accept a response in specific content type / version
func Produce(version Version) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.NegotiateFormat(version) == version {
			c.Header("Content-Type", fmt.Sprintf("%s; charset=utf-8", version))
			c.Next()
			return
		}

		errorMessage := fmt.Sprintf(
			"Unable to produce response of type %q; expected client to accept %q",
			c.Accepted,
			version,
		)

		errorResponseJson(c, http.StatusNotAcceptable, errorMessage)
	}
}

// MetricsMiddleware sends api metrics
func MetricsMiddleware(m MetricsAgent) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := excludeMetrics[c.FullPath()]; !ok {
			t := time.Now()
			defer func() {
				m.Timing("http_request", time.Since(t), []string{
					fmt.Sprintf("full_path:%v", c.FullPath()),
					fmt.Sprintf("path_params:%v", c.Params),
					fmt.Sprintf("method:%v", c.Request.Method),
					fmt.Sprintf("url_path:%v", c.Request.URL.Path),
					fmt.Sprintf("request_content_length:%v", c.Request.ContentLength),
					fmt.Sprintf("response_status_code:%v", c.Writer.Status()),
					fmt.Sprintf("response_size:%v", c.Writer.Size()),
				})
			}()
		}

		c.Next()
	}
}
