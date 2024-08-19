package helpers

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// 2024-08-19T16:31
func DecomposeStringDate(strDate string) (year, month, day, hour, min, sec int, err error) {
	regex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}$`)
	if !regex.MatchString(strDate) {
		return 0, 0, 0, 0, 0, 0, errors.New("date format error")
	}
	firstArr := strings.Split(strDate, "T")
	date := firstArr[0]
	tim := firstArr[1]
	sep1 := strings.Split(date, "-")
	yearStr := sep1[0]
	monthStr := sep1[1]
	dayStr := sep1[2]
	sep2 := strings.Split(tim, ":")

	hourStr := sep2[0]
	minStr := sep2[1]
	secStr := sep2[2]
	year, err = strconv.Atoi(yearStr)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}
	month, err = strconv.Atoi(monthStr)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}
	day, err = strconv.Atoi(dayStr)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}
	hour, err = strconv.Atoi(hourStr)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}
	min, err = strconv.Atoi(minStr)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}
	sec, err = strconv.Atoi(secStr)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}
	return year, month, day, hour, min, sec, nil
}
