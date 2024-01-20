package sensor

import (
	"airport/internal/mqttTools"
)

type ConfigSensor struct {
	Mqtt   mqttTools.ConfigMqtt `mapstructure:"mqtt"`
	Params struct {
		Frequency int    `mapstructure:"frequency"`
		Airport   string `mapstructure:"airport"`
	} `mapstructure:"sensor"`
	Api struct {
		Use bool   `mapstructure:"use"`
		Key string `mapstructure:"key"`
	} `mapstructure:"api"`
}
