package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/21toffy/busha-movie/internal/customerror"
	"github.com/21toffy/busha-movie/internal/database"
	"github.com/21toffy/busha-movie/internal/models"
	"github.com/21toffy/busha-movie/internal/requests"

	// "github.com/21toffy/busha-movie/internal/models"
	// import "gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// @Summary Fetch film characters
// @Description Fetches all characters for a given film
// @Tags character
// @Param id path string true "Film ID"
// @Param sort path string true "sort by (name, gender, height)"
// @Param order path string true "sort by (asc, desc)"
// @Param gender path string true "The gender to filter by (male, female)"
// @Produce json
// @Success 200 {array} requests.Character
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /films/{id}/character/{gender}/{sort}/{order} [get]
func FetchFilmCharacter() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := requests.FetchCharacters(id)
		if err != nil {
			if err == customerror.NotFoundInCache {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  http.StatusNotFound,
					"message": "Film not found",
				})
				return
			} else if err == customerror.ErrCacheMiss {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusBadRequest,
					"message": "Film not found in cache, cache is mostlikely empty",
				})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": err.Error(),
				})
				return
			}
		}

		charLoaded, err := requests.LoadCharacters()
		if err != nil {
			if err == customerror.NotFoundInCache {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  http.StatusNotFound,
					"message": "Film not found",
				})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": err.Error(),
				})
				return
			}

		}

		// filterBy := strings.TrimSpace(c.Query("gender"))
		// sortBy := strings.TrimSpace(c.Query("sort"))
		// sortOrder := strings.TrimSpace(c.Query("order"))
		filterBy := c.Param("gender")
		sortBy := c.Param("sort")
		sortOrder := c.Param("order")

		characters, numberOfChar, totalHeight, totalHeightFeet, totalHeightInches, err := requests.FilterAndSortCharacters(charLoaded, filterBy, sortBy, sortOrder)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data": gin.H{
				"characters": characters,
				"metadata": gin.H{
					"totalMatches":    numberOfChar,
					"totalHeightInCm": totalHeight,
					"feet":            totalHeightFeet,
					"inches":          totalHeightInches,
				},
			},
		})
		return
	}
}

// @Summary Get comments for a film
// @Description Returns a list of comments for a given film ID
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Film ID"
// @Success 200 {object} []models.CommentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /films/{id}/comments [get]
func GetFIlmCommentsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		urlID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": customerror.FailedIdConversion,
			})
			return
		}

		comments, err := database.GetFilmComments(urlID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to fetch comments",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"comments": comments,
		})

	}
}

// GetCommentsHandler returns a handler function that retrieves all comments
// in reverse chronological order from the database and returns them as JSON.
// @Summary Get all comments
// @Tags comments
// @Description Retrieve all comments in reverse chronological order
// @ID get-comments
// @Produce json
// @Success 200 {object} []models.CommentResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /films/comments [get]
func GetCommentsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		comments, err := database.GetCommentsReverseChronological()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to fetch comments",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"comments": comments,
		})
	}
}

// SaveCommentHandler saves a comment for a film
// @Summary Save a comment for a film
// @Description Saves a comment for a film and updates the comment count in cache
// @Tags comments
// @Accept  json
// @Produce  json
// @Param id path int true "Film ID"
// @Param comment body CommentRequest true "Comment payload"
// @Success 200 {object} models.CommentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /films/{id}/comment/create [post]
func SaveCommentHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var comment models.CommentRequest
		var dbComment models.Comment

		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filmFound, err := requests.FetchFilmFromRedis(id)
		if err != nil {
			if err == customerror.NotFoundInCache {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  http.StatusNotFound,
					"message": "Film not found",
				})
				return
			} else if err == customerror.ErrCacheMiss {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusBadRequest,
					"message": "Film not found in cache, cache is mostlikely empty",
				})
				return
			} else {

				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": err.Error(),
				})
				return
			}
		}
		db, err := database.NewDatabase()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": fmt.Errorf("failed to get database instance: %s", err),
			})
			return
		}
		dbComment.MovieID = filmFound
		dbComment.IPAddress = c.Request.RemoteAddr
		dbComment.CreatedAt = time.Now().UTC()
		dbComment.Comment = comment.Comment
		result := db.Create(&dbComment)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save comment",
			})
			return
		}
		redisUpdateErr := requests.UpdateCount(id)
		if redisUpdateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save comment to cache",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   dbComment,
		})
		return
	}
}

// FetchFilmsHandler returns a handler function that retrieves all films
// @Tags films
// @Summary Get all films
// @Description Retrieve all films
// @Produce json
// @Success 200 {object} []requests.Film
// @Failure 500 {object} models.ErrorResponse
// @Router /films/all [get]
func FetchFilmsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		films, err := requests.FetchFilms()
		if err != nil {
			if err == customerror.ErrCacheMiss {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": err.Error(),
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   films,
		})
		return
	}
}

func Home() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.File("internal/templates/home.html")
	}
}
