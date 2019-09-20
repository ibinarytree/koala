package util

import (
	"fmt"
	"net"

	"sync/atomic"
)

var ipAuto atomic.Value

func GetLocalIP() (ip string, err error) {

	ip, ok := ipAuto.Load().(string)
	if ok && len(ip) > 0 {
		return
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				ipAuto.Store(ip)
				return
			}

		}
	}

	err = fmt.Errorf("get local ip failed")
	return
}
