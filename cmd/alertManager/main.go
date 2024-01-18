package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"strconv"
	"flag"

	"airport/internal/mqttTools"
	"airport/internal/config"
)

func main() {
	var (
		configFile = flag.String("config", "alert-manager-config.yaml", "Config file of the sensor")
	)
	flag.Parse()

	config := config.ReadConfig[ConfigStruct](*configFile)

	brokerClient := mqttTools.NewBrokerClient(
		config.MqttId,
		config.MqttUrl,
		config.MqttPort,
		config.MqttLogin,
		config.MqttPassword,
	)
	
	brokerClient.Subscribe("data/#", func(topic string, payload []byte){

		alert := false

		topicElements := strings.Split(topic, "/")
		msgElements := strings.Split(string(payload), ";")
		if (len(topicElements) >= 3 && len(msgElements) >= 2){
			if value, err := strconv.ParseFloat(msgElements[1], 64); err == nil {
				switch topicElements[2] {
				case "Temp":
					alert = value > config.MaxTempValue
				case "Pres":
					alert = value > config.MaxPresValue
				case "Wind":
					alert = value > config.MaxWindValue
				}
			}
		}

		if (alert){
			fmt.Printf("Alerte for topic %s, value : %s \n", topic, string(payload))
			brokerClient.SendMessage(fmt.Sprintf("alert/%s", topic), string(payload), config.MqttQOS)
		}
	})
	
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	<-stopSignal
}
