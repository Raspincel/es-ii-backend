package utils

type Theme struct {
	Theme string
	Words []string
}

type Day struct {
	Day    string
	Themes []Theme
}

type DayRow struct {
	Groups string
}
