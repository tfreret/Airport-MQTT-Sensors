package main

import (
	"airport/internal/apiClient"
	"airport/internal/config"
	"airport/internal/sensor"
	"time"
)

type TempSensor struct {
	sensor.Sensor
}

func (tSensor *TempSensor) GetActualizeMeasure() sensor.Measurement {
	apiResponse, _ := apiClient.GetApiResponse(config.CHECKWX_URL+tSensor.Airport+"/decoded", config.CHECKWX_API_KEY)
	return sensor.Measurement{TypeMesure: "Temp", Value: apiResponse.Data[0].Temperature.Celsius, Timestamp: time.Now().Format(time.RFC3339)}
}

func NewTempSensor(idSensor int, idAirport string) *TempSensor {
	tSensor := &TempSensor{}
	tSensor.Sensor = sensor.NewSensor(tSensor, idSensor, idAirport)
	return tSensor
}
