package main

import (
	"github/megakuul/gorbit/conf"
	"github/megakuul/gorbit/handler"
	"github/megakuul/gorbit/listener"
	"github/megakuul/gorbit/logger"
	"log"
	"os"
)

func main() {
	config, err := conf.LoadConfig()
	if err != nil {
		log.Fatalf("[ Gorbit Panic ]:\n%s\n", err)
		os.Exit(2)
	}

	err = logger.InitLogger("gorbit.log", config.LogOptions, config.MaxLogSizeKB)
	if err != nil {
		log.Fatalf("[ Gorbit Panic ]:\n%s\n", err)
		os.Exit(2)
	}

	go handler.CheckHealth(&config.Endpoints, config.HealthCheckIntervall)

	err = listener.Listen(config)
	if err != nil {
		logger.WriteErrLogger(err)
		os.Exit(2)
	}
}
