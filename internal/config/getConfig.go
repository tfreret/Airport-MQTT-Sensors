package config

import(
	"fmt"
	"os"
	"github.com/spf13/viper"
)

func ReadConfig[Config any](filename string) Config {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error while loading config :\n", err)
		os.Exit(1)
	}

	var config Config

	if err := viper.UnmarshalExact(&config); err != nil {
		fmt.Println("Error while parsing config :\n", err)
		os.Exit(1)
	}
	return config
}