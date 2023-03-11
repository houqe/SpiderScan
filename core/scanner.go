package core

import (
	"EnScan/common"
	"fmt"
)

func Scan(args common.Args) {
	fmt.Println(common.Green("[*] start scan..."))
	hosts, err := common.ParseIP(args.Host)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(common.Yellow("[*] 主机测绘"), hosts)
	PortScan(hosts, args.Ports)
}
