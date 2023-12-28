package sensor

import (
	"fmt"
	"time"
	"sensors/internal/mqttTools"
)

type Sensor struct{
	SensorInterface
	Id int
	Airport string
	brokerClient mqttTools.BrokerClient
}

func NewSensor(concreteSensor SensorInterface, idSensor int, idAirport string) Sensor {
    return Sensor{
		Id : idSensor,
		Airport : idAirport,
		SensorInterface : concreteSensor,
		brokerClient: mqttTools.NewBrokerClient(),
	}
}

func (sensor Sensor) Send(mesure Measurement) {
	sensor.brokerClient.SendMessage(
		fmt.Sprintf("%s/%s/%d", sensor.Airport, mesure.TypeMesure, sensor.Id), 
		fmt.Sprintf("value:%f\ntime:%s\n", mesure.Value, mesure.Timestamp), 
	)
	fmt.Printf("%s/%s/%d\n value:%f\n time:%s\n", sensor.Airport, mesure.TypeMesure, sensor.Id, mesure.Value, mesure.Timestamp)
}

func (sensor Sensor) StartSendingData(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sensor.Send(sensor.GetActualizeMeasure())
		}
	}
}