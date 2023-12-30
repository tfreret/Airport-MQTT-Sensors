package main

import (
    "flag"
)

func main() {
    var (
		idAirport  = flag.String("airport", "", "Airport id of this sensor")
		idSensor   = flag.Int("id", 0, "Id of this sensor")
        frequency  = flag.Int("frequency", 30, "Time frequence of sending data updates in second")
	)

	flag.Parse()

    NewPressureSensor(*idSensor, *idAirport).StartSendingData(*frequency)
}