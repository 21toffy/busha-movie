package api

import (
	docs "github.com/21toffy/busha-movie/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title Busha movie Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func SetupRoutes(router *gin.Engine) {
	router.GET("/films/comments", GetCommentsHandler())
	router.GET("/films/:id/character/:gender/:sort/:order", FetchFilmCharacter())
	router.GET("/films/:id/comments", GetFIlmCommentsHandler())
	router.POST("/films/:id/comment/create", SaveCommentHandler())
	router.GET("/films/all", FetchFilmsHandler())
	router.Use(logRequest())
	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}
