package main

import (
	"flag"
	"fmt"
	"github.com/caius/go-envsensor/internal/envsensor"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type Configuration struct {
	CacheDuration time.Duration
	Location      string
	MQTTBroker    string
	MQTTEnabled   bool
	ProbeDelay    time.Duration
	SensorPin     int
	SensorVersion int
	Valid         bool
	Verbose       bool
	WebEnabled    bool
	WebPort       int
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

func main() {
	// Configuration Section
	config := Configuration{Valid: true}

	// Sensor Configuration
	flag.DurationVar(&config.ProbeDelay, "poll", (time.Second * 10), "How often to poll the sensor for a reading.")
	flag.DurationVar(&config.CacheDuration, "cache", (time.Second * 60), "Max seconds to cache data for.")
	flag.IntVar(&config.SensorPin, "sensor-pin", 0, "GPIO Pin (Physical number) to communicate to sensor on")
	flag.IntVar(&config.SensorVersion, "sensor-version", 11, "Which DHT sensor to talk to. 11 or 22.")
	flag.StringVar(&config.Location, "location", "", "Location identifier for emitted readings")

	// HTTP Configuration
	flag.BoolVar(&config.WebEnabled, "web", false, "Enable webserver")
	flag.IntVar(&config.WebPort, "web-port", 8080, "Port for webserver to listen on")

	// MQTT Configuration
	flag.BoolVar(&config.MQTTEnabled, "mqtt", false, "Enable MQTT publishing")
	flag.StringVar(&config.MQTTBroker, "mqtt-broker", "", "MQTT server address (eg mqtt.local:1883)")
	flag.BoolVar(&config.Verbose, "verbose", false, "Verbose output")

	flag.Parse()

	// Check configuration is correct
	if config.SensorVersion != 11 && config.SensorVersion != 22 {
		log.Error("--sensor-version must be 11 or 22")
		config.Valid = false
	}

	if config.Location == "" {
		log.Error("--location is a required argument")
		config.Valid = false
	}

	if config.Valid != true {
		log.Error("Configuration errors, please see above and fix")
		os.Exit(1)
	}

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
