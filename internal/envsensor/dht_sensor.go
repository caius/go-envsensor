package envsensor

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type DHTSensor struct {
	Emitting     bool
	Pin          int
	Delay        time.Duration
	Version      int
	readingChans []chan Reading
	ticker       *time.Ticker
}

// Public: creates a configured DHTSensor instance
//
func NewDHTSensor(version int, pin int, delay time.Duration) DHTSensor {
	log.WithFields(log.Fields{
		"version": version,
		"pin":     pin,
		"delay":   delay,
	}).Debug("Creating DHTSensor")

	return DHTSensor{
		Version: version,
		Pin:     pin,
		Delay:   delay,
	}
}

// Internal: call readSensor() every Delay seconds and emit reading to readingChans
func (s *DHTSensor) readAndEmit() {
	// Take/Publish first reading *now*
	reading, err := s.readSensor()
	if err != nil {
		log.Error("Error reading sensor")
	} else {
		for _, c := range s.readingChans {
			c <- reading
		}
	}

	// And then continue reading in future
	for _ = range s.ticker.C {
		reading, err := s.readSensor()
		if err != nil {
			log.Error("Error reading sensor")
		} else {
			for _, c := range s.readingChans {
				c <- reading
			}
		}
	}
}

// Start emitting sensor readings into `readingChans` channels.
//
// Only reads at most every `Delay` seconds from the sensor
func (s *DHTSensor) Start(readingChans []chan Reading) {
	if s.Emitting == true {
		// Sensor is running already, stop it before we continue
		s.Stop()
	}

	s.readingChans = readingChans
	s.ticker = time.NewTicker(s.Delay)
	s.Emitting = true

	// Read and emit
	s.readAndEmit()
}

// Stops a running sensor from reading/emitting readingChans
func (s *DHTSensor) Stop() {
	if s.Emitting {
		s.ticker.Stop()
		for _, c := range s.readingChans {
			close(c)
		}
	}
	s.Emitting = false
}
