package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"strconv"

	"airport/internal/mqttTools"
)

func main() {
	const MAX_WIND_VALUE = 10.0
	const MAX_TEMP_VALUE = 6.0
	const MAX_PRES_VALUE = 1025.0

	brokerClient := mqttTools.NewBrokerClient("alertManager")

	
	brokerClient.Subscribe("#", func(topic string, payload []byte){
		alert := false

		topicElements := strings.Split(topic, "/")
		msgElements := strings.Split(string(payload), ";")
		if (len(topicElements) >= 2 && len(strings.Split(string(payload), ";")) >= 2){

			value, _ := strconv.ParseFloat(msgElements[1], 64)

			switch topicElements[1] {
			case "Temp":
				alert = value > MAX_TEMP_VALUE
			case "Pres":
				alert = value > MAX_PRES_VALUE
			case "Wind":
				alert = value > MAX_WIND_VALUE
			}
		}

		if (alert){
			fmt.Printf("Alerte for topic %s, value : %s \n", topic, string(payload))
			brokerClient.SendMessage(fmt.Sprintf("alert/%s", topic), string(payload))
		}
	})
	
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	<-stopSignal
}
