package database

// import (
// 	"github.com/21toffy/busha-movie/internal/models"
// 	// "github.com/21toffy/busha-movie/internal/utils"
// )

// func GetCommentsReverseChronological() ([]models.CommentResponse, error) {
// 	db, err := NewDatabase()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var comments []models.Comment

// 	dbErr := db.Order("created_at DESC").Find(&comments).Error
// 	if dbErr != nil {
// 		return nil, dbErr
// 	}
// 	refac := models.MapCommentListToResponse(comments)
// 	return refac, nil
// }

// func GetFilmComments(id int) ([]models.Comment, error) {
// 	db, err := NewDatabase()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var comments []models.Comment
// 	dbErr := db.Order("created_at").Where("movie_id = ?", id).Find(&comments).Error
// 	if dbErr != nil {
// 		return nil, dbErr
// 	}
// 	return comments, nil
// }

// var count int64

// func GetCommentsByEpisodeIdCount(id int) (int64, error) {
// 	db, err := NewDatabase()
// 	if err != nil {
// 		return 0, err
// 	}

// 	var count int64
// 	rows := db.Model(&models.Comment{}).Where("movie_id = ?", id).Count(&count)
// 	if rows.Error != nil {
// 		return 0, rows.Error
// 	}
// 	return count, nil
// }
