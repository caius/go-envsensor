package envsensor

import (
	"github.com/d2r2/go-dht"
	log "github.com/sirupsen/logrus"
)

// Override: talk to sensor on RPi to get actual reading
func (s *DHTSensor) readSensor() (Reading, error) {
	log.Debug("Reading sensor on ARM")

	var sensorKind dht.SensorType
	switch s.Version {
	case 11:
		sensorKind = dht.DHT11
	case 22:
		sensorKind = dht.DHT22
	default:
		log.WithFields(log.Fields{
			"version": s.Version,
		}).Fatal("Unknown sensor version. Needs to be 11 or 22. Assuming 11.")
		sensorKind = dht.DHT11
	}

	temperature, humidity, retried, err :=
		dht.ReadDHTxxWithRetry(sensorKind, s.Pin, true, 10)

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
		SensorType:  fmt.Sprintf("DHT%d", s.Version),
	}
	return reading, nil
}
