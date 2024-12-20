// Package client_test provides integration tests (not unit tests) for the client package
package client_test

import (
	"forecaster/internal/client"
	"testing"
)

func TestGetGridForecastUrl(t *testing.T) {
	latitude, longitude := 39.7456, -97.0892
	weatherServiceClient := client.NewWeatherServiceClient()
	url, err := weatherServiceClient.GetGridForecastUrl(latitude, longitude)
	if err != nil {
		t.Error("Error occurred:", err)
	}
	if url == "" {
		t.Error("url is an empty string")
	}
	t.Log("gridForecastUrl is", url)
}

func TestGetForecast(t *testing.T) {
	gridForecastUrl := "https://api.weather.gov/gridpoints/TOP/32,81"
	weatherServiceClient := client.NewWeatherServiceClient()
	shortForecast, temperature, err := weatherServiceClient.GetForecast(gridForecastUrl)
	if err != nil {
		t.Error("Error occurred:", err)
	}
	t.Log("shortForecast is", shortForecast)
	t.Log("temperature is", temperature)
}
