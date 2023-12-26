package sensor

import (
	"fmt"
	"time"
)

type Sensor struct{
	SensorInterface
	Id int
	Airport string
}

func (sensor Sensor) Send(mesure Measurement) {
	// TODO use MQTT
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