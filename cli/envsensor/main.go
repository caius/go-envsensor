package main

import (
	"flag"
	"fmt"
	"github.com/caius/go-envsensor/internal/envsensor"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
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
	flag.DurationVar(&config.CacheDuration, "cache", (time.Second * 60), "Max seconds to cache data for.")
	flag.DurationVar(&config.ProbeDelay, "poll", (time.Second * 10), "How often to poll the sensor for a reading.")
	flag.IntVar(&config.SensorPin, "sensor-pin", 0, "GPIO Pin (Physical number) to communicate to sensor on")
	flag.IntVar(&config.SensorVersion, "sensor-version", 11, "Which DHT sensor to talk to. 11 or 22.")
	flag.StringVar(&config.Location, "location", "", "Location identifier for emitted readings")

	// HTTP Configuration
	flag.BoolVar(&config.WebEnabled, "web", false, "Enable webserver")
	flag.IntVar(&config.WebPort, "web-port", 8080, "Port for webserver to listen on")

	// MQTT Configuration
	flag.BoolVar(&config.MQTTEnabled, "mqtt", false, "Enable MQTT publishing")
	flag.StringVar(&config.MQTTBroker, "mqtt-broker", "", "MQTT server (eg mqtt.local:1883)")

	// Other configuration
	flag.BoolVar(&config.Verbose, "verbose", false, "Verbose output")

	flag.Parse()

	log.Info("Welcome to envsensor, where it is our pleasure to probe you today.")
	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	// Check configuration is correct
	if config.SensorVersion != 11 && config.SensorVersion != 22 {
		log.Error("--sensor-version must be 11 or 22")
		config.Valid = false
	}

	if config.Location == "" {
		log.Error("--location is a required argument")
		config.Valid = false
	}

	if config.MQTTEnabled && config.MQTTBroker == "" {
		log.Error("--mqtt-broker is a required argument when --mqtt is passed")
		config.Valid = false
	}

	if config.Valid != true {
		log.Error("Configuration errors, please fix. Check logs for more information.")
		os.Exit(1)
	}

	// Everything that wants to listen needs to put a channel in here
	var readingChannels []chan envsensor.Reading

	if config.MQTTEnabled == false && config.WebEnabled == false {
		log.Info("Neither MQTT nor Web outputs are enabled. Not gonna do much.")
	}

	var publisher envsensor.MQTTPublisher
	var webserver envsensor.WebServer

	// Wire up MQTT if we've a broker to publish to
	if config.MQTTEnabled {
		mqttChannel := make(chan envsensor.Reading)
		readingChannels = append(readingChannels, mqttChannel)

		publisher = envsensor.NewMQTTPublisher(config.MQTTBroker, config.Location)
		go publisher.Start(mqttChannel)
	}

	if config.WebEnabled {
		webChannel := make(chan envsensor.Reading)
		readingChannels = append(readingChannels, webChannel)

		// Serve readings, caching data up to a minute
		port := fmt.Sprintf(":%d", int(config.WebPort))
		webserver = envsensor.NewWebServer(port, config.CacheDuration)
		go webserver.Start(webChannel)
	}

	// And finally kick the sensor off (we have our reading )
	sensor := envsensor.NewDHTSensor(config.SensorVersion, config.SensorPin, config.ProbeDelay)
	go sensor.Start(readingChannels)

	// Trap and cleanup on interrupt (^C)
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			log.Info("Received interrupt, bringing everything down")

			sensor.Stop()

			if config.MQTTEnabled {
				publisher.Stop()
			}

			if config.WebEnabled {
				webserver.Stop()
			}

			cleanupDone <- true
		}
	}()
	<-cleanupDone
	log.Info("Goodbye!")
}
