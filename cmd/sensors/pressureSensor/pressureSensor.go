package main

import (
	"airport/internal/apiClient"
<<<<<<< HEAD
	"airport/internal/sensors/config"
	"airport/internal/sensors/sensor"
	"fmt"
	"log"
=======
	"airport/internal/config"
	"airport/internal/sensor"
>>>>>>> main
	"time"
)

type PressureSensor struct {
	sensor.Sensor
}

func (pSensor *PressureSensor) GetActualizeMeasure() (sensor.Measurement, error) {
	apiResponse, err := apiClient.GetApiResponse(config.CHECKWX_URL+pSensor.Airport+"/decoded", config.CHECKWX_API_KEY)
	if err != nil {
		log.Printf("Erreur lors de l'obtention de la réponse de l'API : %v", err)

		return sensor.Measurement{}, fmt.Errorf("échec lors de l'obtention de la mesure : %w", err)
	}
	return sensor.Measurement{TypeMesure: "Pres", Value: apiResponse.Data[0].Barometer.Hpa, Timestamp: time.Now().Format(time.RFC3339)}, nil
}

func NewPressureSensor(idSensor int, idAirport string) *PressureSensor {
	pSensor := &PressureSensor{}
	pSensor.Sensor = sensor.NewSensor(pSensor, idSensor, idAirport)
	return pSensor
}
