package listener

import (
	"fmt"
	"github/megakuul/gorbit/conf"
	"github/megakuul/gorbit/logger"
	"net"
)

type HandleConnection func(net.Conn, conf.Config)

func Listen(config conf.Config, handleConnection HandleConnection) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", config.ListeningPort))
	if err != nil {
		return err
	}

	fmt.Printf("Listening to port %v\n", config.ListeningPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.WriteWarningLogger(err)
		}

		go handleConnection(conn, config)
	}
}
