package envsensor

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type MQTTPublisher struct {
	Broker   string
	Location string
}

type MQTTReadingMessage struct {
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
}

func NewMQTTPublisher(broker string, location string) MQTTPublisher {
	return MQTTPublisher{
		Broker:   broker,
		Location: location,
	}
}

func (p *MQTTPublisher) subscribeToReadings(readings <-chan Reading, client mqtt.Client) {
	log.Debug("MQTTPublisher subscribing to readings")
	for reading := range readings {
		log.WithFields(log.Fields{
			"reading": reading,
		}).Info("MQTTPublisher received reading")

		msgReading := MQTTReadingMessage{
			Temperature: reading.Temperature,
			Humidity:    reading.Humidity,
		}

		topic := fmt.Sprintf("envsensor/status/%s", p.Location)

		payload, err := json.Marshal(msgReading)
		if err != nil {
			log.Fatal(err)
		} else {
			log.WithFields(log.Fields{
				"topic":   topic,
				"message": string(payload),
			}).Debug("MQTTPublisher publishing")
			client.Publish(topic, 0, false, string(payload)).Wait()
		}
	}
	log.Debug("MQTTPublisher finished listening for readings")
}

func (p *MQTTPublisher) Start(readings <-chan Reading) {
	log.WithFields(log.Fields{
		"broker":   p.Broker,
		"location": p.Location,
	}).Info("MQTTPublisher publishing")

	mqttParams := mqtt.NewClientOptions()
	mqttParams.AddBroker("tcp://mqtt1:1883")
	mqttParams.SetClientID("msub")

	client := mqtt.NewClient(mqttParams)
	client.Connect().Wait()

	log.WithFields(log.Fields{
		"broker": p.Broker,
	}).Info("MQTTPublisher connected to broker")

	p.subscribeToReadings(readings, client)
}
