package main

import (
	"fmt"
	"github.com/caius/go-envsensor/internal/envsensor"
	"time"
)

func main() {
	fmt.Println("Welcome to temphumid, where it is our pleasure to probe you today.")

	// Configuration Section
	// TODO: turn these into cli flags
	// TODO: sensorVersion := 11 // 11 or 22 (DHT11, DHT22)
	gpioPin := 17
	probeDelay := 10              // seconds
	cacheDelay := time.Minute * 1 // return cached data for 1 minute

	// Program internals now
	readingsChan := make(chan envsensor.Reading)

	// Grab the readings every ten seconds
	sensor := envsensor.DHTSensor{}
	go sensor.Start(gpioPin, readingsChan, probeDelay)

	// Serve readings, caching data up to a minute
	server := envsensor.WebServer{}
	server.Start(readingsChan, cacheDelay)
}
