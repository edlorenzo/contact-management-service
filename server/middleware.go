package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		log.Debug().Str("method", c.Request.Method).Int("status", c.Writer.Status()).Str("uri", c.Request.RequestURI).Str("latency", latency.String()).Send()
	}
}
