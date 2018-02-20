package main

import (
	"flag"
	"fmt"
	"github.com/caius/go-envsensor/internal/envsensor"
	log "github.com/sirupsen/logrus"
	"time"
)

type Configuration struct {
	ProbeDelay    time.Duration
	CacheDuration time.Duration
	SensorPin     int
	SensorVersion int
	WebPort       int
	MQTTBroker    string
	Location      string
	Verbose       bool
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

func main() {
	// Configuration Section
	config := new(Configuration)

	flag.DurationVar(&config.ProbeDelay, "poll", (time.Second * 10), "How often to poll the sensor for a reading. Default 10s")
	flag.DurationVar(&config.CacheDuration, "cache", (time.Second * 60), "Max seconds to cache data for. Default 60s")
	flag.IntVar(&config.SensorPin, "sensor-pin", 17, "GPIO Pin (Physical number) to communicate to sensor via")
	flag.IntVar(&config.SensorVersion, "sensor-version", 11, "Which DHT sensor to talk to. 11 or 22.")
	flag.IntVar(&config.WebPort, "web-port", 8080, "Port for webserver to listen on")
	flag.StringVar(&config.MQTTBroker, "mqtt-server", "", "MQTT server location (eg mqtt.local:1883)")
	flag.StringVar(&config.Location, "location", "test", "Location to report stats from")
	flag.BoolVar(&config.Verbose, "verbose", false, "Verbose output")

	flag.Parse()

	log.Info("Welcome to envsensor, where it is our pleasure to probe you today.")
	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	// Program internals now
	var readingChannels []chan envsensor.Reading
	webChannel := make(chan envsensor.Reading)
	readingChannels = append(readingChannels, webChannel)

	// Wire up MQTT if we've a broker to publish to
	if config.MQTTBroker != "" {
		mqttChannel := make(chan envsensor.Reading)
		readingChannels = append(readingChannels, mqttChannel)

		publisher := envsensor.NewMQTTPublisher(config.MQTTBroker, config.Location)
		go publisher.Start(mqttChannel)
	}

	// Start reading from sensor
	sensor := envsensor.NewDHTSensor(config.SensorVersion, config.SensorPin, config.ProbeDelay)
	go sensor.Start(readingChannels)

	// Serve readings, caching data up to a minute
	port := fmt.Sprintf(":%d", int(config.WebPort))
	server := envsensor.NewWebServer(port, config.CacheDuration)
	server.Start(webChannel)
}
