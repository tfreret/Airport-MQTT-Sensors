package main

import (
	"airport/internal/mqttTools"
)

type ConfigStruct struct {
	Mqtt     mqttTools.ConfigMqtt `mapstructure:"mqtt"`
	MaxValue struct {
		MaxWindValue float64 `mapstructure:"maxWindValue"`
		MaxPresValue float64 `mapstructure:"maxPresValue"`
		MaxTempValue float64 `mapstructure:"maxTempValue"`
	} `mapstructure:"maxValue"`
}
