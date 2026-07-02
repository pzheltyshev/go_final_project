package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func afterNow(d time.Time, now time.Time) bool {
	return d.Format("20060102") > now.Format("20060102")
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if len(repeat) == 0 {
		return "", errors.New("Empty repeat")
	}

	ds, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}

	availableTypes := "dy"

	parts := strings.Split(repeat, " ")

	if !strings.Contains(availableTypes, parts[0]) {
		return "", errors.New("Invalid repeat type")
	}

	var newDate time.Time

	switch parts[0] {
	case "y":
		if len(parts) > 1 {
			return "", errors.New("Invalid year repeat format")
		}

		for {
			newDate = ds.AddDate(1, 0, 0)

			if afterNow(newDate, now) {
				break
			}
		}

		return newDate.Format("20060102"), nil

	case "d":
		if len(parts) != 2 {
			return "", errors.New("Invalid day repeat format")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}
		if interval <= 0 || interval > 400 {
			return "", errors.New("Interval must be from 1 to 400")
		}

		for {
			newDate = ds.AddDate(0, 0, interval)

			if afterNow(newDate, now) {
				break
			}
		}

		return newDate.Format("20060102"), nil
	default:
		return "", errors.New("Unsupported repeat type")
	}
}
