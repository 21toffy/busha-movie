package utils

import (
	// "fmt"
	"encoding/json"
	"errors"

	"github.com/21toffy/busha-movie/internal/customerror"

	// "github.com/21toffy/busha-movie/internal/requests"

	"strconv"
	"strings"
	"time"
)

func StringToTime(stringedTime string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000000Z"
	t, err := time.Parse(layout, stringedTime)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func GetSingleNumberFromUrl(url string) (int, error) {
	if url == "" {
		return 0, customerror.EmptyString
	}

	numStr := strings.TrimRight(url, "/")
	numStr = numStr[strings.LastIndex(numStr, "/")+1:]

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, customerror.FailedIdConversion
	}

	return num, nil
}

func GetMovieNumber(arr []string) (int, error) {
	if len(arr) == 0 {
		return 0, errors.New("array is empty")
	}

	lastUrl := arr[len(arr)-1]
	numStr := strings.TrimRight(lastUrl, "/")
	numStr = numStr[strings.LastIndex(numStr, "/")+1:]

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}

	return num, nil
}

type Person struct {
	Id int `json:"id"`

	Name           string    `json:"name"`
	Gender         string    `json:"gender"`
	Created        time.Time `json:"created"`
	URL            string    `json:"url"`
	HeightInCM     int       `json:"height_in_cm"`
	HeightInFeet   int       `json:"height_in_feet"`
	HeightInInches int       `json:"height_in_inches"`
	Movies         []int     `json:"movies"`
}

// func NewCharacterFromApiResponse(apiResp requests.ApiResponse) (*Person, error) {
// 	// Extract the ID number from the URL.
// 	id, err := GetNumberFromUrl(apiResp.Url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	name := apiResp.Name
// 	gender := apiResp.Gender
// 	createdStr := apiResp.Created
// 	created, err := time.Parse(time.RFC3339, createdStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	heightCmStr := apiResp.Height
// 	heightCm, err := strconv.Atoi(heightCmStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	heightFeet := CmToFeet(heightCm)
// 	heightInches := CmToInches(heightCm)

// 	// Extract the movie URLs and extract the ID numbers from them.
// 	movieIds := []int{}
// 	for _, movieUrl := range apiResp.Films {
// 		movieId, err := GetNumberFromUrl(movieUrl)
// 		if err != nil {
// 			return nil, err
// 		}
// 		movieIds = append(movieIds, movieId)
// 	}

// 	// Create a new Person object with the extracted values.
// 	character := &Person{
// 		Id:             id,
// 		Name:           name,
// 		Gender:         gender,
// 		Created:        created,
// 		HeightInCM:     heightCm,
// 		HeightInFeet:   heightFeet,
// 		HeightInInches: heightInches,
// 		Movies:         movieIds,
// 	}

// 	return character, nil
// }

func ConvertPerson(data string) (*Person, error) {
	p := new(Person)

	// Unmarshal the json into a map[string]interface{}
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil, err
	}

	urlSlice, ok := m["url"].([]string)
	if !ok {
		return nil, customerror.NotFoundInCache
	}

	characterID, err := GetMovieNumber(urlSlice)

	if err != nil {
		return nil, err
	}

	p.Id = characterID
	p.Name = m["name"].(string)
	p.Gender = m["gender"].(string)
	p.URL = m["url"].(string)
	createdStr := m["created"].(string)
	created, err := time.Parse(time.RFC3339, createdStr)
	if err != nil {
		return nil, err
	}
	p.Created = created

	// Convert height from cm to feet and inches
	heightStr := m["height"].(string)
	heightInCM, err := strconv.Atoi(heightStr)
	if err != nil {
		return nil, err
	}
	p.HeightInCM = heightInCM
	p.HeightInFeet = CmToFeet(heightInCM)
	p.HeightInInches = CmToInches(heightInCM)
	films := m["films"].([]interface{})
	movies := make([]int, len(films))
	for i, _ := range films {
		urlSlice, ok := m["url"].([]string)
		if !ok {
			return nil, customerror.NotFoundInCache
		}
		movieNum, err := GetMovieNumber(urlSlice)
		if err != nil {
			return nil, err
		}
		movies[i] = movieNum
	}
	p.Movies = movies

	return p, nil
}

func CmToFeet(cm int) int {
	return int(float64(cm) / 30.48)
}

// Convert cm to inches
func CmToInches(cm int) int {
	return int(float64(cm) / 2.54)
}

func GetNumberFromUrl(arr []string) (int, error) {
	if len(arr) == 0 {
		return 0, errors.New("array is empty")
	}

	lastUrl := arr[len(arr)-1]
	numStr := strings.TrimRight(lastUrl, "/")
	numStr = numStr[strings.LastIndex(numStr, "/")+1:]

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}

	return num, nil
}
