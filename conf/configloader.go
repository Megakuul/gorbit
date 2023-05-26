package conf

import (
	"github.com/spf13/viper"
)

type Config struct {
	ListeningPort int    `mapstructure:"listeningport"`
	LogOptions    string `mapstructure:"logoptions"`
	Endpoints     []Endpoint
}

type Endpoint struct {
	Port       int    `mapstructure:"port"`
	Hostname   string `mapstructure:"hostname"`
	Weight     int    `mapstructure:"weight"`
	Timeout_ms int    `mapstructure:"Timeout_ms"`
}

func LoadConfig() (Config, error) {
	var config Config

	viper.SetConfigName("gorbit-conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
