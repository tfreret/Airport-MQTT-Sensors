package sensor

type SensorInterface interface {
	Send(mesure Measurement)
	StartSendingData()
	GetActualizeMeasure() (Measurement, error)
}
