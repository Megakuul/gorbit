package listener

import (
	"fmt"
	"github/megakuul/gorbit/conf"
	"github/megakuul/gorbit/handler"
	"github/megakuul/gorbit/logger"
	"net"
)

func Listen(config conf.Config) error {
	addr, err := net.ResolveTCPAddr("tcp",
		fmt.Sprintf(":%v", config.ListeningPort),
	)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	logger.WriteInformationLogger("Listening to port %v\n", config.ListeningPort)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.WriteWarningLogger(err)
		}

		go handler.HandleConnection(conn, config)
	}
}
