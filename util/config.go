package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	TargetList string `mapstructure:"TARGET_LIST"`
	Port       int    `mapstructure:"PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.SetEnvPrefix("lb")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
