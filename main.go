package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// Weather structure
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
		Region  string `json:"region"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int     `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
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

	q := "Denpasar"
	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	url := "https://weatherapi-com.p.rapidapi.com/forecast.json?q=" + q + "&days=3"
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

	// fmt.Println(weather)

	location, curerrent, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf("%s, %s: %.0fC, %s\n",
		location.Name,
		location.Region,
		curerrent.TempC,
		curerrent.Condition.Text)


		for _, hour := range hours {
			date := time.Unix(int64(hour.TimeEpoch), 0)
		
			// Time now & future
			if date.Before(time.Now()) {
				continue
			}
		
			message := fmt.Sprintf(
				"Time: %s - %.0fC, %s, %.0f%% chance of rain\n",
				date.Format("15:04"),
				hour.TempC,
				hour.Condition.Text,
				hour.ChanceOfRain,
			)
		
			if hour.ChanceOfRain < 40 {
				color.Cyan(message)
			} else {
				color.Red(message)
			}
		}		
}
