package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServerCfg struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type LoggerCfg struct {
	LogLevel  string `mapstructure:"log_level"`
	LogOutput string `mapstructure:"log_output"`
}

type DatabaseCfg struct {
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
	LogLevel string `mapstructure:"log_level"`
}

type MovieApiCfg struct {
	ApiKey string `mapstructure:"api_key"`
}

type Config struct {
	Server   ServerCfg   `mapstructure:"server"`
	Logger   LoggerCfg   `mapstructure:"logger"`
	Database DatabaseCfg `mapstructure:"database"`
	MovieApi MovieApiCfg `mapstructure:"movieapi"`
}

var (
	cfg *Config
)

const (
	configName = "config"
	configType = "yaml"
	configPath = "$HOME/configs"
)

func init() {
	var err error

	cfg, err = setupConfig(configName, configType, configPath)
	if err != nil {
		panic(fmt.Errorf("fatal error in configuration file: %s", err))
	}
}

func setupConfig(cname, ctype, cpath string) (*Config, error) {
	var c *Config

	viper.SetConfigName(cname)
	viper.SetConfigType(ctype)
	viper.AddConfigPath(cpath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func GetConfig() *Config {
	return cfg
}
