package common

func Banner() {
	banner := `
 _____ _      ____  ____ ____  _     
/  __// \  /|/ ___\/   _Y  _ \/ \  /|
|  \  | |\ |||    \|  / | / \|| |\ ||
|  /_ | | \||\___ ||  \_| |-||| | \||
\____\\_/  \|\____/\____|_/ \|\_/  \|
	EnScan version: ` + version + `
`
	print(banner)
}

func Flag(args *Args) {
	Banner()
	flag.StringVar(&Args.Host, "h", "", "IP address of the host you want to scan,for example: 192.168.11.11 | 192.168.11.11-255 | 192.168.11.11,192.168.11.12")

}
