package mqttTools

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
)

type BrokerClient struct {
	client mqtt.Client
}

func NewBrokerClient(id string, url string, port int, login string, password string) BrokerClient {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprint(url, ":", port))

	opts.SetClientID(id)
	opts.SetUsername(login)
	opts.SetPassword(password)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	return BrokerClient{client: client}
}

func (brokerClient BrokerClient) SendMessage(topic, message string, qos byte) {
	token := brokerClient.client.Publish(topic, qos, false, message)
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
