package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/21toffy/busha-movie/internal/api"
	"github.com/21toffy/busha-movie/internal/cache"
	"github.com/21toffy/busha-movie/internal/config"
	"github.com/21toffy/busha-movie/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Error initializing config: %s", err.Error())
	}

	db, err := database.NewDatabase()
	if err != nil {
		fmt.Errorf("Failed to create a new database instance: %v", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Errorf("Failed to get database instance: %v", err)
		return
	}
	defer sqlDB.Close()

	// Initialize cache
	cache.InitRedisCache()

	r := gin.Default()
	api.SetupRoutes(r)

	// start the server
	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.GetPort()), r); err != nil {
		fmt.Errorf("Failed to start the server: %v", err)
	}
}
