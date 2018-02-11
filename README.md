# Envsensor

Environmental Sensor Server. Intended for Raspberry Pi's sampling temperature & humidity from the physical location using DHT11 or DHT22 (/AM2302) sensors and exposing the realtime data via HTTP.

There's two parts to the app, firstly it samples the sensor every so often (configurable) and keeps a record of the latest reading taken.

The second part is the webserver, which exposes the data as JSON at `/` and also for [Prometheus][] at `/metrics`.

[Prometheus]: https://prometheus.io/

## Usage

Install it on your RPi:

```shell
sudo apt install golang
GOPATH=$HOME/go go get -u github.com/caius/go-envsensor/cli/envsensor
```

Then you'll want to figure out which GPIO pin your sensor is connected to and which model of sensor you have (DHT22 is more accurate!)

```shell
# DHT22 with data connected to physical pin 12, BCM pin 18, GPIO.1
$HOME/go/bin/envsensor -sensor-pin 18 -sensor-version 22
```

See `init/systemd.txt` for an example Systemd unit file you can use to run it as a service (aka run on boot & restart if it crashes.)

## Architecture

The sensor is sampled in a goroutine, controlled via `DHTSensor`. (This is split across three files, so it can be developed on a mac without the DHT library installed.) This is kicked off by a `time.Ticker` channel every delay seconds. Any successful reading is sent into the readings channel.

The other half of the process is the webserver, which has a reading stored in a variable and both the endpoint handlers just return the data in that variable. It listens to the readings channel, any reading received is wrapped to a `CachedReading` which then expires a minute after it was created.

## License

Apache 2, see `LICENSE` for further details.
