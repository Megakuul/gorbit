package main

import (
	"fmt"
	"github/megakuul/gorbit/conf"
	"io"
	"net"
	"os"
)

func main() {
	config, err := conf.LoadConfig()
	if err != nil {
		fmt.Printf("Gorbit Panic:\n%s\n", err)
		os.Exit(2)
	}

	fmt.Printf("ListeningPort: %d\n", config.ListeningPort)
	for i, endpoint := range config.Endpoints {
		fmt.Printf("Endpoint %d, Port: %d, Hostname: %s\n", i, endpoint.Port, endpoint.Hostname)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", config.ListeningPort))
	if err != nil {
		panic("Unable to bind port")
	}

	fmt.Printf("Listening to port %v\n", config.ListeningPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("Unable to accept connection")
		}

		go handle(conn)
	}
}

func handle(src net.Conn) {
	defer src.Close()

	dst, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		panic("Unable to reach host")
	}

	defer dst.Close()

	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			panic(err)
		}
	}()

	if _, err := io.Copy(src, dst); err != nil {
		panic(err)
	}
}
