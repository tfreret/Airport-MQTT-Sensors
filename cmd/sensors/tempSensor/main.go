package main

import (
	"airport/internal/sensor"
	"flag"
	"fmt"
)

func main() {
	var (
		configFile = flag.String("config", "sensor-temperature-config.yaml", "Config file of the sensor")
	)
	flag.Parse()
	config := sensor.ReadSensorConfig(*configFile)

	fmt.Println("Using config :", config)
	NewTempSensor(config).StartSendingData()
}
