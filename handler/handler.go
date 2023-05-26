package handler

import (
	"fmt"
	"github/megakuul/gorbit/conf"
	"github/megakuul/gorbit/logger"
	"io"
	"net"
	"sync"
	"time"
)

type Session struct {
	Conn         net.Conn
	Creationtime time.Time
	Timeout      time.Duration
}

type LoadBalancer struct {
	Sessions     map[string]*Session
	SessionMutex sync.RWMutex
}

func (lb *LoadBalancer) HandleConnection(srcCon net.Conn, config conf.Config) {

	session := &Session{
		Conn:         srcCon,
		Creationtime: time.Now(),
		Timeout:      time.Duration(config.Endpoints[0].Timeout_ms) * time.Millisecond,
	}

	lb.SessionMutex.Lock()
	lb.Sessions[srcCon.LocalAddr().String()] = session
	lb.SessionMutex.Unlock()

	defer func() {
		lb.SessionMutex.Lock()
		delete(lb.Sessions, srcCon.LocalAddr().String())
		lb.SessionMutex.Unlock()
	}()

	defer srcCon.Close()

	dstCon, err := net.Dial("tcp", fmt.Sprintf("%s:%v", config.Endpoints[0].Hostname, config.Endpoints[0].Port))
	defer dstCon.Close()
	if err != nil {
		logger.WriteInformationLogger(
			fmt.Sprintf("%s failed to reach %s", srcCon.LocalAddr().String(), dstCon.LocalAddr().String()),
		)
		return
	}

	go func() {
		if _, err := io.Copy(dstCon, srcCon); err != nil {
			logger.WriteInformationLogger(
				fmt.Sprintf("%s", err),
			)
		}
	}()

	if _, err := io.Copy(srcCon, dstCon); err != nil {
		logger.WriteInformationLogger(
			fmt.Sprintf("%s", err),
		)
	}
}
