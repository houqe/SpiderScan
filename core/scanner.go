package core

import (
	"EnScan/common"
	"fmt"
)

func Scan(args common.Args) {
	fmt.Println("[*] start scan...")
	hostsList, err := common.ParseIP(args.Host)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(hostsList) > 3 {
		fmt.Println("[*] 主机测绘", hostsList[0], hostsList[1], hostsList[2], "...")
	} else {
		fmt.Println("[*] 主机测绘", hostsList)
	}
	if len(hostsList) > 0 {
		CheckLive(hostsList, common.Ping)
	}

	if args.Ports != "" || common.DftPorts || common.WPorts && len(hostsList) > 0 {
		if args.Ports == "" {
			if common.DftPorts {
				args.Ports = common.DefaultPorts
			}
			if common.WPorts {
				args.Ports = common.WebPorts
			}
		}
		AliveAddress := PortScan(hostsList, args.Ports, common.Timeout, common.Threads)
		fmt.Printf("[*] Port Open number %d\n", len(AliveAddress))
		//for _, addr := range AliveAddress {
		//	fmt.Printf(color.GreenString("[+] %s\n"), addr)
		//}
	}

}
