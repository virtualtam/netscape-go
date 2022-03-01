package decoder

import (
	"strconv"
	"time"
)

func decodeDate(input string) (time.Time, error) {
	// First, attempt to parse the date as a UNIX epoch, which is the most
	// commonly used format.
	unixTime, err := strconv.ParseInt(input, 10, 64)
	if err == nil {
		return decodeUnixDate(unixTime), nil
	}

	// Attempt to parse the date as RFC3339
	date, err := decodeRFC3339Date(input)
	if err == nil {
		return date, nil
	}

	return time.Time{}, err
}

// decodeRFC3339Date returns the time.Time corresponding to a RFC3339
// representation.
func decodeRFC3339Date(input string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, input)
	if err == nil {
		return date.UTC(), nil
	}

	date, err = time.Parse(time.RFC3339Nano, input)
	if err == nil {
		return date.UTC(), nil
	}

	return time.Time{}, err
}

// decodeUnixDate returns the time.Time corresponding to a UNIX timestamp.
//
// Dates are usually specified in seconds, but some browsers and bookmarking
// services may use milliseconds, microseconds or nanoseconds.
//
// To address these cases, we ensure the resulting time.Time is comprised in a
// reasonable interval (ie not further than N years in the future).
func decodeUnixDate(unixTime int64) time.Time {
	const rangeYears = 30

	date := time.Unix(unixTime, 0).UTC()

	if date.After(time.Now().AddDate(rangeYears, 0, 0)) {
		date = time.UnixMilli(unixTime).UTC()
	}

	if date.After(time.Now().AddDate(rangeYears, 0, 0)) {
		date = time.UnixMicro(unixTime).UTC()
	}

	if date.After(time.Now().AddDate(rangeYears, 0, 0)) {
		date = time.Unix(0, unixTime).UTC()
	}

	return date
}
