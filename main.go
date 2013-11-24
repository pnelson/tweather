package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var weatherUrl = fmt.Sprintf("http://api.wunderground.com/api/%s/%%s/q/%s/%s.json",
	os.Getenv("WUNDERGROUND_API_KEY"),
	os.Getenv("TWEATHER_COUNTRY"),
	os.Getenv("TWEATHER_CITY"),
)

type Weather struct {
	Observation struct {
		Condition      string `json:"weather"`
		Temperature    int    `json:"temp_c"`
		TemperatureMod string `json:"feelslike_c"`
		Wind           int    `json:"wind_kph"`
		Visibility     string `json:"visibility_km"`
	} `json:"current_observation"`
	Forecast struct {
		Simple struct {
			Day []struct {
				High struct {
					Celsius string
				}
				Low struct {
					Celsius string
				}
				Precipitation int
			} `json:"forecastday"`
		} `json:"simpleforecast"`
	}
}

func (w *Weather) String() string {
	o := &w.Observation
	f := &w.Forecast.Simple.Day[0]
	return fmt.Sprintf("%s at %dºC (%sºC)\nWind %dkm/h — Visibility %skm\n▼%sºC ▲%sºC ☂%d%%",
		o.Condition, o.Temperature, o.TemperatureMod,
		o.Wind, o.Visibility,
		f.Low.Celsius, f.High.Celsius, f.Precipitation,
	)
}

func weatherFrom(features []string) (*Weather, error) {
	url := fmt.Sprintf(weatherUrl, strings.Join(features, "/"))
	rv, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rv.Body.Close()
	if rv.StatusCode != http.StatusOK {
		return nil, errors.New(rv.Status)
	}
	weather := new(Weather)
	err = json.NewDecoder(rv.Body).Decode(weather)
	if err != nil {
		return nil, err
	}
	return weather, nil
}

func main() {
	weather, err := weatherFrom([]string{"conditions", "forecast"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(weather)
}
