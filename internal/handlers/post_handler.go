package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iruiz/gin-blog-api/internal/metrics"
	"github.com/iruiz/gin-blog-api/internal/middleware"
	"github.com/iruiz/gin-blog-api/internal/models"
	"gorm.io/gorm"
)

type PostHandler struct {
	db      *gorm.DB
	metrics *metrics.Metrics
}

func NewPostHandler(db *gorm.DB, observability *metrics.Metrics) *PostHandler {
	return &PostHandler{db: db, metrics: observability}
}

type PaginatedResponse struct {
	Data interface{} `json:"data"`
	Meta MetaData    `json:"meta"`
}

type MetaData struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func (h *PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	query := h.db.Model(&models.Post{})
	if pub := c.Query("published"); pub != "" {
		if published, err := strconv.ParseBool(pub); err == nil {
			query = query.Where("published = ?", published)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		middleware.RespondError(c, http.StatusInternalServerError, "failed to count posts")
		return
	}

	var posts []models.Post
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		middleware.RespondError(c, http.StatusInternalServerError, "failed to list posts")
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data: posts,
		Meta: MetaData{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}

func (h *PostHandler) Get(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var post models.Post
	if err := h.db.Preload("Comments").First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			middleware.RespondError(c, http.StatusNotFound, "post not found")
			return
		}
		middleware.RespondError(c, http.StatusInternalServerError, "failed to get post")
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) Create(c *gin.Context) {
	var input models.PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "invalid payload: "+err.Error())
		return
	}

	post := models.Post{
		Title:     input.Title,
		Content:   input.Content,
		Author:    input.Author,
		Published: input.Published,
	}

	if err := h.db.Create(&post).Error; err != nil {
		middleware.RespondError(c, http.StatusInternalServerError, "failed to create post")
		return
	}

	h.metrics.Business.PostsCreatedTotal.Inc()
	if post.Published {
		h.metrics.Business.PostsPublishedTotal.Inc()
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var post models.Post
	if err := h.db.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			middleware.RespondError(c, http.StatusNotFound, "post not found")
			return
		}
		middleware.RespondError(c, http.StatusInternalServerError, "failed to get post")
		return
	}

	var input models.PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "invalid payload: "+err.Error())
		return
	}

	post.Title = input.Title
	post.Content = input.Content
	post.Author = input.Author
	wasPublished := post.Published
	post.Published = input.Published

	if err := h.db.Save(&post).Error; err != nil {
		middleware.RespondError(c, http.StatusInternalServerError, "failed to update post")
		return
	}

	if !wasPublished && post.Published {
		h.metrics.Business.PostsPublishedTotal.Inc()
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) Publish(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var post models.Post
	if err := h.db.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			middleware.RespondError(c, http.StatusNotFound, "post not found")
			return
		}
		middleware.RespondError(c, http.StatusInternalServerError, "failed to get post")
		return
	}

	wasPublished := post.Published
	post.Published = true
	if err := h.db.Save(&post).Error; err != nil {
		middleware.RespondError(c, http.StatusInternalServerError, "failed to publish post")
		return
	}

	if !wasPublished {
		h.metrics.Business.PostsPublishedTotal.Inc()
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var post models.Post
	if err := h.db.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			middleware.RespondError(c, http.StatusNotFound, "post not found")
			return
		}
		middleware.RespondError(c, http.StatusInternalServerError, "failed to get post")
		return
	}

	if err := h.db.Delete(&post).Error; err != nil {
		middleware.RespondError(c, http.StatusInternalServerError, "failed to delete post")
		return
	}

	h.metrics.Business.PostsDeletedTotal.Inc()

	c.Status(http.StatusNoContent)
}

func parseID(raw string) (uint, error) {
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}
