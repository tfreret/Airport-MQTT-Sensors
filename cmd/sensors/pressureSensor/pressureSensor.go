package main

import (
	"airport/internal/apiClient"
	"airport/internal/config"
	"airport/internal/sensor"
	"time"
)

type PressureSensor struct {
	sensor.Sensor
}

func (pSensor *PressureSensor) GetActualizeMeasure() sensor.Measurement {
	apiResponse, _ := apiClient.GetApiResponse(config.CHECKWX_URL+pSensor.Airport+"/decoded", config.CHECKWX_API_KEY)
	return sensor.Measurement{TypeMesure: "Pres", Value: apiResponse.Data[0].Barometer.Hpa, Timestamp: time.Now().Format(time.RFC3339)}
}

func NewPressureSensor(idSensor int, idAirport string) *PressureSensor {
	pSensor := &PressureSensor{}
	pSensor.Sensor = sensor.NewSensor(pSensor, idSensor, idAirport)
	return pSensor
}
