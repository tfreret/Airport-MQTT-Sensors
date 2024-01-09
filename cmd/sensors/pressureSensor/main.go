package main

import (
	"airport/internal/sensor"
	"flag"
	"fmt"
)

func main() {
	var (
		configFile = flag.String("config", "sensor-pressure-config.env", "Config file of the sensor")
	)
	flag.Parse()
	config := sensor.ReadConfig(*configFile)

	fmt.Println("Using config :", config)
	NewPressureSensor(config).StartSendingData()
}
