package sensor

type ConfigMqtt struct {
	MqttUrl      string `mapstructure:"url"`
	MqttPort     int    `mapstructure:"port"`
	MqttQOS      byte   `mapstructure:"qos"`
	MqttId       string `mapstructure:"id"`
	MqttLogin    string `mapstructure:"login"`
	MqttPassword string `mapstructure:"password"`
}

type ConfigUtilities struct {
	Frequency int    `mapstructure:"frequency"`
	Airport   string `mapstructure:"airport"`
}

type ConfigApi struct {
	Key string `mapstructure:"key"`
}

type ConfigSensor struct {
	Mqtt   ConfigMqtt      `mapstructure:"mqtt"`
	Params ConfigUtilities `mapstructure:"sensor"`
	Api    ConfigApi       `mapstructure:"api"`
}