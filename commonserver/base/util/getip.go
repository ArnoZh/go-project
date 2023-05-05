// Package util .

package util

import (
	"net"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// GetIntranetIP 获取本机内网IP，只返回其中一个
func GetIntranetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logrus.Errorf("get intranet ip failed: %v", err)
		return ""
	}
	var ips []string
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	if len(ips) > 0 {
		return ips[0]
	}
	return ""
}

// GetNetLastNumber 获取网卡最后一段
func GetNetLastNumber() int {
	strip := GetIntranetIP()
	substrs := strings.Split(strip, ".")
	cnt := len(substrs)
	if cnt > 0 {
		data, _ := strconv.Atoi(substrs[cnt-1])
		return data
	}
	return -1
}
