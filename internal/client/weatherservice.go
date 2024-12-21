package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type PointMetadata struct {
	PointMetaDataProperties `json:"properties"`
}

type PointMetaDataProperties struct {
	ForecastGridData string `json:"forecastGridData"`
}

type Forecast struct {
	ForecastProperties `json:"properties"`
}

type ForecastProperties struct {
	Periods []Period `json:"periods"`
}

type Period struct {
	ShortForecast string  `json:"shortForecast"`
	Temperature   float64 `json:"temperature"`
}

var ErrForecastGridDataNotFound = errors.New("forecastGridData not available for given coordinates")
var ErrForecastGridDataMissing = errors.New("forecastGridData is missing in result")

type WeatherServiceGetter interface {
	GetGridForecastUrl(latitude float64, longitude float64) (string, error)
	GetForecast(gridForecastUrl string) (string, float64, error)
}

type Client struct {
	httpClient http.Client
}

func NewWeatherServiceClient() *Client {
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
		return "", 0, fmt.Errorf("[http.NewRequest] error: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return "", 0, fmt.Errorf("[httpClient.Do] error: %w", err)
	}
	defer res.Body.Close()

	var forecast Forecast
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&forecast)
	if err != nil {
		return "", 0, fmt.Errorf("[httpClient.Do] error: %w", err)
	}

	if len(forecast.ForecastProperties.Periods) == 0 {
		return "", 0, errors.New("periods field is missing in properties")
	}

	firstPeriod := forecast.ForecastProperties.Periods[0]
	shortForecast := firstPeriod.ShortForecast
	temperature := firstPeriod.Temperature

	return shortForecast, temperature, nil
}

func (c Client) GetGridForecastUrl(latitude float64, longitude float64) (string, error) {
	url := fmt.Sprintf("https://api.weather.gov/points/%f,%f", latitude, longitude)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("[httpClient.Do] error: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		log.Println("GET", url, "returned 404")
		return "", ErrForecastGridDataNotFound
	}

	var pointMetadata PointMetadata
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&pointMetadata)
	if err != nil {
		return "", fmt.Errorf("[Decode] error: %w", err)
	}

	forecastGridDataUrl := pointMetadata.PointMetaDataProperties.ForecastGridData
	if forecastGridDataUrl == "" {
		log.Println("GET", url, "result is missing forecastGridData")
		return "", ErrForecastGridDataMissing
	}

	return forecastGridDataUrl, nil
}
