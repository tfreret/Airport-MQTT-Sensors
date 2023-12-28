package main

import (
    "sensors/internal/sensor"
    "time"
)

type PressureSensor struct {
    sensor.Sensor
}

func (pSensor PressureSensor) GetActualizeMeasure() sensor.Measurement {
    // TODO fetch from api or get from json
    return sensor.Measurement{TypeMesure: "Pres", Value: 0.66, Timestamp: time.Now().Format(time.RFC3339)}
}

func NewPressureSensor(idSensor int, idAirport string) (pSensor PressureSensor) {
    pSensor.Sensor = sensor.NewSensor(pSensor, idSensor, idAirport)
    return pSensor
}