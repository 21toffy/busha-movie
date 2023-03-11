package database

import (
	"github.com/21toffy/busha-movie/internal/models"
	// "github.com/21toffy/busha-movie/internal/utils"
	// "github.com/21toffy/busha-movie/internal/database"
)

func MapCommentListToResponse(commentList []models.Comment) []models.CommentResponse {
	responseList := make([]models.CommentResponse, len(commentList))
	for i, comment := range commentList {
		responseList[i] = models.CommentResponse{
			ID:        comment.ID,
			MovieID:   comment.MovieID,
			IPAddress: comment.IPAddress,
			Comment:   comment.Comment,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
	}
	return responseList
}

func GetCommentsReverseChronological() ([]models.CommentResponse, error) {
	db, err := NewDatabase()
	if err != nil {
		return nil, err
	}
	var comments []models.Comment

	dbErr := db.Order("created_at DESC").Find(&comments).Error
	if dbErr != nil {
		return nil, dbErr
	}
	refac := MapCommentListToResponse(comments)
	return refac, nil
}

func GetFilmComments(id int) ([]models.Comment, error) {
	db, err := NewDatabase()
	if err != nil {
		return nil, err
	}
	var comments []models.Comment
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
	rows := db.Model(&models.Comment{}).Where("movie_id = ?", id).Count(&count)
	if rows.Error != nil {
		return 0, rows.Error
	}
	return count, nil
}
