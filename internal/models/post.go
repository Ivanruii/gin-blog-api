package models

import "time"

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:200;not null" json:"title" binding:"required,min=3,max=200"`
	Content   string    `gorm:"type:text;not null" json:"content" binding:"required,min=10"`
	Author    string    `gorm:"size:100;not null" json:"author" binding:"required"`
	Published bool      `gorm:"default:false" json:"published"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
}

type PostInput struct {
	Title     string `json:"title" binding:"required,min=3,max=200"`
	Content   string `json:"content" binding:"required,min=10"`
	Author    string `json:"author" binding:"required"`
	Published bool   `json:"published"`
}
