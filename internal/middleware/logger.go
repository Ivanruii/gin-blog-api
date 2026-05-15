package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iruiz/gin-blog-api/internal/applog"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		applog.Logger.Printf("[%s] %s %d %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
		)
	}
}
