package mqttTools

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type BrokerClient struct {
	client mqtt.Client
}

func NewBrokerClient(ConfigMqtt ConfigMqtt) BrokerClient {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprint(ConfigMqtt.MqttUrl, ":", ConfigMqtt.MqttPort))

	opts.SetClientID(ConfigMqtt.MqttId)
	opts.SetUsername(ConfigMqtt.MqttLogin)
	opts.SetPassword(ConfigMqtt.MqttPassword)
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

func (brokerClient BrokerClient) Subscribe(topic string, qos byte, callBack func(topic string, message []byte)) {
	if token := brokerClient.client.Subscribe(topic, qos,
		func(client mqtt.Client, msg mqtt.Message) {
			callBack(msg.Topic(), msg.Payload())
		}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}
