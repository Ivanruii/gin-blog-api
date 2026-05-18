package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/iruiz/gin-blog-api/internal/handlers"
	"github.com/iruiz/gin-blog-api/internal/metrics"
	"github.com/iruiz/gin-blog-api/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, observability *metrics.Metrics) *gin.Engine {
	r := gin.New()
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.Logger())
	r.Use(middleware.Metrics(observability))
	r.Use(gin.Recovery())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	postHandler := handlers.NewPostHandler(db, observability)
	commentHandler := handlers.NewCommentHandler(db, observability)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", handlers.HealthHandler)
		v1.GET("/posts", postHandler.List)
		v1.GET("/posts/:id", postHandler.Get)
		v1.POST("/posts", postHandler.Create)
		v1.PUT("/posts/:id", postHandler.Update)
		v1.PATCH("/posts/:id/publish", postHandler.Publish)
		v1.DELETE("/posts/:id", postHandler.Delete)

		v1.GET("/posts/:id/comments", commentHandler.ListByPost)
		v1.POST("/posts/:id/comments", commentHandler.CreateForPost)
		v1.DELETE("/comments/:id", commentHandler.Delete)
	}

	return r
}
