package mqttTools

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type BrokerClient struct {
	client mqtt.Client
}

func NewBrokerClient(ConfigMqtt ConfigMqtt) BrokerClient {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

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

func (brokerClient BrokerClient) Subscribe(topic string, callBack func(topic string, message []byte), qos byte) {
	if token := brokerClient.client.Subscribe(topic, qos,
		func(client mqtt.Client, msg mqtt.Message) {
			callBack(msg.Topic(), msg.Payload())
		}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}
