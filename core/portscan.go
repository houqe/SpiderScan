package core

import (
	"EnScan/common"
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"sync"
	"time"
)

type Addr struct {
	ip   string
	port int
}

func PortScan(hostsList []string, ports string, timout int, threads int) []string {
	var AliveAddress []string
	//println(len(hostsList))
	portsList := common.ParsePort(ports)
	fmt.Println("[*] 端口测绘", portsList)
	Addrs := make(chan Addr, len(hostsList)*len(portsList))
	results := make(chan string, len(hostsList)*len(portsList))
	var wg sync.WaitGroup

	go func() {
		for found := range results {
			AliveAddress = append(AliveAddress, found)
			wg.Done()
		}
	}()

	for i := 0; i < threads; i++ {
		go func() {
			for addr := range Addrs {
				PortConnect(addr, results, timout, &wg)
				wg.Done()
			}
		}()
	}

	for _, port := range portsList {
		for _, host := range hostsList {
			wg.Add(1)
			Addrs <- Addr{host, port}
		}
	}

	wg.Wait()
	close(Addrs)
	close(results)
	return AliveAddress
}

func PortConnect(addr Addr, results chan string, timout int, wg *sync.WaitGroup) {
	host, port := addr.ip, addr.port
	conn, err := common.WrapperTcpWithTimeout("tcp4", fmt.Sprintf("%s:%v", host, port), time.Duration(timout)*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err == nil {
		address := host + ":" + strconv.Itoa(port)
		result := fmt.Sprintf("%s open", address)
		fmt.Printf(color.GreenString("[+] %s\n"), result)
		wg.Add(1)
		results <- result
	}
}
