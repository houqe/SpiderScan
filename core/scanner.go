package core

import (
	"SpiderScan/common"
	"fmt"
)

func Scan(args common.Args) {
	fmt.Println("[*] start scan...")
	hostsList, err := common.ParseIP(args.Host)
	if err != nil {
		log.Warn(err.Error())
		return
	}
	if common.Log {
		log.Info(fmt.Sprintf("[*] 主机测绘: %s", hostsList))
	}
	if len(hostsList) > 3 {
		fmt.Println("[*] 主机测绘", hostsList[0], hostsList[1], hostsList[2], "...")
	} else {
		fmt.Println("[*] 主机测绘", hostsList)
	}
	var AliveHosts []string
	if len(hostsList) > 0 {
		AliveHosts = CheckLive(hostsList, common.Ping)
	}

	if args.Ports != "" || common.DftPorts || common.WPorts && len(AliveHosts) > 0 {
		if args.Ports == "" {
			if common.DftPorts {
				args.Ports = common.DefaultPorts
			}
			if common.WPorts {
				args.Ports = common.WebPorts
			}
		}
		AliveAddress := PortScan(AliveHosts, args.Ports, common.Timeout, common.Threads)
		if common.Log {
			log.Info(fmt.Sprintf("[*] 端口探活完成，存活数量：%d", len(AliveAddress)))
		}
		fmt.Printf("[*] Port open number %d\n", len(AliveAddress))
		//for _, addr := range AliveAddress {
		//	fmt.Printf(color.GreenString("[+] %s\n"), addr)
		//}
	}

}
