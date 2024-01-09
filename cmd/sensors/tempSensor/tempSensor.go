package main

import (
	"airport/internal/apiClient"
	"airport/internal/config"
	"airport/internal/sensor"
	"fmt"
	"log"
	"time"
)

type TempSensor struct {
	sensor.Sensor
}

func (tSensor *TempSensor) GetActualizeMeasure() (sensor.Measurement, error) {
	apiResponse, err := apiClient.GetApiResponse(config.CHECKWX_URL+tSensor.SensorAirport+"/decoded", config.CHECKWX_API_KEY)
	if err != nil {
		log.Printf("Erreur lors de l'obtention de la réponse de l'API : %v", err)

		return sensor.Measurement{}, fmt.Errorf("échec lors de l'obtention de la mesure : %w", err)
	}
	if len(apiResponse.Data) == 0 {
		return sensor.Measurement{}, fmt.Errorf("réponse de l'API invalide")
	}
	return sensor.Measurement{TypeMesure: "Temp", Value: apiResponse.Data[0].Temperature.Celsius, Timestamp: time.Now().UTC().Format(time.RFC3339)}, nil
}

func NewTempSensor(config sensor.ConfigSensor) *TempSensor {
	tSensor := &TempSensor{}
	tSensor.Sensor = sensor.NewSensor(tSensor, config)
	return tSensor
}
