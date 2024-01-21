package sensor

import (
	"airport/internal/mqttTools"
	"airport/internal/randomSensor"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Sensor struct {
	SensorInterface
	ConfigSensor
	brokerClient mqttTools.BrokerClient
	randomSensor.NumberGenerator
}

func NewSensor(concreteSensor SensorInterface, config ConfigSensor, generator randomSensor.NumberGenerator) Sensor {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	client := mqttTools.NewBrokerClient(
		config.Mqtt)
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
		fmt.Sprintf("%s;%s", mesure.Timestamp, strconv.FormatFloat(math.Round(mesure.Value*10/10), 'f', -2, 64)),
		sensor.Mqtt.MqttQOS,
	)
	log.Printf("Send message on topic : '/data/%s/%s/%s' value: '%s:%s'", sensor.Params.Airport, mesure.TypeMesure, sensor.Mqtt.MqttId, mesure.Timestamp, strconv.FormatFloat(math.Round(mesure.Value*10/10), 'f', -2, 64))
}

func (sensor Sensor) StartSendingData() {
	ticker := time.NewTicker(time.Duration(sensor.Params.Frequency) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			measurement, err := sensor.GetActualizeMeasure()
			if err != nil {
				log.Errorf("Erreur lors de l'obtention de la mesure : %v", err)
				return
			}
			sensor.Send(measurement)
		}
	}
}
