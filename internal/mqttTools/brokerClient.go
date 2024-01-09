package mqttTools

import (
	"airport/internal/config"
	"fmt"
	"os"
	"github.com/eclipse/paho.mqtt.golang"
)

type BrokerClient struct {
	client mqtt.Client
}

func NewBrokerClient() BrokerClient {
	broker := config.BROKER_URL

	opts := mqtt.NewClientOptions().AddBroker(broker)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Erreur de connexion au broker MQTT:", token.Error())
		os.Exit(1)
	}
	return BrokerClient{client: client}
}

func (brokerClient BrokerClient) SendMessage(topic, message string) {
	token := brokerClient.client.Publish(topic, config.BROKER_QoS, false, message)
	token.Wait()
}
