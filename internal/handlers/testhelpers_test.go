package handlers

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iruiz/gin-blog-api/internal/metrics"
	"github.com/iruiz/gin-blog-api/internal/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// Each test gets a fresh in-memory SQLite database, so state cannot leak
	// across tests and execution order does not affect results.
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	require.NoError(t, db.Exec("PRAGMA foreign_keys = ON").Error)
	require.NoError(t, db.AutoMigrate(&models.Post{}, &models.Comment{}))

	return db
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(gin.Recovery())

	observability := metrics.NewMetrics(prometheus.NewRegistry())
	postHandler := NewPostHandler(db, observability)
	commentHandler := NewCommentHandler(db, observability)

	v1 := r.Group("/api/v1")
	{
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

func seedPost(t *testing.T, db *gorm.DB, published bool) models.Post {
	t.Helper()

	post := models.Post{
		Title:     "Seed post",
		Content:   "Seed post content with enough length",
		Author:    "Seeder",
		Published: published,
	}
	require.NoError(t, db.Create(&post).Error)
	return post
}

func seedComment(t *testing.T, db *gorm.DB, postID uint) models.Comment {
	t.Helper()

	comment := models.Comment{
		PostID:  postID,
		Author:  "Reader",
		Content: "Great post!",
	}
	require.NoError(t, db.Create(&comment).Error)
	return comment
}
