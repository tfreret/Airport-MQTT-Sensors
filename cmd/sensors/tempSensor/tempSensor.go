package main

import (
	"airport/internal/apiClient"
	"airport/internal/sensors/config"
	"airport/internal/sensors/sensor"
	"fmt"
	"log"
	"time"
)

type TempSensor struct {
	sensor.Sensor
}

func (tSensor *TempSensor) GetActualizeMeasure() (sensor.Measurement, error) {
	apiResponse, err := apiClient.GetApiResponse(config.CHECKWX_URL+tSensor.Airport+"/decoded", config.CHECKWX_API_KEY)
	if err != nil {
		log.Printf("Erreur lors de l'obtention de la réponse de l'API : %v", err)

		return sensor.Measurement{}, fmt.Errorf("échec lors de l'obtention de la mesure : %w", err)
	}
	return sensor.Measurement{TypeMesure: "Temp", Value: apiResponse.Data[0].Temperature.Celsius, Timestamp: time.Now().Format(time.RFC3339)}, nil
}

func NewTempSensor(idSensor int, idAirport string) *TempSensor {
	tSensor := &TempSensor{}
	tSensor.Sensor = sensor.NewSensor(tSensor, idSensor, idAirport)
	return tSensor
}
