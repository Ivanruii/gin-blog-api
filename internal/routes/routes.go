package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/iruiz/gin-blog-api/internal/handlers"
	"gorm.io/gorm"
)

func Setup(_ *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", handlers.HealthHandler)
	}

	return r
}
