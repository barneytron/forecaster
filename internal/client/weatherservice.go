package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type WeatherServiceClient interface {
	GetGridForecastUrl(latitude float64, longitude float64) (string, error)
	GetForecast(gridForecastUrl string) (string, float64, error)
}

type Client struct {
	httpClient http.Client
}

func NewWeatherServiceClient() WeatherServiceClient {
	return &Client{
		httpClient: http.Client{
			Timeout: time.Second * 60,
		},
	}
}

// GetForecast calls the grid forecast api endpoint.
// Using the first "Period" of the forecast, the method returns the shortForecast string, temperature, and any error encountered.
func (c Client) GetForecast(gridForecastUrl string) (string, float64, error) {

	url := fmt.Sprintf("%s/forecast", gridForecastUrl)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return "", 0, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return "", 0, err
	}

	var result map[string]any

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	if err != nil {
		log.Println(err)
		return "", 0, err
	}

	firstPeriod, ok := result["properties"].(map[string]any)["periods"].([]any)[0].(map[string]any)
	if !ok {
		return "", 0, errors.New("periods field is missing in properties")
	}
	
	shortForecast := fmt.Sprint(firstPeriod["shortForecast"])
	temperature, ok := firstPeriod["temperature"].(float64)
	if !ok {
		return "", 0, errors.New("forecastGridData field is missing in properties")
	}

	return shortForecast, temperature, nil
}

func (c Client) GetGridForecastUrl(latitude float64, longitude float64) (string, error) {
	url := fmt.Sprintf("https://api.weather.gov/points/%f,%f", latitude, longitude)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var result map[string]any
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	if err != nil {
		log.Println(err)
		return "", err
	}

	forecastGridDataUrl, ok := result["properties"].(map[string]any)["forecastGridData"].(string)
	if !ok {
		return "", errors.New("forecastGridData field is missing in properties")
	}

	return forecastGridDataUrl, nil
}
