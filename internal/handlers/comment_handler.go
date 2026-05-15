package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iruiz/gin-blog-api/internal/middleware"
	"github.com/iruiz/gin-blog-api/internal/models"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

func (h *CommentHandler) ListByPost(c *gin.Context) {
	postID, err := parseID(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var post models.Post
	if err := h.db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			middleware.RespondError(c, http.StatusNotFound, "post not found")
			return
		}
		middleware.RespondError(c, http.StatusInternalServerError, "failed to get post")
		return
	}

	var comments []models.Comment
	if err := h.db.Where("post_id = ?", postID).Order("created_at DESC").Find(&comments).Error; err != nil {
		middleware.RespondError(c, http.StatusInternalServerError, "failed to list comments")
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) CreateForPost(c *gin.Context) {
	postID, err := parseID(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var post models.Post
	if err := h.db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			middleware.RespondError(c, http.StatusNotFound, "post not found")
			return
		}
		middleware.RespondError(c, http.StatusInternalServerError, "failed to get post")
		return
	}

	var input models.CommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "invalid payload: "+err.Error())
		return
	}

	comment := models.Comment{
		PostID:  postID,
		Author:  input.Author,
		Content: input.Content,
	}

	if err := h.db.Create(&comment).Error; err != nil {
		middleware.RespondError(c, http.StatusInternalServerError, "failed to create comment")
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var comment models.Comment
	if err := h.db.First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			middleware.RespondError(c, http.StatusNotFound, "comment not found")
			return
		}
		middleware.RespondError(c, http.StatusInternalServerError, "failed to get comment")
		return
	}

	if err := h.db.Delete(&comment).Error; err != nil {
		middleware.RespondError(c, http.StatusInternalServerError, "failed to delete comment")
		return
	}

	c.Status(http.StatusNoContent)
}
