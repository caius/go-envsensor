package envsensor

import (
	"github.com/d2r2/go-dht"
	log "github.com/sirupsen/logrus"
)

// Override: talk to sensor on RPi to get actual reading
func (s *DHTSensor) readSensor() (Reading, error) {
	log.Debug("Reading sensor on ARM")

	temperature, humidity, retried, err :=
		dht.ReadDHTxxWithRetry(dht.DHT11, s.Pin, true, 10)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Got error reading from sensor")
		return Reading{}, err
	}

	log.WithFields(log.Fields{
		"temperature": temperature,
		"humidity":    humidity,
		"retried":     retried,
	}).Debug("Read sensor successfully")

	reading := Reading{
		Temperature: temperature,
		Humidity:    humidity,
	}
	return reading, nil
}
