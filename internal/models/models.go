package models

import (
	"time"

	// "github.com/21toffy/busha-movie/internal/database"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	MovieID   int    `json:"movie_id"`
	IPAddress string `json:"ip_address"`
	Comment   string `json:"comment" validate:"required,max=500"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CommentsResponse struct {
	Status   int       `json:"status"`
	Comments []Comment `json:"comments"`
}

type CommentResponse struct {
	ID        uint       `json:"id"`
	MovieID   int        `json:"movie_id"`
	IPAddress string     `json:"ip_address"`
	Comment   string     `json:"comment" validate:"required,max=500"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
