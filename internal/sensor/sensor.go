package sensor

import (
	"airport/internal/mqttTools"
	"airport/internal/randomSensor"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Sensor struct {
	SensorInterface
	ConfigSensor
	brokerClient mqttTools.BrokerClient
	randomSensor.NumberGenerator
}

type ConfigMqtt struct {
	MqttUrl      string `mapstructure:"url"`
	MqttPort     int    `mapstructure:"port"`
	MqttQOS      byte   `mapstructure:"qos"`
	MqttId       string `mapstructure:"id"`
	MqttLogin    string `mapstructure:"login"`
	MqttPassword string `mapstructure:"password"`
}

type ConfigApi struct {
	Key string `mapstructure:"key"`
}

type ConfigUtilities struct {
	Frequency int    `mapstructure:"frequency"`
	Airport   string `mapstructure:"airport"`
}

type ConfigSensor struct {
	Mqtt   ConfigMqtt      `mapstructure:"mqtt"`
	Params ConfigUtilities `mapstructure:"sensor"`
	Api    ConfigApi       `mapstructure:"api"`
}

func NewSensor(concreteSensor SensorInterface, config ConfigSensor, generator randomSensor.NumberGenerator) Sensor {
	client := mqttTools.NewBrokerClient(
		config.Mqtt.MqttId,
		config.Mqtt.MqttUrl,
		config.Mqtt.MqttPort,
		config.Mqtt.MqttLogin,
		config.Mqtt.MqttPassword)
	return Sensor{
		ConfigSensor:    config,
		SensorInterface: concreteSensor,
		brokerClient:    client,
		NumberGenerator: generator,
	}
}

func (sensor Sensor) Send(mesure Measurement) {
	sensor.brokerClient.SendMessage(
		fmt.Sprintf("data/%s/%s/%s", sensor.Params.Airport, mesure.TypeMesure, sensor.Mqtt.MqttId),
		fmt.Sprintf("%s;%f\n", mesure.Timestamp, mesure.Value),
		sensor.Mqtt.MqttQOS,
	)
	fmt.Printf("data/%s/%s/%s\n value:%f\n time:%s\n", sensor.Params.Airport, mesure.TypeMesure, sensor.Mqtt.MqttId, mesure.Value, mesure.Timestamp)
}

func (sensor Sensor) StartSendingData() {
	ticker := time.NewTicker(time.Duration(sensor.Params.Frequency) * time.Second)
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

func ReadSensorConfig(filename string) ConfigSensor {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
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
