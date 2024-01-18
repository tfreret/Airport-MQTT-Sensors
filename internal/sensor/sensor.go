package sensor

import (
	"airport/internal/mqttTools"
	"airport/internal/randomSensor"
	"fmt"
	"math"
	"time"
)

type Sensor struct {
	SensorInterface
	ConfigSensor
	brokerClient mqttTools.BrokerClient
	randomSensor.NumberGenerator
}

func NewSensor(concreteSensor SensorInterface, config ConfigSensor, generator randomSensor.NumberGenerator) Sensor {
	client := mqttTools.NewBrokerClient(
		config.Mqtt.MqttId,
		config.Mqtt.MqttUrl,
		config.Mqtt.MqttPort,
		config.Mqtt.MqttLogin,
		config.Mqtt.MqttPassword)
	return Sensor{
		ConfigSensor:    config,
		SensorInterface: concreteSensor,
		brokerClient:    client,
		NumberGenerator: generator,
	}
}

func (sensor Sensor) Send(mesure Measurement) {
	sensor.brokerClient.SendMessage(
		fmt.Sprintf("data/%s/%s/%s", sensor.Params.Airport, mesure.TypeMesure, sensor.Mqtt.MqttId),
		fmt.Sprintf("%s;%f", mesure.Timestamp, math.Round(mesure.Value*10)/10),
		sensor.Mqtt.MqttQOS,
	)
	fmt.Printf("data/%s/%s/%s\n value:%f\n time:%s\n", sensor.Params.Airport, mesure.TypeMesure, sensor.Mqtt.MqttId, math.Round(mesure.Value*10)/10, mesure.Timestamp)
}

func (sensor Sensor) StartSendingData() {
	ticker := time.NewTicker(time.Duration(sensor.Params.Frequency) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			measurement, err := sensor.GetActualizeMeasure()
			if err != nil {
				fmt.Printf("Erreur lors de l'obtention de la mesure : %v", err)
				return
			}
			sensor.Send(measurement)
		}
	}
}
