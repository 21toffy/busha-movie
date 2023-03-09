package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/21toffy/busha-movie/internal/customerror"
	"github.com/21toffy/busha-movie/internal/database"
	"github.com/21toffy/busha-movie/internal/requests"
	"github.com/gin-gonic/gin"
)

func handlePing(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

type CommentRequest struct {
	Comment string `json:"comment, omitempty" validate:"required,max=500" binding:"required"`
}

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

		filterBy := strings.TrimSpace(c.Query("gender"))
		sortBy := strings.TrimSpace(c.Query("sort"))
		sortOrder := strings.TrimSpace(c.Query("order"))

		characters, numberOfChar, totalHeight, err := requests.FilterAndSortCharacters(charLoaded, filterBy, sortBy, sortOrder)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		// totalHeightInCm, totalHeightInFeet, totalHeightInInches, totalMatches, err := requests.GetMetadata(characters)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"status":  http.StatusInternalServerError,
		// 		"message": err.Error(),
		// 	})
		// 	return
		// }

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data": gin.H{
				"characters": characters,
				"metadata": gin.H{
					"totalMatches":    numberOfChar,
					"totalHeightInCm": totalHeight,
					// "totalHeightInFeet": gin.H{
					// 	"feet":   totalHeightInFeet,
					// 	"inches": totalHeightInInches,
					// },
				},
			},
		})

		// c.JSON(http.StatusOK, gin.H{
		// 	"status": http.StatusOK,
		// 	"data":   charLoaded,
		// })
		return
	}
}

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

func SaveCommentHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var comment CommentRequest
		var dbComment database.Comment

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
