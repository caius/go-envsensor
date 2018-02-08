package envsensor

import (
	"time"
)

type DHTSensor struct {
	Pin      int
	Emitting bool
	ticker   *time.Ticker
	results  chan<- Reading
}

// Emits sensor readings into results channel.
//
// Only reads at most every `delay` seconds from the sensor
//
// This is faked out to work on my mac with no sensor attached
func (s *DHTSensor) EmitTo(results chan<- Reading, delay int) {
	if s.Emitting == true {
		// Sensor is running already, stop it before we continue
		s.Stop()
	}

	s.ticker = time.NewTicker(time.Second * time.Duration(delay))
	s.results = results
	go func() {
		for _ = range s.ticker.C {
			println("received tick, emitting reading")
			s.results <- Reading{
				Temperature: 15.2,
				Humidity:    39.0,
			}
		}
	}()
	s.Emitting = true
}

func (s *DHTSensor) Stop() {
	if s.Emitting {
		s.ticker.Stop()
		close(s.results)
	}
	s.Emitting = false
}
