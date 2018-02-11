package envsensor

import (
	"testing"
	"time"
)

func TestNewCachedReading(t *testing.T) {
	temp := float32(15.2)
	humd := float32(35.9)

	cr := NewCachedReading(Reading{Temperature: temp, Humidity: humd}, time.Second)
	if cr.Temperature != temp {
		t.Errorf("Temperature was incorrect, got %f, want %f", cr.Temperature, temp)
	}
	if cr.Humidity != humd {
		t.Errorf("Humidity was incorrect, got %f, want %f", cr.Humidity, temp)
	}

	if cr.IsStale() != false {
		t.Errorf("stale was incorrect, expected false got true")
	}
}

func TestIsStale(t *testing.T) {
	cr := NewCachedReading(Reading{Temperature: 5, Humidity: 5}, -1)

	if cr.IsStale() != true {
		t.Errorf("stale was incorrect, expected true got false")
	}
}

func TestIsNotStale(t *testing.T) {
	cr := NewCachedReading(Reading{Temperature: 5, Humidity: 5}, time.Second)

	if cr.IsStale() != false {
		t.Errorf("stale was incorrect, expected false got true")
	}
}
