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
	readingsChan chan<- Reading
	ticker       *time.Ticker
}

func NewDHTSensor(version int, pin int, delay time.Duration) DHTSensor {
	log.Debug("Creating DHT Sensor with version=%d, pin=%d, delay=%s", version, pin, delay)
	return DHTSensor{
		Version: version,
		Pin:     pin,
		Delay:   delay,
	}
}

// Internal: call readSensor() every Delay seconds and emit reading to readingsChan
func (s *DHTSensor) readAndEmit() {
	reading, err := s.readSensor()
	if err != nil {
		log.Error("Error reading sensor")
	} else {
		// Take first reading *now*
		s.readingsChan <- reading
	}

	// And then continue reading in future
	for _ = range s.ticker.C {
		reading, err := s.readSensor()
		if err != nil {
			log.Error("Error reading sensor")
		} else {
			// Take first reading *now*
			s.readingsChan <- reading
		}
	}
}

// Start emitting sensor readings into `readingsChan` channel.
//
// Only reads at most every `Delay` seconds from the sensor
func (s *DHTSensor) Start(readingsChan chan<- Reading) {
	if s.Emitting == true {
		// Sensor is running already, stop it before we continue
		s.Stop()
	}

	s.readingsChan = readingsChan
	s.ticker = time.NewTicker(s.Delay)
	s.Emitting = true

	// Read and emit
	s.readAndEmit()
}

// Stops a running sensor from reading/emitting readingsChan
func (s *DHTSensor) Stop() {
	if s.Emitting {
		s.ticker.Stop()
		close(s.readingsChan)
	}
	s.Emitting = false
}
