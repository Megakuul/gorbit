package conf

import (
	"github.com/spf13/viper"
)

type Config struct {
	ListeningPort int `mapstructure:"listeningport"`
	Endpoints     []Endpoint
}

type Endpoint struct {
	Port     int    `mapstructure:"port"`
	Hostname string `mapstructure:"hostname"`
	Weight   int    `mapstructure:"weight"`
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
