package sensor

import (
	"airport/internal/mqttTools"
	"airport/internal/randomSensor"
	"fmt"
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
		config.MqttId,
		config.MqttUrl,
		config.MqttPort,
		config.MqttLogin,
		config.MqttPassword)
	return Sensor{
		ConfigSensor:    config,
		SensorInterface: concreteSensor,
		brokerClient:    client,
		NumberGenerator: generator,
	}
}

func (sensor Sensor) Send(mesure Measurement) {
	sensor.brokerClient.SendMessage(
		fmt.Sprintf("data/%s/%s/%s", sensor.Airport, mesure.TypeMesure, sensor.MqttId),
		fmt.Sprintf("%s;%f", mesure.Timestamp, mesure.Value),
		sensor.MqttQOS,
	)
	fmt.Printf("data/%s/%s/%s\n value:%f\n time:%s\n", sensor.Airport, mesure.TypeMesure, sensor.MqttId, mesure.Value, mesure.Timestamp)
}

func (sensor Sensor) StartSendingData() {
	ticker := time.NewTicker(time.Duration(sensor.Frequency) * time.Second)
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