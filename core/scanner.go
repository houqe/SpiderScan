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
	fmt.Println("[*] 主机测绘", hosts)
}
