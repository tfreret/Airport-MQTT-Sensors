package main

import (
	"airport/internal/config"
	"airport/internal/sensor"
	"flag"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	var (
		configFile = flag.String("config", "sensor-pressure-config.yaml", "Config file of the sensor")
	)
	flag.Parse()
	configSensor := config.ReadConfig[sensor.ConfigSensor](*configFile)

	log.Println("Using config : ", configSensor)

	NewPressureSensor(configSensor).StartSendingData()
}
