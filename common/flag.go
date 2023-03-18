package common

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
)

func Banner() {
	banner := `
	 _____ _      ____  ____ ____  _     
	/  __// \  /|/ ___\/   _Y  _ \/ \  /|
	|  \  | |\ |||    \|  / | / \|| |\ ||
	|  /_ | | \||\___ ||  \_| |-||| | \||
	\____\\_/  \|\____/\____|_/ \|\_/  \|
		EnScan version: ` + version + `
`
	fmt.Println(color.BlueString(banner))
}

func Flag(args *Args) {
	Banner()
	flag.StringVar(&args.Host, "h", "", "IP address of the host you want to scan,for example: 192.168.11.11 | 192.168.11.11-255 | 192.168.11.0/24 | 192.168.11.11,192.168.11.12")
	flag.StringVar(&args.Ports, "p", DefaultPorts, "Select a port,for example: 22 | 1-65535 | 22,80,3306")
	flag.BoolVar(&Ping, "ping", false, "using ping replace icmp")
	flag.IntVar(&Threads, "t", 600, "Use Thread nums")
	flag.IntVar(&Timeout, "time", 3, "Set timeout")
	flag.Parse()
}
