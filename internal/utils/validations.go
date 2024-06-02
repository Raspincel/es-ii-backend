package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseGroups(groups string) []Theme {
	groupsArray := strings.Split(groups, ";")

	if len(groupsArray) != 4 {
		return nil
	}

	var themes []Theme = make([]Theme, 4)

	for index, group := range groupsArray {
		if len(group) == 0 {
			return nil
		}

		groupArray := strings.Split(group, ": ")

		if len(groupArray) != 2 {
			return nil
		}

		wordsArray := strings.Split(groupArray[1], ",")

		if len(wordsArray) != 4 {
			return nil
		}

		for _, word := range wordsArray {
			if len(word) == 0 {
				return nil
			}
		}

		themes[index] = Theme{
			Theme: groupArray[0],
			Words: wordsArray,
		}
	}

	return themes
}

func ParseDate(date string) (string, year, month, day string) {
	splitDate := strings.SplitN(date, "-", 3)

	if len(splitDate) != 3 {
		return "", "", "", ""
	}

	year, month, day = splitDate[0], splitDate[1], splitDate[2]

	if len(year) != 4 || len(month) > 2 || len(day) > 2 {
		return "", "", "", ""
	}

	yearInt, err := strconv.Atoi(year)

	if err != nil {
		return "", "", "", ""
	}

	monthInt, err := strconv.Atoi(month)

	if err != nil {
		return "", "", "", ""
	}

	dayInt, err := strconv.Atoi(day)

	if err != nil {
		return "", "", "", ""
	}

	if yearInt < 0 || monthInt < 1 || monthInt > 12 || dayInt < 1 || dayInt > 31 {
		return "", "", "", ""
	}

	if len(day) == 1 {
		day = "0" + day
	}

	if len(month) == 1 {
		month = "0" + month
	}

	return fmt.Sprintf("%s-%s-%s", year, month, day), year, month, day
}
