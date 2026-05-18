package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iruiz/gin-blog-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCommentHandler_ListByPost_OK(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, true)
	c1 := seedComment(t, db, post.ID)
	c2 := seedComment(t, db, post.ID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/posts/%d/comments", post.ID), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusOK, w.Code)
	var comments []models.Comment
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &comments))
	assert.Len(t, comments, 2)
	assert.ElementsMatch(t, []uint{c1.ID, c2.ID}, []uint{comments[0].ID, comments[1].ID})
	assert.Equal(t, post.ID, comments[0].PostID)
	assert.Equal(t, post.ID, comments[1].PostID)
}

func TestCommentHandler_ListByPost_Empty(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, true)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/posts/%d/comments", post.ID), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusOK, w.Code)
	var comments []models.Comment
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &comments))
	assert.Len(t, comments, 0)
}

func TestCommentHandler_ListByPost_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, true)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/posts/%d/comments", post.ID+999), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCommentHandler_ListByPost_InvalidID(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/abc/comments", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommentHandler_CreateForPost_OK(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, true)

	body := map[string]string{
		"author":  "Reader",
		"content": "Nice article",
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/posts/%d/comments", post.ID), bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusCreated, w.Code)
	var got models.Comment
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, post.ID, got.PostID)
	assert.Equal(t, "Reader", got.Author)
	assert.Equal(t, "Nice article", got.Content)
}

func TestCommentHandler_CreateForPost_ValidationError(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, true)

	body := map[string]string{
		"content": "Missing author",
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/posts/%d/comments", post.ID), bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommentHandler_CreateForPost_MalformedJSON(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, true)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/posts/%d/comments", post.ID), bytes.NewReader([]byte("{bad json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommentHandler_CreateForPost_InvalidID(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)

	body := map[string]string{
		"author":  "Reader",
		"content": "Nice article",
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/abc/comments", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommentHandler_Delete_OK(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, true)
	comment := seedComment(t, db, post.ID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/comments/%d", comment.ID), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)

	var deleted models.Comment
	assert.ErrorIs(t, db.First(&deleted, comment.ID).Error, gorm.ErrRecordNotFound)
}

func TestCommentHandler_Delete_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, true)
	comment := seedComment(t, db, post.ID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/comments/%d", comment.ID+999), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCommentHandler_Delete_AfterPriorDeletion(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, true)
	comment := seedComment(t, db, post.ID)

	firstReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/comments/%d", comment.ID), nil)
	firstW := httptest.NewRecorder()
	router.ServeHTTP(firstW, firstReq)
	require.Equal(t, http.StatusNoContent, firstW.Code)

	secondReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/comments/%d", comment.ID), nil)
	secondW := httptest.NewRecorder()

	// Act
	router.ServeHTTP(secondW, secondReq)

	// Assert
	assert.Equal(t, http.StatusNotFound, secondW.Code)
}

func TestCommentHandler_Delete_InvalidID(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/abc", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
