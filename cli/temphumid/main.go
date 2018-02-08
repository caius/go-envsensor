package main

import (
	"fmt"
	"github.com/d2r2/go-dht"
	"github.com/urfave/cli"
	"os"
)

func printReadingFrom(pin int) {
	// TODO: make sensor configurable from args
	temperature, humidity, retried, err :=
		dht.ReadDHTxxWithRetry(dht.DHT11, pin, true, 10)
	if err != nil {
		panic(err)
	}
	// Print temperature and humidity
	fmt.Printf("Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
		temperature, humidity, retried)
}

func main() {
	var pin int

	app := cli.NewApp()
	app.Name = "envsensor"
	app.Usage = "Serve up environment readings from sensors"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "pin",
			Usage:       "GPIO pin (physical number) to read from",
			Destination: &pin,
		},
	}
	app.Action = func(c *cli.Context) error {
		if pin == 0 {
			panic("Pick a pin you dolt")
		}

		printReadingFrom(pin)

		return nil
	}

	app.Run(os.Args)
}
