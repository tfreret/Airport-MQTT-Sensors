package main

import (
    "sensors/internal/common/sensor"
    "time"
)

type TempSensor struct {
    sensor.Sensor
}

func (tSensor TempSensor) GetActualizeMeasure() sensor.Measurement {
    // TODO fetch from api or get from json
    return sensor.Measurement{TypeMesure: "Temp", Value: 0.66, Timestamp: time.Now().Format(time.RFC3339)}
}

func NewTempSensor(idSensor int, idAirport string) *TempSensor {
    tSensor := TempSensor{ sensor.Sensor{Id: idSensor, Airport: idAirport} }
    tSensor.Sensor.SensorInterface = tSensor
    return &tSensor
}