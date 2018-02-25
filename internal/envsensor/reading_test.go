package envsensor

import (
	"testing"
)

func TestNewReading(t *testing.T) {
	temp := float32(15.52)
	humid := float32(39.2)
	sensor := "DHT22"

	r := NewReading(temp, humid, sensor)

	if r.Temperature != temp {
		t.Fatalf("Temperature is wrong. Expected %f, got %f", temp, r.Temperature)
	}

	if r.Humidity != humid {
		t.Fatalf("Humidity is wrong. Expected %f, got %f", humid, r.Humidity)
	}

	if r.SensorType != sensor {
		t.Fatalf("Sensor Type is wrong. Expected %s, got %s", sensor, r.SensorType)
	}

	if r.ReadAt.IsZero() {
		t.Fatal("Expected ReadAt not to be zero time")
	}
}
