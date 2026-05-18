package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	Author    string    `gorm:"size:100;not null" json:"author" binding:"required"`
	Content   string    `gorm:"size:500;not null" json:"content" binding:"required,min=1,max=500"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentInput struct {
	Author  string `json:"author" binding:"required"`
	Content string `json:"content" binding:"required,min=1,max=500"`
}
