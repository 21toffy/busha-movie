package config

import (
	// "fmt"
	"github.com/spf13/viper"
)

func InitConfig() error {
	// Set the name of the config file (without extension)
	viper.SetConfigName("config")
	// Set the paths to look for the config file
	viper.AddConfigPath(".")
	viper.AddConfigPath("./internal/config")
	// Set the file formats to search for
	viper.SetConfigType("toml")

	// Read in the config file
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Set default values for any missing configuration options
	viper.SetDefault("port", "8081")
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "5432")
	viper.SetDefault("db.user", "postgres")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.database", "movie_db")
	viper.SetDefault("redis.addr", "movieapi_redis:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", "0")

	return nil
}

func GetPort() string {
	return viper.GetString("port")
}
