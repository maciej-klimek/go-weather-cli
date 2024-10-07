package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		WindKmh   float64 `json:"wind_kph"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		ForecastDay []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
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
	var userInput string = "Poland,%20Krakow"
	handleHttpRequest(userInput)
}

func handleHttpRequest(userInput string) {
	resp, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=0dae3e9126f34758a55201836240710&q=" + userInput + "&days=1&aqi=yes&alerts=no")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("Error while fetching weather data")
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var weatherForecast Weather

	err = json.Unmarshal(body, &weatherForecast)
	if err != nil {
		panic(err)
	}
	location, current, hours := weatherForecast.Location, weatherForecast.Current, weatherForecast.Forecast.ForecastDay[0].Hour

	fmt.Printf(
		"%s, %s, %s: %.1fC, %s\n -----------------------------------\n",
		time.Now().Format("15:05"),
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)
		fmt.Printf(
			"%s, %.1fC, Rain: %.1f%%, %s \n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)
	}
}
