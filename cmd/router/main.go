package main

import (
	"fmt"
	"os"
	"os/signal"

	"airport/internal/mqttTools"
)

func main() {
	brokerClient := mqttTools.NewBrokerClient("router")

	
	brokerClient.Subscribe("#", func(topic string, payload []byte){
		fmt.Printf("Received message on topic : '%s'\n%s", topic, string(payload))
	})
	
	// Attente de signaux d'arrêt
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	<-stopSignal

	// // Désabonnement et déconnexion
	// if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
	// 	log.Fatal(token.Error())
	// }
	// client.Disconnect(250)
}
