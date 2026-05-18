package middleware

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func RespondError(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, ErrorResponse{Error: msg, Code: status})
}
