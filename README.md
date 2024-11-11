# forecaster

## About
forecaster is an HTTP server that serves the forecasted weather using the National Weather Service API Web Service as the data source.

The endpoint /forecast accepts a POST request having latitude and longitude coordinates and returns the short forecast for that area and a characterization of whether the temperature (F) is “hot”, “cold”, or “moderate”. Hot is greater than 90F, moderate is 70F to 89F, and cold is less than 70F.

The server is hardcoded to listen on port 8080.

No third party Go packages were used in this project.

## Requirements
go version go1.23.3 linux/amd64

## How to run
Inside the *forecaster* directory, execute:

```console
go run ./...
```

## How to test
With curl:
```console
curl --header "Content-Type: application/json" --data '{"latitude":39.7456,"longitude":-97.0892}' http://localhost:8080/forecast
```
```console
Sunny and cold temperature
```

With HTTPie:
```console
http POST http://localhost:8080/forecast latitude:=39.7456 longitude:=-97.0892
```
```console
HTTP/1.1 200 OK
Content-Length: 27
Content-Type: text/plain; charset=utf-8
Date: Sun, 10 Nov 2024 20:32:03 GMT

Sunny and cold temperature
```
