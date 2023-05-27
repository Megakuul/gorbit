package handler

import (
	"fmt"
	"github/megakuul/gorbit/conf"
	"net"
	"time"
)

func SelectEndpoint(endpoints *[]conf.Endpoint) (endpointIndex int, err error) {
	selectedIndex := -1

	for i := 0; i < len(*endpoints); i++ {
		if !(*endpoints)[i].Healthy {
			continue
		}
		if selectedIndex == -1 {
			selectedIndex = i
			continue
		}
		if (*endpoints)[selectedIndex].Sessions > (*endpoints)[i].Sessions {
			selectedIndex = i
			continue
		}
	}

	if selectedIndex == -1 {
		return -1, fmt.Errorf("no healthy endpoint available")
	} else {
		return selectedIndex, nil
	}
}

func CheckHealth(endpoints *[]conf.Endpoint, sintervall int) {
	for {
		for i := 0; i < len(*endpoints); i++ {
			if err := SendHealthCheck(&(*endpoints)[i]); err != nil {
				(*endpoints)[i].Healthy = false
			} else {
				(*endpoints)[i].Healthy = true
			}
		}
		time.Sleep(time.Duration(sintervall) * time.Second)
	}
}

func SendHealthCheck(endpoint *conf.Endpoint) error {

	addr, err := net.ResolveTCPAddr("tcp",
		fmt.Sprintf("%s:%v", endpoint.Hostname, endpoint.Port),
	)
	if err != nil {
		return err
	}

	// Create a TCP Connection to the desired Endpoint
	dstCon, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	defer dstCon.Close()

	return nil
}
