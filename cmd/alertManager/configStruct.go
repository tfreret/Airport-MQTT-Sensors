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
	mqttTools.ConfigMqtt		`mapstructure:"mqtt"`
	ConfigMaxValue				`mapstructure:"maxValue"`
}