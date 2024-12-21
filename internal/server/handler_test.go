package server_test

import (
	"errors"
	"fmt"
	"forecaster/internal/server"
	"forecaster/mocks"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	path                  = "/forecast"
	latitude              = 38.672222
	longitude             = -121.157778
	forecastUrl           = "http://weatherserver"
	shortForecast         = "Sunny"
	temperature   float64 = 49
)

type HandlerTest struct {
	suite.Suite
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerTest))
}

func (suite *HandlerTest) TestHandleForecastRequest_HappyPath_ShouldReturn200() {
	client := new(mocks.WeatherServiceGetter)
	client.On("GetGridForecastUrl", mock.Anything, mock.Anything).Return(forecastUrl, nil)
	client.On("GetForecast", forecastUrl).Return(shortForecast, temperature, nil)
	recorder := httptest.NewRecorder()
	jsonPayload := fmt.Sprintf(`{"latitude":%v, "longitude":%v}`, latitude, longitude)
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(jsonPayload))
	handler := server.NewHttpHandler(client)

	handler.HandleForecastRequest(recorder, request)

	suite.Equal(http.StatusOK, recorder.Result().StatusCode)
	suite.Equal("Sunny and cold temperature", strings.TrimSpace(recorder.Body.String()))
}

func (suite *HandlerTest) TestHandleForecastRequest_GetGridForecastUrlErrors_ShouldReturn500() {
	client := new(mocks.WeatherServiceGetter)
	client.On("GetGridForecastUrl", mock.Anything, mock.Anything).Return("", errors.New("unexpected error"))
	recorder := httptest.NewRecorder()
	jsonPayload := fmt.Sprintf(`{"latitude":%v, "longitude":%v}`, latitude, longitude)
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(jsonPayload))
	handler := server.NewHttpHandler(client)

	handler.HandleForecastRequest(recorder, request)

	suite.Equal(http.StatusInternalServerError, recorder.Result().StatusCode)
}

func (suite *HandlerTest) TestHandleForecastRequest_InvalidJsonPayload_ShouldReturn400() {
	client := new(mocks.WeatherServiceGetter)
	client.On("GetGridForecastUrl", mock.Anything, mock.Anything).Return(forecastUrl, nil)
	client.On("GetForecast", forecastUrl).Return(shortForecast, temperature, nil)
	recorder := httptest.NewRecorder()
	jsonPayload := fmt.Sprintf(`{invalid"latitude":%v "longitude":%v`, latitude, longitude)
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(jsonPayload))
	handler := server.NewHttpHandler(client)

	handler.HandleForecastRequest(recorder, request)

	suite.Equal(http.StatusBadRequest, recorder.Result().StatusCode)
}

func (suite *HandlerTest) TestIsLatitudeValid_InvalidLatitude_ShouldReturnFalse() {
	latitude := float64(100)

	got := server.IsLatitudeValid(latitude)

	suite.Equal(false, got)
}

func (suite *HandlerTest) TestIsLatitudeValid_ValidLatitude_ShouldReturnTrue() {
	latitude := float64(-90)

	got := server.IsLatitudeValid(latitude)

	suite.Equal(true, got)
}

func (suite *HandlerTest) TestCharacterizeTemperature_GreaterThan90_ShouldReturnHot() {
	temperature := float64(100)

	got := server.CharacterizeTemperature(temperature)

	suite.Equal("hot", got)
}

func (suite *HandlerTest) TestCharacterizeTemperature_Between70And89_ShouldReturnModerate() {
	temperature := float64(70)

	got := server.CharacterizeTemperature(temperature)

	suite.Equal("moderate", got)
}

func (suite *HandlerTest) TestCharacterizeTemperature_LessThan70_ShouldReturnCold() {
	temperature := float64(60)

	got := server.CharacterizeTemperature(temperature)

	suite.Equal("cold", got)
}
