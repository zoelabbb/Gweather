package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Weather structure
type Weather struct {
	Location struct {
		Name string `json:"name"`
		Country string `json:"country"`
		Region string `json:"region"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int `json:"code"`
		}`json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Day struct {
				Condition struct {
					Text string `json:"text"`
					Icon string `json:"icon"`
					Code int `json:"code"`
				} `json:"condition"`
			} `json:"day"`
		} `json:"forecastday"`
	}`json:"forecast"`
}

func main() {

	// Load env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Get env file
	rapidAPIKey := os.Getenv("RAPID_API_KEY")
	rapidAPIHost := os.Getenv("RAPID_API_HOST")

	url := "https://weatherapi-com.p.rapidapi.com/forecast.json?q=Indonesia&days=3"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("X-RapidAPI-Key", rapidAPIKey)
	req.Header.Add("X-RapidAPI-Host", rapidAPIHost)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("API call failed with " + res.Status)
	} else {
		fmt.Println("API call succeeded with " + res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(body))

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	fmt.Println(weather)
}