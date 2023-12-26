package main

import (
    "sensors/internal/common/sensor"
    "time"
)

type WindSensor struct {
    sensor.Sensor
}

func (wSensor WindSensor) GetActualizeMeasure() sensor.Measurement {
    // TODO fetch from api or get from json
    return sensor.Measurement{TypeMesure: "Temp", Value: 0.66, Timestamp: time.Now().Format(time.RFC3339)}
}

func NewWindSensor(idSensor int, idAirport string) *WindSensor {
    wSensor := WindSensor{ sensor.Sensor{Id: idSensor, Airport: idAirport} }
    wSensor.Sensor.SensorInterface = wSensor
    return &wSensor
}