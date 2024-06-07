package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/raspincel/es-ii-backend/internal/db"
	"github.com/raspincel/es-ii-backend/internal/gpt"
	"github.com/raspincel/es-ii-backend/internal/utils"
)

var client *http.Client

func init() {
	client = &http.Client{}
}

func StartServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /day/{date}", getDate)
	mux.HandleFunc("GET /range", getRange)
	mux.HandleFunc("GET /range/validate/{date}", isDateInRange)

	port := utils.GetEnv("PORT")
	fmt.Println("Running on port " + port)
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}

func getDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	day := r.URL.Path[strings.LastIndex(r.URL.Path, "/"):]
	day = strings.Replace(day, "/", "", 1)

	date, err, _, _ := utils.ParseDate(day)

	if err == "" {
		w.Write([]byte("Invalid date format (yyyy-mm-dd)"))
		return
	}

	themes := db.GetDay(client, date)
	fmt.Println("Themes from db", themes)
	attempts := 0

	for themes == nil {
		if attempts == 3 {
			w.Write([]byte("Error getting valid themes"))
			return
		}

		content := gpt.RequestGroups()

		if content == "" {
			attempts++
			continue
		}

		themes = utils.ParseGroups(content)

		if themes != nil {
			db.SaveDay(client, day, content)
			break
		}

		attempts++
	}

	answer := utils.Day{
		Day:    day,
		Themes: themes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}

func getRange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sixMonthsAgo := utils.GetSixMonthsAgo()
	firstYear := sixMonthsAgo.Year()
	firstMonth := int(sixMonthsAgo.Month())
	firstDay := sixMonthsAgo.Day()

	now := time.Now()
	lastYear := now.Year()
	lastMonth := int(now.Month())
	lastDay := now.Day()

	firstDate, _, _, _ := utils.ParseDate(fmt.Sprintf("%d-%d-%d", firstYear, firstMonth, firstDay))
	lastDate, _, _, _ := utils.ParseDate(fmt.Sprintf("%d-%d-%d", lastYear, lastMonth, lastDay))

	jsonString := fmt.Sprintf(`{
		"firstDay": "%s",
		"lastDay": "%s"
	}`, firstDate, lastDate)
	w.Write([]byte(jsonString))
}

func isDateInRange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	unformattedDay := r.URL.Path[strings.LastIndex(r.URL.Path, "/"):]
	unformattedDay = strings.Replace(unformattedDay, "/", "", 1)

	date, y, m, d := utils.ParseDate(unformattedDay)

	if date == "" {
		w.Write([]byte(`{
			"date": "", 
			"isValid": false,
			"error": "Invalid date format (yyyy-mm-dd)"
		}`))
		return
	}

	year := utils.ConvertToNum(y)
	month := utils.ConvertToNum(m)
	day := utils.ConvertToNum(d)

	newDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	sixMonthsAgo := utils.GetSixMonthsAgo()
	comparison := sixMonthsAgo.Compare(newDate)

	isBeforeToday := newDate.Before(time.Now())

	jsonString := fmt.Sprintf(`{"date": "%s", "isValid": %t}`, newDate, comparison <= 0 && isBeforeToday)
	w.Write([]byte(jsonString))
}
