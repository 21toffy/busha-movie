package models

import (
	"time"

	// "github.com/21toffy/busha-movie/internal/database"
	"gorm.io/gorm"
)

type ApiResponse struct {
	Count    int               `json:"count"`
	Next     string            `json:"next"`
	Previous interface{}       `json:"previous"`
	Results  []ApiResultStruct `json:"results"`
}

type Film struct {
	Title        string   `json:"title"`
	EpisodeId    int      `json:"episode_id"`
	OpeningCrawl string   `json:"opening_crawl"`
	ReleaseDate  string   `json:"release_date"`
	Characters   []string `json:"characters"`
	CommentCount int64    `json:"film_count"`
}

type FilmsResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []Film  `json:"results"`
}

type ApiResultStruct struct {
	Name      string   `json:"name"`
	Height    string   `json:"height"`
	Mass      string   `json:"mass"`
	HairColor string   `json:"hair_color"`
	SkinColor string   `json:"skin_color"`
	EyeColor  string   `json:"eye_color"`
	BirthYear string   `json:"birth_year"`
	Gender    string   `json:"gender"`
	Homeworld string   `json:"homeworld"`
	Films     []string `json:"films"`
	Species   []string `json:"species"`
	Vehicles  []string `json:"vehicles"`
	Starships []string `json:"starships"`
	Created   string   `json:"created"`
	Edited    string   `json:"edited"`
	Url       string   `json:"url"`
}

type APIResponseCharacter struct {
	Url     string    `json:"url"`
	Name    string    `json:"name"`
	Gender  string    `json:"gender"`
	Created time.Time `json:"created"`
	URL     string    `json:"url"`
	Height  int       `json:"height"`
	Movies  []string  `json:"movies"`
}

type Character struct {
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	Gender  string    `json:"gender"`
	Created time.Time `json:"created"`
	// URL            string    `json:"url"`
	HeightInCM     int   `json:"height_in_cm"`
	HeightInFeet   int   `json:"height_in_feet"`
	HeightInInches int   `json:"height_in_inches"`
	Movies         []int `json:"movies"`
}

// I want the struct to be in this form

type NewCharacter struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Height  string `json:"height"`
	Gender  string `json:"gender"`
	Films   []int  `json:"films"`
	Created string `json:"created"`
}

type CharactersResponse struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Character `json:"results"`
}

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

type CommentRequest struct {
	Comment string `json:"comment,omitempty" validate:"required,max=500" binding:"required"`
}
