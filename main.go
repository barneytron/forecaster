package main

import (
	"forecaster/internal/client"
	"forecaster/internal/server"
)

// Write an HTTP server that serves the forecasted weather. Your server should expose an endpoint that:
// Accepts latitude and longitude coordinates
// Returns the short forecast for that area for Today (“Partly Cloudy” etc)
// Returns a characterization of whether the temperature is “hot”, “cold”, or “moderate” (use your discretion on mapping temperatures to each type)
// Use the National Weather Service API Web Service as a data source.

func main() {
	weatherServiceClient := client.NewWeatherServiceClient()
	httpServer := server.NewHttpServer(weatherServiceClient)
	httpServer.Start()
}
