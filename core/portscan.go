package core

import (
	"EnScan/common"
	"fmt"
)

func PortScan(hosts []string, ports string) []string {
	portList := common.ParsePort(ports)
	fmt.Println(common.Yellow("[*] 端口测绘"), portList)
	return nil
}
