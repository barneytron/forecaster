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

type Coordinates struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type HttpServer struct {
	weatherServiceClient client.WeatherServiceClient
}

func NewHttpServer(w client.WeatherServiceClient) *HttpServer {
	return &HttpServer{
		weatherServiceClient: w,
	}
}

func characterize(temperature float64) string {
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

func (h HttpServer) handler(w http.ResponseWriter, r *http.Request) {
	var coordinates Coordinates
	var err error

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&coordinates)

	if coordinates.Latitude == nil || coordinates.Longitude == nil ||
		!IsLatitudeValid(*coordinates.Latitude) || !IsLongitudeValid(*coordinates.Longitude) {
		log.Println("invalid request:", r.Body)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	gridForecastUrl, err := h.weatherServiceClient.GetGridForecastUrl(
		*coordinates.Latitude,
		*coordinates.Longitude)
	if err != nil {
		log.Println("error occurred:", err)
		if errors.Is(err, client.ErrForecastGridDataNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortForecast, temperature, err := h.weatherServiceClient.GetForecast(gridForecastUrl)
	if err != nil {
		log.Println("error occurred:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, shortForecast, " and ", characterize(temperature), " temperature\n")

}

func (h HttpServer) Start() {
	http.HandleFunc("/forecast", h.handler)
	fmt.Println("Server started and listening")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
