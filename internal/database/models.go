package database

import (
	// "fmt"
	// "github.com/21toffy/busha-movie/internal/database"
	// "time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	MovieID   int    `json:"movie_id"`
	IPAddress string `json:"ip_address"`
	Comment   string `json:"comment" validate:"required,max=500"`
}

func GetCommentsReverseChronological() ([]Comment, error) {
	db, err := NewDatabase()
	if err != nil {
		return nil, err
	}
	var comments []Comment
	dbErr := db.Order("created_at DESC").Find(&comments).Error
	if dbErr != nil {
		return nil, dbErr
	}
	return comments, nil
}

func GetFilmComments(id int) ([]Comment, error) {
	db, err := NewDatabase()
	if err != nil {
		return nil, err
	}
	var comments []Comment
	dbErr := db.Order("created_at").Where("movie_id = ?", id).Find(&comments).Error
	if dbErr != nil {
		return nil, dbErr
	}
	return comments, nil
}

var count int64

func GetCommentsByEpisodeIdCount(id int) (int64, error) {
	db, err := NewDatabase()
	if err != nil {
		return 0, err
	}

	var count int64
	rows := db.Model(&Comment{}).Where("movie_id = ?", id).Count(&count)
	if rows.Error != nil {
		return 0, rows.Error
	}
	return count, nil

	// var comments []Comment
	// commentCount, dbErr := db.Order("created_at").Where("movie_id = ?", id).Find(&comments).Count(&count)
	// if dbErr != nil {
	// 	return 0, dbErr
	// }
	// return commentCount, nil
}
