package handler

import (
	"fmt"
	"github/megakuul/gorbit/conf"
	"net"
	"time"
)

func SelectEndpoint(endpoints []conf.Endpoint) (endpointIndex int, err error) {
	selectedIndex := -1

	for i, endpoint := range endpoints {
		fmt.Printf("%v", endpoint)
		if !endpoint.Healthy {
			continue
		}
		if endpoints[selectedIndex].Sessions > endpoint.Sessions {
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

func CheckHealth(endpoints []conf.Endpoint, sintervall int) {
	for {
		for _, endpoint := range endpoints {
			if err := SendHealthCheck(endpoint); err != nil {
				endpoint.Healthy = false
			} else {
				endpoint.Healthy = true
			}
		}
		time.Sleep(time.Duration(sintervall) * time.Second)
	}
}

func SendHealthCheck(endpoint conf.Endpoint) error {

	addr, err := net.ResolveTCPAddr("tcp",
		fmt.Sprintf("%s:%v", endpoint.Hostname, endpoint.Port),
	)
	if err != nil {
		return err
	}

	// Create a TCP Connection to the desired Endpoint
	dstCon, err := net.DialTCP("tcp", nil, addr)
	defer dstCon.Close()
	if err != nil {
		return err
	}

	return nil
}
