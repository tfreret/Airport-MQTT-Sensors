package sensor

type SensorInterface interface{
	Send(mesure Measurement)
	StartSendingData(interval int)
	GetActualizeMeasure() Measurement
}