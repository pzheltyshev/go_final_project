package api

import (
	"errors"
	"fmt"
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

func GetMonths(repeat string) ([]int, error) {

	if len(repeat) == 0 {
		return []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, nil
	}

	months := []int{}

	monthsStr := strings.Split(repeat, ",")

	for _, monthStr := range monthsStr {
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			return []int{}, err
		}

		if month < 1 || month > 12 {
			return []int{}, errors.New("Wrong month format")
		}

		months = append(months, month)
	}

	sort.Ints(months)

	return months, nil
}

func GetNums(repeat string) ([]int, error) {

}

func GetMonthDates(startDate time.Time, curDate time.Time, months []int, nums []int) []int {
	// стартовая дата нужна чтобы сравнить дату в первой итерации потому что числа могут быть меньше чем день в текущем месяце
	// тек месяц используется для понимания в каком месяце мы на данный момент находимся ! лучше заменить на дату!
	// нужно проверить что текущий обрабатываемый месяц есть в month . если нет, то пустое возвращаем и по сути на месяц прокрутим вперед
	// из nums нужно получить месяца для текущего прорабатываемого месяца
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

		count := len(mask)

		for {
			testDate := newDate.AddDate(0, 0, mask[i]+dx)

			fmt.Println("testDate " + testDate.Format(DateFormat) + " i " + strconv.Itoa(i) + " dx " + strconv.Itoa(dx) +
				" mask " + strconv.Itoa(mask[i]))
			if afterNow(testDate, now) {
				newDate = testDate
				break
			}

			if i == count-1 {
				dx += 7
				i = 0
			} else {
				i++
			}

		}

		return newDate.Format(DateFormat), nil
	case "m":
		if len(parts) != 2 {
			return "", errors.New("Invalid month repeat format")
		}

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
