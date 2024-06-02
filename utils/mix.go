package utils

import (
	"strconv"
	"time"
)

func ConvertToNum(date string) int {
	num, err := strconv.Atoi(date)

	if err != nil {
		panic(err)
	}

	return num
}

func GetSixMonthsAgo() time.Time {
	t := time.Now().AddDate(0, -6, 0)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
