package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServerCfg struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type LoggerCfg struct {
	LogLevel  string `yaml:"log_level"`
	LogOutput string `yaml:"log_output"`
}

type Config struct {
	Server ServerCfg `yaml:"server"`
	Logger LoggerCfg `yaml:"logger"`
}

var (
	cfg Config
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/configs")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error in configuration file: %s", err))
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
}

func GetConfig() Config {
	return cfg
}
