package sensor

import (
	"airport/internal/mqttTools"
	"fmt"
	"time"
)

type Sensor struct {
	SensorInterface
	Id           int
	Airport      string
	brokerClient mqttTools.BrokerClient
}

func NewSensor(concreteSensor SensorInterface, idSensor int, idAirport string) Sensor {
	return Sensor{
		Id:              idSensor,
		Airport:         idAirport,
		SensorInterface: concreteSensor,
		brokerClient:    mqttTools.NewBrokerClient(),
	}
}

func (sensor Sensor) Send(mesure Measurement) {
	sensor.brokerClient.SendMessage(
		fmt.Sprintf("%s/%s/%d", sensor.Airport, mesure.TypeMesure, sensor.Id),
		fmt.Sprintf("%s;%f", mesure.Timestamp, mesure.Value),
	)
	fmt.Printf("%s/%s/%d\nvalue:%f\ntime:%s", sensor.Airport, mesure.TypeMesure, sensor.Id, mesure.Value, mesure.Timestamp)
}

func (sensor Sensor) StartSendingData(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
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
