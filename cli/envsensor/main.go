package main

import (
	"flag"
	"fmt"
	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/caius/go-envsensor/internal/envsensor"
	"os"
	"time"
)

type Configuration struct {
	ProbeDelay    time.Duration
	CacheDuration time.Duration
	SensorPin     int
	SensorVersion int
	WebPort       int
	Logger        bunyan.Logger
	Verbose       bool
}

func createLogger(verbose bool) bunyan.Logger {
	logLevel := bunyan.LogLevelInfo
	if verbose {
		logLevel = bunyan.LogLevelDebug
	}
	loggerConfig := bunyan.Config{
		Name:   "envsensor",
		Stream: os.Stdout,
		Level:  logLevel,
	}

	log, err := bunyan.CreateLogger(loggerConfig)
	if err != nil {
		panic(err)
	}

	return log
}

func main() {
	// Configuration Section
	config := new(Configuration)

	flag.DurationVar(&config.ProbeDelay, "poll", (time.Second * 10), "How often to poll the sensor for a reading. Default 10s")
	flag.DurationVar(&config.CacheDuration, "cache", (time.Second * 60), "Max seconds to cache data for. Default 60s")
	flag.IntVar(&config.SensorPin, "sensor-pin", 17, "GPIO Pin (Physical number) to communicate to sensor via")
	flag.IntVar(&config.SensorVersion, "sensor-version", 11, "Which DHT sensor to talk to. 11 or 22.")
	flag.IntVar(&config.WebPort, "web-port", 8080, "Port for webserver to listen on")
	flag.BoolVar(&config.Verbose, "verbose", false, "Verbose output")

	flag.Parse()

	config.Logger = createLogger(config.Verbose)

	config.Logger.Info("Welcome to envsensor, where it is our pleasure to probe you today.")

	// Program internals now
	readingsChan := make(chan envsensor.Reading)

	// Start reading from sensor
	sensor := envsensor.NewDHTSensor(config.SensorVersion, config.SensorPin, config.ProbeDelay)
	go sensor.Start(readingsChan)

	// Serve readings, caching data up to a minute
	server := envsensor.WebServer{}
	port := fmt.Sprintf(":%d", int(config.WebPort))
	server.Start(readingsChan, port, config.CacheDuration)
}
