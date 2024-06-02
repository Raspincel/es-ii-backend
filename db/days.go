package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Raspincel/es-ii-backend/utils"
)

func GetDay(client *http.Client, day string) []utils.Theme {
	supabaseUrl := utils.GetEnv("SUPABASE_URL")
	supabaseKey := utils.GetEnv("SUPABASE_KEY")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/days?day=eq.%s&select=groups", supabaseUrl, day), nil)

	if err != nil {
		log.Fatalln(err)
		return nil
	}

	req.Header.Set("apikey", supabaseKey)
	res, err := client.Do(req)

	defer func() {
		res.Body.Close()
	}()

	if err != nil {
		log.Fatalln(err)
		return nil
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatalln(err)
		return nil
	}

	bodyStr := string(body[:])
	replacer := strings.NewReplacer("[", "", "]", "")

	var jsonMap utils.DayRow
	json.Unmarshal([]byte(replacer.Replace(bodyStr)), &jsonMap)

	themes := utils.ParseGroups(jsonMap.Groups)

	return themes
}

func SaveDay(client *http.Client, day string, content string) {
	supabaseUrl := utils.GetEnv("SUPABASE_URL")
	supabaseKey := utils.GetEnv("SUPABASE_KEY")

	jsonBody := fmt.Sprintf(`{"day": "%s", "groups": "%s"}`, day, content)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/days", supabaseUrl), bytes.NewBuffer([]byte(jsonBody)))

	if err != nil {
		log.Fatalln(err)
		return
	}

	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	defer func() {
		res.Body.Close()
	}()

	if err != nil {
		log.Fatalln(err)
		return
	}

	// body, err := io.ReadAll(res.Body)

	// if err != nil {
	// 	log.Fatalln(err)
	// 	return
	// }

	// bodyStr := string(body[:])
	// replacer := strings.NewReplacer("[", "", "]", "")

	// fmt.Println(replacer.Replace(bodyStr))
}
