package api

import (
	"github.com/21toffy/busha-movie/internal/logger"
	"github.com/gin-gonic/gin"
)

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Incoming request")
		c.Next()
	}
}
