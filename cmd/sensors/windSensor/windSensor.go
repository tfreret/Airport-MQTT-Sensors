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

type WindSensor struct {
	sensor.Sensor
}

func (wSensor *WindSensor) GetActualizeMeasure() (sensor.Measurement, error) {
	apiResponse, err := apiClient.GetApiResponse(config.CHECKWX_URL+wSensor.Airport+"/decoded", config.CHECKWX_API_KEY)
	if err != nil {
		log.Printf("Erreur lors de l'obtention de la réponse de l'API : %v", err)

		return sensor.Measurement{}, fmt.Errorf("échec lors de l'obtention de la mesure : %w", err)
	}
	return sensor.Measurement{TypeMesure: "Wind", Value: apiResponse.Data[0].Wind.SpeedKph, Timestamp: time.Now().Format(time.RFC3339)}, nil
}

func NewWindSensor(idSensor int, idAirport string) *WindSensor {
	wSensor := &WindSensor{}
	wSensor.Sensor = sensor.NewSensor(wSensor, idSensor, idAirport)
	return wSensor
}
