package main

import (
	"fmt"
	"github.com/caius/go-envsensor/internal/envsensor"
	"time"
)

func main() {
	fmt.Println("Welcome to temphumid, where it is our pleasure to probe you today.")

	readings := make(chan envsensor.Reading)

	// Grab the readings every second
	sensor := envsensor.DHTSensor{}
	sensor.Start(17, readings, 1)

	go func() {
		for r := range readings {
			fmt.Printf("Got reading! t=%f, h=%f\n", r.Temperature, r.Humidity)
		}
	}()

	// Fake waiting for stuff
	time.Sleep(time.Second * time.Duration(2))
	sensor.Stop()
}
