package envsensor

import (
	"fmt"
)

// Override: fake sensor read for debugging on mac
//
// Obviously we can't attach a GPIO easily, so just fake out emitting every time
// we'd read a sensor instead.
func (s *DHTSensor) readSensor() (Reading, error) {
	sensor := fmt.Sprintf("DHT%d", s.Version)
	reading := NewReading(float32(15.2), float32(39.0), sensor)
	return reading, nil
}
