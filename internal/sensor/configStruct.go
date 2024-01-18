package sensor

import(
	"airport/internal/mqttTools"
)

type ConfigUtilities struct {
	Frequency int    `mapstructure:"frequency"`
	Airport   string `mapstructure:"airport"`
}

type ConfigApi struct {
	Key string `mapstructure:"key"`
}

type ConfigSensor struct {
	mqttTools.ConfigMqtt		`mapstructure:"mqtt"`
	ConfigUtilities				`mapstructure:"sensor"`
	ConfigApi					`mapstructure:"api"`
}