package envsensor

import (
	"fmt"
	"os"
	"testing"
)

func TestMQTTPublisherClientId(t *testing.T) {
	pub := NewMQTTPublisher("mqtt.test:1883", "kitchen")

	hostname, err := os.Hostname()
	if err != nil {
		t.Fatal("hostname failed")
	}

	if hostname == "" {
		t.Fatal("hostname is empty")
	}

	if pub.clientId() != fmt.Sprintf("envsensor_%s", hostname) {
		t.Fatal("clientID is not as expected")
	}
}
