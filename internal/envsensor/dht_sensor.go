package envsensor

import (
	"time"
)

type DHTSensor struct {
	Emitting     bool
	pin          int
	readingsChan chan<- Reading
	ticker       *time.Ticker
}

// Start emitting sensor readings into `readingsChan` channel.
//
// Only reads at most every `delay` seconds from the sensor
func (s *DHTSensor) Start(pin int, readingsChan chan<- Reading, delay int) {
	if s.Emitting == true {
		// Sensor is running already, stop it before we continue
		s.Stop()
	}

	s.pin = pin
	s.readingsChan = readingsChan
	s.ticker = time.NewTicker(time.Second * time.Duration(delay))
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
