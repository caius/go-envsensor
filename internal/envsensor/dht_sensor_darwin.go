package envsensor

// Debug mode for mac. Obviously we can't attach a GPIO easily, so
// just fake out emitting every time we'd read a sensor instead
func (s *DHTSensor) readAndEmit() {
	for _ = range s.ticker.C {
		s.readingsChan <- Reading{
			Temperature: 15.2,
			Humidity:    39.0,
		}
	}
}
