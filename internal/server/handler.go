package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"forecaster/internal/client"
	"log"
	"math"
	"net/http"
)

type ForecastHandler interface {
	HandleForecastRequest(w http.ResponseWriter, r *http.Request)
}

type HttpHandler struct {
	weatherServiceClient client.WeatherServiceGetter
}

func NewHttpHandler(client client.WeatherServiceGetter) *HttpHandler {
	return &HttpHandler{
		weatherServiceClient: client,
	}
}

type Coordinates struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

func CharacterizeTemperature(temperature float64) string {
	if temperature >= 90 {
		return "hot"
	} else if temperature >= 70 {
		return "moderate"
	} else {
		return "cold"
	}
}

// Latitude must be a number between -90 and 90
func IsLatitudeValid(latitude float64) bool {
	abs := math.Abs(latitude)
	return abs <= 90
}

// Longitude must a number between -180 and 180
func IsLongitudeValid(longitude float64) bool {
	abs := math.Abs(longitude)
	return abs <= 180
}

func (h *HttpHandler) HandleForecastRequest(w http.ResponseWriter, r *http.Request) {
	var coordinates Coordinates
	var err error

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&coordinates)
	if err != nil {
		log.Println("error occurred decoding json payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cannot decode json payload"))
		return
	}

	if coordinates.Latitude == nil || coordinates.Longitude == nil ||
		!IsLatitudeValid(*coordinates.Latitude) || !IsLongitudeValid(*coordinates.Longitude) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Latitude must be between -90 and 90, and longitude must be between -180 and 180"))
		return
	}

	gridForecastUrl, err := h.weatherServiceClient.GetGridForecastUrl(
		*coordinates.Latitude,
		*coordinates.Longitude)
	if err != nil {
		if errors.Is(err, client.ErrForecastGridDataNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		log.Println("internal server error occurred:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortForecast, temperature, err := h.weatherServiceClient.GetForecast(gridForecastUrl)
	if err != nil {
		log.Println("internal server error occurred:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, shortForecast, " and ", CharacterizeTemperature(temperature), " temperature\n")
}
