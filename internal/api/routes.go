package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/ping", handlePing)
	router.GET("/films", FetchFilmsHandler())
	router.POST("/film-comment/:id", SaveCommentHandler())
	router.GET("/film-comment/", GetCommentsHandler())
	router.GET("/film/:id/comment/", GetFIlmCommentsHandler())
	router.GET("/film/:id/character/", FetchFilmCharacter())
	// router.GET("/films/:id/characters", handlers.FetchFilmCharacters())

	router.Use(logRequest())
}
