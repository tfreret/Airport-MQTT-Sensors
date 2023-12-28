package main

import (
    "sensors/internal/sensor"
    "time"
)

type TempSensor struct {
    sensor.Sensor
}

func (tSensor TempSensor) GetActualizeMeasure() sensor.Measurement {
    // TODO fetch from api or get from json
    return sensor.Measurement{TypeMesure: "Temp", Value: 0.66, Timestamp: time.Now().Format(time.RFC3339)}
}

func NewTempSensor(idSensor int, idAirport string) (tSensor TempSensor) {
    tSensor.Sensor = sensor.NewSensor(tSensor, idSensor, idAirport)
    return tSensor
}