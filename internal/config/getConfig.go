package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var log = logrus.New()

func ReadConfig[Config any](filename string) Config {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Error("Error while loading config : ", err)
		os.Exit(1)
	}

	var config Config

	if err := viper.UnmarshalExact(&config); err != nil {
		log.Println("Error while parsing config : ", err)
		os.Exit(1)
	}
	return config
}

// We need to parse .env because it's easier docker usage
func ReadEnv[Config any](filename string) Config {
	viper.SetConfigName(filename)
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error while loading config : ", err)
		os.Exit(1)
	}

	var config Config

	if err := viper.UnmarshalExact(&config); err != nil {
		log.Println("Error while parsing config : ", err)
		os.Exit(1)
	}
	return config
}
