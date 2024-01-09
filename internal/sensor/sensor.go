package sensor

import (
	"airport/internal/mqttTools"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

type Sensor struct {
	SensorInterface
	ConfigSensor
	brokerClient mqttTools.BrokerClient
}

type ConfigSensor struct {
	MqttUrl         string `mapstructure:"MQTT_URL"`
	MqttPort        int    `mapstructure:"MQTT_PORT"`
	MqttQOS         byte   `mapstructure:"MQTT_QOS"`
	MqttId          string `mapstructure:"MQTT_ID"`
	MqttLogin       string `mapstructure:"MQTT_LOGIN"`
	MqttPassword    string `mapstructure:"MQTT_PASSWORD"`
	SensorAirport   string `mapstructure:"SENSOR_AIRPORT"`
	SensorFrequency int    `mapstructure:"SENSOR_FREQUENCY"`
	ApiKey          string `mapstructure:"API_KEY"`
}

func NewSensor(concreteSensor SensorInterface, config ConfigSensor) Sensor {
	client := mqttTools.NewBrokerClient(
		config.MqttId,
		config.MqttUrl,
		config.MqttPort,
		config.MqttLogin,
		config.MqttPassword)
	return Sensor{
		ConfigSensor:    config,
		SensorInterface: concreteSensor,
		brokerClient:    client,
	}
}

func (sensor Sensor) Send(mesure Measurement) {
	sensor.brokerClient.SendMessage(
		fmt.Sprintf("%s/%s/%s", sensor.SensorAirport, mesure.TypeMesure, sensor.MqttId),
		fmt.Sprintf("value:%f\ntime:%s\n", mesure.Value, mesure.Timestamp),
		sensor.MqttQOS,
	)
	fmt.Printf("%s/%s/%s\n value:%f\n time:%s\n", sensor.SensorAirport, mesure.TypeMesure, sensor.MqttId, mesure.Value, mesure.Timestamp)
}

func (sensor Sensor) StartSendingData() {
	ticker := time.NewTicker(time.Duration(sensor.SensorFrequency) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			measurement, err := sensor.GetActualizeMeasure()
			if err != nil {
				fmt.Printf("Erreur lors de l'obtention de la mesure : %v", err)
				return
			}
			sensor.Send(measurement)
		}
	}
}

func ReadConfig(filename string) ConfigSensor {
	viper.SetConfigName(filename)
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error while loading config :\n", err)
		os.Exit(1)
	}

	var config ConfigSensor

	if err := viper.UnmarshalExact(&config); err != nil {
		fmt.Println("Error while parsing config :\n", err)
		os.Exit(1)
	}
	return config
}
