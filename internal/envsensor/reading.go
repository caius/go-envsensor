package envsensor

import (
	"time"
)

type Reading struct {
	Temperature float32
	Humidity    float32
	SensorType  string
	ReadAt      time.Time
}

func NewReading(temp float32, humid float32, sensor string) Reading {
	return Reading{
		Temperature: temp,
		Humidity:    humid,
		SensorType:  sensor,
		ReadAt:      time.Now(),
	}
}
