package api

import (
	"errors"
	"io"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

func afterNow(d time.Time, now time.Time) bool {
	return d.Format(DateFormat) > now.Format(DateFormat)
}

func GetMaskWeek(dstart time.Time, repeat string) ([]int, error) {

	var days []int = []int{}

	strWeek := strings.Split(repeat, ",")

	for _, strDay := range strWeek {
		day, err := strconv.Atoi(strDay)
		if err != nil {
			return nil, err
		}
		if day > 7 || day < 1 {
			return nil, errors.New("Wrong day format")
		}

		if !slices.Contains(days, day) {
			days = append(days, day)
		}
	}

	sort.Ints(days)

	startDay := dstart.Weekday()
	if startDay == time.Sunday {
		startDay = 7
	}

	for i, day := range days {
		days[i] = day - int(startDay)

		if days[i] < 0 {
			days[i] += 7
		}
	}

	sort.Ints(days)

	return days, nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if len(repeat) == 0 {
		return "", errors.New("Empty repeat")
	}

	ds, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	availableTypes := "dyw"

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

		newDate = ds
		for {
			newDate = newDate.AddDate(1, 0, 0)

			if afterNow(newDate, now) {
				break
			}
		}

		return newDate.Format(DateFormat), nil

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

		newDate = ds

		for {
			newDate = newDate.AddDate(0, 0, interval)

			if afterNow(newDate, now) {
				break
			}
		}

		return newDate.Format(DateFormat), nil
	case "w":
		if len(parts) != 2 {
			return "", errors.New("Invalid week repeat format")
		}

		mask, err := GetMaskWeek(ds, parts[1])
		if err != nil {
			return "", err
		}

		newDate = ds

		dx := 0
		i := 0

		for {
			newDate = newDate.AddDate(0, 0, mask[i]+dx)

			if afterNow(newDate, now) {
				break
			}

			i++

			if i == 6 {
				dx += 7
				i = 0
			}

		}

		return newDate.Format(DateFormat), nil
	default:
		return "", errors.New("Unsupported repeat type")
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	nowStr := r.FormValue("now")

	now := time.Now()

	if nowStr != "" {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "Invalid now format", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	io.WriteString(w, nextDate)

}
