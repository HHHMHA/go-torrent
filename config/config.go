package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	DefaultDownloadLocation string `mapstructure:"defaultDownloadLocation"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("settings")

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	downloadsDirectory := homeDirectory + "/Downloads" // won't cause an issue for windows

	viper.SetDefault("DefaultDownloadLocation", downloadsDirectory)
	//viper.SetConfigType("json")

	//viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&config)
	return
}
