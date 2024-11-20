package http_test

import (
	"github.com/JoseBeteta/surfe/app/infrastructure/common/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)

	config := http.Config{
		HTTPPort:       ":1234",
		ReadTimeout:    1 * time.Second,
		HandlerTimeout: 2 * time.Second,
		WriteTimeout:   3 * time.Second,
		IdleTimeout:    15 * time.Second,
	}
	server := http.NewServer(config, engine)

	assert.Equal(t, ":1234", server.Addr)
	assert.Equal(t, 1*time.Second, server.ReadTimeout)
	assert.Equal(t, 3*time.Second, server.WriteTimeout)
	assert.Equal(t, 15*time.Second, server.IdleTimeout)
}
