package http

import (
	stdLogger "log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Config are the configurations related to the HTTP Server
type Config struct {
	HTTPPort       string `env:"HTTP_PORT"`
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	HandlerTimeout time.Duration
}

func NewServer(cfg Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      http.TimeoutHandler(router, cfg.HandlerTimeout, ""),
		ErrorLog:     stdLogger.New(os.Stderr, "http: ", stdLogger.LstdFlags),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
