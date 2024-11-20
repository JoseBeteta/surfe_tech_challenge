package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Version HTTP API version
type Version = string

// V1 HTTP API version 1
const V1 Version = "application/vnd.surfe.v1+json"

// ErrEmptyBody error when body is empty
var ErrEmptyBody = errors.New("missing request body")

func okResponseJson(c *gin.Context, obj any) {
	c.JSON(http.StatusOK, obj)
}

func errorResponseJson(c *gin.Context, statusCode int, err string) {
	c.AbortWithStatusJSON(statusCode, gin.H{"error": err})
}

func handleRecovery(c *gin.Context, err any) {
	if msg, ok := err.(string); ok {
		errorResponseJson(c, http.StatusInternalServerError, msg)
		return
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}
