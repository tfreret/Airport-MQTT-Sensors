package mqttTools

import (
	"airport/internal/config"
	"log"
	"github.com/eclipse/paho.mqtt.golang"
)

type BrokerClient struct {
	client mqtt.Client
}

func NewBrokerClient(idClient ...string) BrokerClient {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.BROKER_URL)
	opts.SetUsername(config.BROKER_USERNAME)
	opts.SetPassword(config.BROKER_PASSWORD)

	if len(idClient) == 1 {
		opts.SetClientID(idClient[0])
	}
		
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	return BrokerClient{client: client}
}

func (brokerClient BrokerClient) SendMessage(topic, message string) {
	token := brokerClient.client.Publish(topic, config.BROKER_QoS, false, message)
	token.Wait()
}


func (brokerClient BrokerClient) Subscribe(topic string, callBack func(topic string, message []byte)) {
	if token := brokerClient.client.Subscribe(topic, 0, 
		func(client mqtt.Client, msg mqtt.Message) {
			callBack(msg.Topic(), msg.Payload())
		}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}