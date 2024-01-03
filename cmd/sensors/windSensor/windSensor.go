package main

import (
	"airport/internal/apiClient"
	"airport/internal/config"
	"airport/internal/sensor"
	"time"
)

type WindSensor struct {
	sensor.Sensor
}

func (wSensor *WindSensor) GetActualizeMeasure() sensor.Measurement {
	apiResponse, _ := apiClient.GetApiResponse(config.CHECKWX_URL+wSensor.Airport+"/decoded", config.CHECKWX_API_KEY)
	return sensor.Measurement{TypeMesure: "Wind", Value: apiResponse.Data[0].Wind.SpeedKph, Timestamp: time.Now().Format(time.RFC3339)}
}

func NewWindSensor(idSensor int, idAirport string) *WindSensor {
	wSensor := &WindSensor{}
	wSensor.Sensor = sensor.NewSensor(wSensor, idSensor, idAirport)
	return wSensor
}
