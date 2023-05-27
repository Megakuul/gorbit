package conf

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	ListeningPort        int    `mapstructure:"listeningport"`
	LogOptions           string `mapstructure:"logoptions"`
	BufferSizeKB         int    `mapstructure:"buffersizekb"`
	HealthCheckIntervall int    `mapstructure:"HealthCheckIntervall"`
	Endpoints            []Endpoint
}

type Endpoint struct {
	Port         int    `mapstructure:"port"`
	Hostname     string `mapstructure:"hostname"`
	Weight       int    `mapstructure:"weight"`
	Healthy      bool
	Sessions     int
	SessionMutex sync.RWMutex
}

func (e *Endpoint) MutAppendSession() {
	e.SessionMutex.Lock()
	defer e.SessionMutex.Unlock()
	e.Sessions++
}

func (e *Endpoint) MutRemoveSession() {
	e.SessionMutex.Lock()
	defer e.SessionMutex.Unlock()
	e.Sessions--
}

func LoadConfig() (Config, error) {
	var config Config

	viper.SetConfigName("gorbit-conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("buffersizekb", 32)
	viper.SetDefault("HealthCheckIntervall", 5)
	viper.SetDefault("logoptions", "ERROR|WARNING")

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
