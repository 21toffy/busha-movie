package customerror

import (
	"errors"
)

var ErrCacheMiss = errors.New("cache missing error")
var CacheSetError = errors.New("cache set error")
var UnmarshalingError = errors.New("error during unmarshalling")
var OtherCacheError = errors.New("something went wrong")
var NotFoundInCache = errors.New("not found in caches")
var FailedIdConversion = errors.New("Error during conversion of ID")

var DecodeError = errors.New("Decoding error")
var FailedCacheFetch = errors.New("Failed to fetch from cache")

var FailedFetch = errors.New("Failed to fetch from sawpi")
var FailedCacheSetting = errors.New("Failed to set data to cache")
var InvalidReleaseDate = errors.New("invalid release date")
var EmptyString = errors.New("empty url passed")
var InvalidSortParam = errors.New("Invalid sorting parameter")

var ReadBodyError = errors.New("Could not read error")

var FailedTimeConversion = errors.New("An error occured during time conversion")
