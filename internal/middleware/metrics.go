package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iruiz/gin-blog-api/internal/metrics"
)

func Metrics(observability *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		observability.HTTP.RequestsInFlight.Inc()
		defer observability.HTTP.RequestsInFlight.Dec()

		start := time.Now()
		c.Next()

		duration := time.Since(start).Seconds()
		path := observability.RouteLabel(c.FullPath())

		status := c.Writer.Status()
		statusStr := strconv.Itoa(status)

		observability.HTTP.RequestsTotal.WithLabelValues(c.Request.Method, path, statusStr).Inc()
		observability.HTTP.RequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)

		if errorType := observability.ErrorTypeForStatus(status); errorType != "" {
			observability.HTTP.ErrorsTotal.WithLabelValues(errorType).Inc()
		}
	}
}
