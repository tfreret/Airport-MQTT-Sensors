package main

import (
	"airport/internal/apiClient"
	"airport/internal/config"
	"airport/internal/randomSensor"
	"airport/internal/sensor"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type PressureSensor struct {
	sensor.Sensor
}

func (pSensor *PressureSensor) GetActualizeMeasure() (sensor.Measurement, error) {
	if config.USE_API {
		apiResponse, err := apiClient.GetApiResponse(config.CHECKWX_URL+pSensor.Airport+"/decoded", pSensor.Key)
		if err != nil {
			log.Printf("Erreur lors de l'obtention de la réponse de l'API : %v", err)

			return sensor.Measurement{}, fmt.Errorf("échec lors de l'obtention de la mesure : %w", err)
		}

		if len(apiResponse.Data) == 0 {
			return sensor.Measurement{}, fmt.Errorf("réponse de l'API invalide")
		}
		return sensor.Measurement{TypeMesure: "Pres", Value: apiResponse.Data[0].Barometer.Hpa, Timestamp: time.Now().UTC().Format(time.RFC3339)}, nil
	} else {
		return sensor.Measurement{TypeMesure: "Pres", Value: pSensor.NumberGenerator.GenerateRandomNumber(), Timestamp: time.Now().UTC().Format(time.RFC3339)}, nil
	}
}

func NewPressureSensor(configSensor sensor.ConfigSensor) *PressureSensor {
	pSensor := &PressureSensor{}
	var min, max float64 = 950, 1100
	start := min + rand.Float64()*(max-min)
	nbGenerator := randomSensor.NewNumberGenerator(start, min, max)
	pSensor.Sensor = sensor.NewSensor(pSensor, configSensor, *nbGenerator)
	return pSensor
}
