package main

import (
	"airport/internal/sensor"
	"flag"
	"fmt"
)

func main() {
	var (
		configFile = flag.String("config", "sensor-pressure-config.yaml", "Config file of the sensor")
	)
	flag.Parse()
	config := sensor.ReadConfig(*configFile)

	fmt.Println("Using config :", config)
	NewPressureSensor(config).StartSendingData()
}
