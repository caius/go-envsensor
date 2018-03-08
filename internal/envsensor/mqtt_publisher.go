package envsensor

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"time"
)

type MQTTPublisher struct {
	Broker   string
	Location string
	client   mqtt.Client
}

type MQTTReadingMessage struct {
	Temperature float32   `json:"temperature"`
	Humidity    float32   `json:"humidity"`
	SensorType  string    `json:"sensor"`
	ReadAt      time.Time `json:"read_at"`
}

func NewMQTTPublisher(broker string, location string) MQTTPublisher {
	return MQTTPublisher{
		Broker:   broker,
		Location: location,
	}
}

func (p *MQTTPublisher) subscribeToReadings(readings <-chan Reading) {
	log.Debug("MQTTPublisher subscribing to readings")

	for reading := range readings {
		log.WithFields(log.Fields{
			"reading": reading,
		}).Info("MQTTPublisher received reading")

		msgReading := MQTTReadingMessage{
			Temperature: reading.Temperature,
			Humidity:    reading.Humidity,
			SensorType:  reading.SensorType,
			ReadAt:      reading.ReadAt,
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

			if p.client.IsConnected() {
				p.client.Publish(topic, 0, false, string(payload)).Wait()
			}
		}
	}

	log.Debug("MQTTPublisher finished listening for readings")
}

func (p *MQTTPublisher) clientId() string {
	return fmt.Sprintf("envsensor_%s", p.Location)
}

func (p *MQTTPublisher) Start(readings <-chan Reading) {
	log.WithFields(log.Fields{
		"broker":   p.Broker,
		"location": p.Location,
	}).Info("MQTTPublisher publishing")

	mqttParams := mqtt.NewClientOptions()
	mqttParams.AddBroker(fmt.Sprintf("tcp://%s", p.Broker))
	mqttParams.SetClientID(p.clientId())

	p.client = mqtt.NewClient(mqttParams)

	p.client.Connect().Wait()
	log.WithFields(log.Fields{
		"broker": p.Broker,
	}).Info("MQTTPublisher connected to broker")

	// Do our job
	p.subscribeToReadings(readings)
}

func (p *MQTTPublisher) Stop() {
	log.Info("MQTTPublisher received stop")
	p.client.Disconnect(250)
}
