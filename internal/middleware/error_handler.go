package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recover() != nil {
				RespondError(c, http.StatusInternalServerError, "internal server error")
			}
		}()

		c.Next()
	}
}
