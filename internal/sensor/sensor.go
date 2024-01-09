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

func NewSensor(concreteSensor SensorInterface, config ConfigSensor) Sensor {
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
	}
}

func (sensor Sensor) Send(mesure Measurement) {
	sensor.brokerClient.SendMessage(
		fmt.Sprintf("%s/%s/%s", sensor.Params.Airport, mesure.TypeMesure, sensor.Mqtt.MqttId),
		fmt.Sprintf("value:%f\ntime:%s\n", mesure.Value, mesure.Timestamp),
		sensor.Mqtt.MqttQOS,
	)
	fmt.Printf("%s/%s/%s\n value:%f\n time:%s\n", sensor.Params.Airport, mesure.TypeMesure, sensor.Mqtt.MqttId, mesure.Value, mesure.Timestamp)
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

func ReadConfig(filename string) ConfigSensor {
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
