package utils

import (
	// "fmt"
	"errors"

	"github.com/21toffy/busha-movie/internal/customerror"

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
