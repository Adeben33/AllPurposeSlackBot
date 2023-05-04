package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	SlackAppToken     string `mapstructure:"SLACK_APP_TOKEN"`
	SlackAppChannelID string `mapstructure:"SLACK_APP_CHANNEL_ID"`
	AppEnv            string `mapstructure:"APP_ENV"`
	SlackBotToken     string `mapstructure:"SLACK_BOT_TOKEN"`
	ApplevelToken     string `mapstructure:"APP_LEVEL"`
}

var config *Config

func GetConfig() *Config {
	return config
}

func LoadConfig(path string) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return
		} else {
			// Config file was found but another error was produced
			return
		}
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)

	}
}
