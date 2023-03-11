package requests

import (
	"encoding/json"
	"math"

	"github.com/21toffy/busha-movie/internal/cache"
	"github.com/21toffy/busha-movie/internal/customerror"
	"github.com/21toffy/busha-movie/internal/database"

	"github.com/21toffy/busha-movie/internal/models"
	"github.com/21toffy/busha-movie/internal/utils"

	"log"

	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

// Filter characters by gender if filter parameter is present
// Sort characters by the specified field and order
// Calculate and return metadata
func FilterAndSortCharacters(characters []models.Character, filter string, sortBy string, sortOrder string) ([]models.Character, int, float64, float64, float64, error) {
	// Filter characters by gender if filter parameter is present
	if filter != "" {
		filteredCharacters := make([]models.Character, 0)
		for _, char := range characters {
			if char.Gender == filter {
				filteredCharacters = append(filteredCharacters, char)
			}
		}
		characters = filteredCharacters
	}

	// Sort characters by the specified field and order
	switch sortBy {
	case "name":
		if sortOrder == "asc" {
			sort.Slice(characters, func(i, j int) bool { return characters[i].Name < characters[j].Name })
		} else if sortOrder == "desc" {
			sort.Slice(characters, func(i, j int) bool { return characters[i].Name > characters[j].Name })
		}
	case "gender":
		if sortOrder == "asc" {
			sort.Slice(characters, func(i, j int) bool { return characters[i].Gender < characters[j].Gender })
		} else if sortOrder == "desc" {
			sort.Slice(characters, func(i, j int) bool { return characters[i].Gender > characters[j].Gender })
		}
	case "height":
		if sortOrder == "asc" {
			sort.Slice(characters, func(i, j int) bool { return characters[i].HeightInCM < characters[j].HeightInCM })
		} else if sortOrder == "desc" {
			sort.Slice(characters, func(i, j int) bool { return characters[i].HeightInCM > characters[j].HeightInCM })
		}
	default:
		return nil, 0, 0, 0, 0, customerror.InvalidSortParam
	}

	// Calculate metadata
	totalHeightInCm := 0
	totalHeightInFeet := 0
	totalHeightInInches := 0
	for _, char := range characters {
		totalHeightInCm += int(char.HeightInCM)
		totalHeightInFeet += int(utils.CmToFeet(char.HeightInCM))
		totalHeightInInches += utils.CmToInches(char.HeightInCM)
	}
	totalHeightInFeet += totalHeightInInches / 12.0
	totalHeightInInches = int(math.Mod(float64(totalHeightInInches), 12.0))

	return characters, len(characters), float64(totalHeightInCm), float64(totalHeightInFeet), float64(totalHeightInInches), nil
}

// load characters to redis or
// fetch from swapi and load to redis
func FetchCharacters(movieID string) error {
	films := []models.Film{}

	redisInstance := cache.NewRedisCache()
	err := redisInstance.Get("films", &films)
	if err != nil {
		return err
	}
	var targetFilm models.Film
	for _, film := range films {
		if strconv.Itoa(film.EpisodeId) == movieID {
			targetFilm = film
			break
		}
	}
	if targetFilm.EpisodeId == 0 {
		return customerror.NotFoundInCache
	}
	return nil
}

// update comment count in redis if a new comment was added
func UpdateCount(filmId string) error {
	films := []models.Film{}
	redisInstance := cache.NewRedisCache()
	err := redisInstance.Get("films", &films)

	for i, filmObj := range films {
		urlID, err := strconv.Atoi(filmId)
		if err != nil {
			return customerror.FailedIdConversion
		}
		if filmObj.EpisodeId == urlID {
			films[i].CommentCount = filmObj.CommentCount + 1
			break
		}
	}
	err = redisInstance.Set("films", films, time.Hour)
	if err != nil {
		return customerror.FailedCacheSetting
	}
	return nil
}

func FetchFilmFromRedis(id string) (int, error) {
	films := []models.Film{}
	redisInstance := cache.NewRedisCache()
	err := redisInstance.Get("films", &films)
	if err != nil {
		if err != nil {
			return 0, err
		}
		return 0, err
	}
	var foundFilm *models.Film
	for _, filmObj := range films {
		urlID, err := strconv.Atoi(id)
		if err != nil {
			return 0, customerror.FailedIdConversion
		}
		if filmObj.EpisodeId == urlID {
			foundFilm = &filmObj
			break
		}
	}

	if foundFilm == nil {
		return 0, customerror.NotFoundInCache
	}
	return foundFilm.EpisodeId, nil
}

// load films to redis or fetch from swapi and
// load to redis also return just error if ther is one
func LoadFilms() error {
	films := []models.Film{}
	redisInstance := cache.NewRedisCache()

	err := redisInstance.Get("films", &films)

	if err == nil {
		return nil
	}

	if err != nil {
		return nil
	}

	res, err := http.Get("https://swapi.dev/api/films/")
	if err != nil {
		return customerror.FailedFetch
	}
	defer res.Body.Close()

	var filmsResponse models.FilmsResponse
	if err := json.NewDecoder(res.Body).Decode(&filmsResponse); err != nil {
		return customerror.DecodeError
	}

	for _, film := range filmsResponse.Results {
		f := models.Film{
			Title:        film.Title,
			EpisodeId:    film.EpisodeId,
			OpeningCrawl: film.OpeningCrawl,
			ReleaseDate:  film.ReleaseDate,
			Characters:   film.Characters,
		}
		films = append(films, f)
	}

	err = redisInstance.Set("films", films, time.Hour)
	if err != nil {
		return customerror.FailedCacheSetting
	}
	return nil

}

// fetch from redis or load characters from swapi
func LoadCharacters() ([]models.Character, error) {
	redisInstance := cache.NewRedisCache()
	var characters []models.Character

	err := redisInstance.Get("characters", &characters)
	if err == nil && len(characters) > 0 {
		return characters, nil
	}

	apiUrl := "https://swapi.dev/api/people/"
	characterChan := make(chan models.Character)
	doneChan := make(chan bool)
	errorChan := make(chan error)
	var wg sync.WaitGroup

	go func() {
		for char := range characterChan {
			characters = append(characters, char)
		}
		doneChan <- true
	}()

	for {
		res, err := http.Get(apiUrl)
		if err != nil {
			return nil, customerror.FailedFetch
		}

		var response models.ApiResponse
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			return nil, customerror.DecodeError
		}

		for _, char := range response.Results {
			wg.Add(1)
			go func(char models.ApiResultStruct) {
				defer wg.Done()

				movieIds := []int{}
				for _, singleID := range char.Films {
					movieId, err := utils.GetSingleNumberFromUrl(singleID) //atm
					if err != nil {
						log.Println("Failed get number from URL:", singleID)
						// errorChan <- err
						return
					}
					movieIds = append(movieIds, movieId)
				}

				height, err := strconv.Atoi(char.Height)
				if err != nil {
					// errorChan <- customerror.FailedIdConversion
					log.Println("Failed convert this height to an integer:", char.Height)

					return
				}

				timeGotten, timeErr := utils.StringToTime(char.Created)
				if timeErr != nil {
					log.Println("Failed convert time to a time object:", char.Created)
					// errorChan <- customerror.FailedTimeConversion
					return
				}

				charID, err := utils.GetSingleNumberFromUrl(char.Url)

				if err != nil {
					log.Println("Failed get character ID from character URL:", char.Url)
					// errorChan <- err
					return
				}

				c := models.Character{
					Id:             charID,
					Name:           char.Name,
					Gender:         char.Gender,
					Created:        timeGotten,
					HeightInCM:     height,
					HeightInFeet:   utils.CmToFeet(height),
					HeightInInches: utils.CmToInches(height),
					Movies:         movieIds,
				}

				characterChan <- c
			}(char)
		}

		apiUrl = response.Next
		if apiUrl == "" {
			break
		}

	}

	go func() {
		wg.Wait()
		close(characterChan)
	}()

	select {
	case <-doneChan:
		err = redisInstance.Set("characters", characters, time.Hour)
		if err != nil {
			return nil, customerror.FailedCacheSetting
		}
		return characters, nil
	case err := <-errorChan:
		return nil, err
	}
}

// fetch movies from redis or swapui endpoint and return movie
func FetchFilms() ([]models.Film, error) {
	films := []models.Film{}
	newMovies := []models.Film{}

	redisInstance := cache.NewRedisCache()
	err := redisInstance.Get("films", &films)
	if err == nil {
		return films, nil
	}
	if err != nil && err != customerror.ErrCacheMiss {
		return nil, err
	}
	res, err := http.Get("https://swapi.dev/api/films/")
	if err != nil {
		return nil, customerror.FailedFetch
	}
	defer res.Body.Close()

	var filmsResponse models.FilmsResponse
	if err := json.NewDecoder(res.Body).Decode(&filmsResponse); err != nil {
		return nil, customerror.DecodeError
	}

	for _, film := range filmsResponse.Results {
		commentCount, err := database.GetCommentsByEpisodeIdCount(film.EpisodeId)
		if err != nil {
			return nil, err
		}

		newMovie := models.Film{
			Title:        film.Title,
			EpisodeId:    film.EpisodeId,
			OpeningCrawl: film.OpeningCrawl,
			ReleaseDate:  film.ReleaseDate,
			Characters:   film.Characters,
			CommentCount: commentCount,
		}
		newMovies = append(newMovies, newMovie)
	}

	err = redisInstance.Set("films", newMovies, time.Hour)
	if err != nil {
		return nil, customerror.FailedCacheSetting
	}

	sort.Slice(newMovies, func(i, j int) bool {
		return newMovies[i].ReleaseDate < newMovies[j].ReleaseDate
	})

	return newMovies, nil
}
