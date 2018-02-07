package main

import "fmt"
import "github.com/morus12/dht22"
import "log"

func main() {
	sensor := dht22.New("GPIO_17")
	temperature, err := sensor.Temperature()
	if err != nil {
		log.Fatal(err)
	}
	humidity, err := sensor.Humidity()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Temperature %f\n", temperature)
	fmt.Printf("Humidity %f\n", humidity)
}
