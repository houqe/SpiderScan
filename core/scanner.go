package core

import (
	"EnScan/common"
	"fmt"
)

func Scan(args common.Args) {
	fmt.Println("[*] start scan...")
	hosts, err := common.ParseIP(args.Host)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(hosts) > 3 {
		fmt.Println("[*] 主机测绘", hosts[0], hosts[1], hosts[2], "...")
	} else {
		fmt.Println("[*] 主机测绘", hosts)
	}
	PortScan(hosts, args.Ports)
	if common.Ping && len(hosts) > 0 {
		//println("进行主机探活")
		CheckLive(hosts, common.Ping)
	}
}
