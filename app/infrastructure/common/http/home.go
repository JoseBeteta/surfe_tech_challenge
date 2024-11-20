package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const homePath = "/"

func RegisterHomeHandler(router *gin.Engine) {
	router.GET(homePath, newHomeHandler())
}

func newHomeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Congratulations, you have reached DT one!")
	}
}
