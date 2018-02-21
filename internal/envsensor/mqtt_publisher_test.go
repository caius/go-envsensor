package envsensor

import (
	"testing"
)

func TestMQTTPublisherClientId(t *testing.T) {
	pub := NewMQTTPublisher("mqtt.test:1883", "kitchen")

	if pub.clientId() != "envsensor_kitchen" {
		t.Fatal("clientID is not as expected")
	}
}
