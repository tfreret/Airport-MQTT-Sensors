package main

import (
	"airport/internal/apiClient"
	"airport/internal/config"
	"airport/internal/sensor"
	"fmt"
	"log"
	"time"
)

type PressureSensor struct {
	sensor.Sensor
}

func (pSensor *PressureSensor) GetActualizeMeasure() (sensor.Measurement, error) {
	apiResponse, err := apiClient.GetApiResponse(config.CHECKWX_URL+pSensor.Params.Airport+"/decoded", pSensor.Api.Key)
	if err != nil {
		log.Printf("Erreur lors de l'obtention de la réponse de l'API : %v", err)

		return sensor.Measurement{}, fmt.Errorf("échec lors de l'obtention de la mesure : %w", err)
	}

	if len(apiResponse.Data) == 0 {
		return sensor.Measurement{}, fmt.Errorf("réponse de l'API invalide")
	}
	return sensor.Measurement{TypeMesure: "Pres", Value: apiResponse.Data[0].Barometer.Hpa, Timestamp: time.Now().UTC().Format(time.RFC3339)}, nil
}

func NewPressureSensor(config sensor.ConfigSensor) *PressureSensor {
	pSensor := &PressureSensor{}
	pSensor.Sensor = sensor.NewSensor(pSensor, config)
	return pSensor
}
