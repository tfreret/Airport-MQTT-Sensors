package mqttTools

type ConfigMqtt struct {
	MqttUrl      string `mapstructure:"url"`
	MqttPort     int    `mapstructure:"port"`
	MqttQOS      byte   `mapstructure:"qos"`
	MqttId       string `mapstructure:"id"`
	MqttLogin    string `mapstructure:"login"`
	MqttPassword string `mapstructure:"password"`
}