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

type paginatedPostsResponse struct {
	Data []models.Post `json:"data"`
	Meta MetaData      `json:"meta"`
}

func TestPostHandler_Create_OK(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)

	body := map[string]interface{}{
		"title":     "My first post",
		"content":   "This is long enough content for validation",
		"author":    "Ivan",
		"published": false,
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusCreated, w.Code)

	var got models.Post
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, body["title"], got.Title)
	assert.Equal(t, body["content"], got.Content)
	assert.Equal(t, body["author"], got.Author)
	assert.False(t, got.Published)
	assert.NotZero(t, got.ID)
}

func TestPostHandler_Create_ValidationError(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)

	body := map[string]interface{}{
		"title":   "ab",
		"content": "This is long enough content for validation",
		"author":  "Ivan",
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostHandler_Create_MalformedJSON(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader([]byte("{bad json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostHandler_List_PaginationAndFilter(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)

	p1 := seedPost(t, db, true)
	p2 := seedPost(t, db, true)
	seedPost(t, db, false)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts?page=1&limit=10&published=true", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusOK, w.Code)

	var resp paginatedPostsResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, int64(2), resp.Meta.Total)
	assert.Len(t, resp.Data, 2)
	assert.True(t, resp.Data[0].Published)
	assert.True(t, resp.Data[1].Published)
	assert.ElementsMatch(t, []uint{p1.ID, p2.ID}, []uint{resp.Data[0].ID, resp.Data[1].ID})
}

func TestPostHandler_Get_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, false)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/posts/%d", post.ID+999), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostHandler_Get_InvalidID(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/abc", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostHandler_Update_OK(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, false)

	body := map[string]interface{}{
		"title":     "Updated title",
		"content":   "Updated content with enough length",
		"author":    "Updated author",
		"published": true,
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/posts/%d", post.ID), bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusOK, w.Code)
	var got models.Post
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "Updated title", got.Title)
	assert.Equal(t, "Updated content with enough length", got.Content)
	assert.Equal(t, "Updated author", got.Author)
	assert.True(t, got.Published)
}

func TestPostHandler_Update_MalformedJSON(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, false)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/posts/%d", post.ID), bytes.NewReader([]byte("{bad json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostHandler_Publish_OK(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, false)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/posts/%d/publish", post.ID), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusOK, w.Code)
	var got models.Post
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.True(t, got.Published)
}

func TestPostHandler_Delete_CascadeComments(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	router := setupTestRouter(db)
	post := seedPost(t, db, false)
	comment := seedComment(t, db, post.ID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/posts/%d", post.ID), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusNoContent, w.Code)

	var deletedPost models.Post
	var deletedComment models.Comment
	assert.ErrorIs(t, db.First(&deletedPost, post.ID).Error, gorm.ErrRecordNotFound)
	assert.ErrorIs(t, db.First(&deletedComment, comment.ID).Error, gorm.ErrRecordNotFound)
}
