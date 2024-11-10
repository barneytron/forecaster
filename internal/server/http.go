package server

import (
	"encoding/json"
	"fmt"
	"forecaster/internal/client"
	"log"
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

func (h HttpServer) handler(w http.ResponseWriter, r *http.Request) {
	var coordinates Coordinates
	var err error

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&coordinates)

	if coordinates.Latitude == nil || coordinates.Longitude == nil {
		log.Println("invalid request:", r.Body)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	gridForecastUrl, err := h.weatherServiceClient.GetGridForecastUrl(
		*coordinates.Latitude,
		*coordinates.Longitude)
	if err != nil {
		log.Println("error occurred:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	shortForecast, temperature, err := h.weatherServiceClient.GetForecast(gridForecastUrl)
	if err != nil {
		log.Println("error occurred:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, shortForecast, " and ", characterize(temperature), " temperature\n")

}

func (h HttpServer) Start() {
	http.HandleFunc("/forecast", h.handler)
	fmt.Println("Server started and listening")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
