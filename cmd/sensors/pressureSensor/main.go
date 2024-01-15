package main

import (
	"airport/internal/sensor"
	"airport/internal/config"
	"flag"
	"fmt"
)

func main() {
	var (
		configFile = flag.String("config", "sensor-pressure-config.yaml", "Config file of the sensor")
	)
	flag.Parse()
	config := config.ReadSensorConfig[sensor.ConfigSensor](*configFile)

	fmt.Println("Using config :", config)
	NewPressureSensor(config).StartSendingData()
}
