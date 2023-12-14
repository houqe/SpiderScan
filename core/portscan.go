package core

import (
	"SpiderScan/common"
	"SpiderScan/crack"
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"strings"
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
	if common.Log {
		log.Info(fmt.Sprintf("[*] 端口测绘: %d", portsList))
	}
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
		if err != nil {
			if common.Log {
				log.Info(fmt.Sprintf("[-] %s:%d 连接失败", host, port))
			}
		} else {
			conn.Close()
		}
	}()
	if err == nil {
		address := host + ":" + strconv.Itoa(port)
		if common.Log {
			log.Info(fmt.Sprintf("[+] %s:%d 连接成功", host, port))
		}
		protocol := ""
		if common.Service || common.Crack {
			protocol = parseProtocol(conn, host, strconv.Itoa(port))
			if common.Crack {
				// 支持遍历字典的扫描类型
				protocols := []string{"ssh", "mysql", "redis", "telnet"}
				for _, proto := range protocols {
					if strings.Contains(strings.ToLower(protocol), proto) {
						crack.Run(host, strconv.Itoa(port), proto, timout, 5)
						break
					}
				}
			}
		}
		result := fmt.Sprintf("%s %s", address, protocol)
		fmt.Printf(color.GreenString("[+] %s\n"), result)
		wg.Add(1)
		results <- result
	}
}
