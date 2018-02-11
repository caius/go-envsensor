package envsensor

// Override: fake sensor read for debugging on mac
//
// Obviously we can't attach a GPIO easily, so just fake out emitting every time
// we'd read a sensor instead.
func (s *DHTSensor) readSensor() Reading {
	return Reading{
		Temperature: 15.2,
		Humidity:    39.0,
	}
}
