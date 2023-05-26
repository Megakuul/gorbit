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

	err = logger.InitLogger("gorbit.log", config.LogOptions)
	if err != nil {
		log.Fatalf("[ Gorbit Panic ]:\n%s\n", err)
		os.Exit(2)
	}

	loadBalancer := &handler.LoadBalancer{
		Sessions: make(map[string]*handler.Session),
	}

	err = listener.Listen(config, loadBalancer.HandleConnection)
	if err != nil {
		logger.WriteErrLogger(err)
		os.Exit(2)
	}
}
