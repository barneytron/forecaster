package server

import (
	"fmt"
	"log"
	"net/http"
)

type HttpServer struct {
	httpHandler ForecastHandler
}

func NewHttpServer(handler ForecastHandler) *HttpServer {
	return &HttpServer{
		httpHandler: handler,
	}
}

func (h HttpServer) Start() {
	http.HandleFunc("/forecast", h.httpHandler.HandleForecastRequest)
	fmt.Println("Server started and listening")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
