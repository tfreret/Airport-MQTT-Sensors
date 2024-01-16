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
	Mqtt   mqttTools.ConfigMqtt		`mapstructure:"mqtt"`
	Params ConfigUtilities			`mapstructure:"sensor"`
	Api    ConfigApi				`mapstructure:"api"`
}