package main

import(
	"airport/internal/mqttTools"
)

type ConfigMaxValue struct {
	MaxWindValue float64    `mapstructure:"maxWindValue"`
	MaxPresValue float64    `mapstructure:"maxPresValue"`
	MaxTempValue float64    `mapstructure:"maxTempValue"`
}

type ConfigStruct struct {
	Mqtt   mqttTools.ConfigMqtt		`mapstructure:"mqtt"`
	MaxValue ConfigMaxValue			`mapstructure:"maxValue"`
}