package core

import (
	"EnScan/common"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/net/icmp"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	AliveHosts []string
	ExitHosts  = make(map[string]struct{})
	OS         = runtime.GOOS
	livewg     sync.WaitGroup
)

func CheckLive(hostsList []string, ping bool) []string {
	chanHosts := make(chan string, len(hostsList))
	go func() {
		for ip := range chanHosts {
			if _, ok := ExitHosts[ip]; !ok {
				ExitHosts[ip] = struct{}{}
				if ping {
					fmt.Printf(color.GreenString("[+] Target %-15s is alive for Ping\n"), ip)
				} else {
					fmt.Printf(color.GreenString("[+] Target %-15s is alive for ICMP\n"), ip)
				}
				AliveHosts = append(AliveHosts, ip)
			}
			livewg.Done()
		}
	}()

	if ping == true {
		RunPing(hostsList, chanHosts)
	} else {
		//优先尝试监听本地icmp,批量探测
		conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err == nil {
			if common.Log {
				log.Info("[*] 本地监听ICMP探测")
			}
			RunIcmp1(hostsList, conn, chanHosts)
		} else {
			//无监听icmp探测
			if common.Log {
				log.Info("[*] 本地监听ICMP探测不可用,使用无监听探测方式")
			}
			conn, err := net.DialTimeout("ip4:icmp", "127.0.0.1", 3*time.Second)
			defer func() {
				if conn != nil {
					conn.Close()
				}
			}()
			if err == nil {
				RunIcmp2(hostsList, chanHosts)
			} else {
				//使用ping探测
				if common.Log {
					log.Info("[*] 本地监听ICMP探测不可用,无监听探测方式不可用,使用PING探测")
				}
				fmt.Println(color.YellowString("Failed to send ICMP packet"))
				fmt.Println("start ping")
				RunPing(hostsList, chanHosts)
			}
		}
	}
	livewg.Wait()
	close(chanHosts)
	fmt.Printf("[*] Host alive number %d\n", len(AliveHosts))
	if common.Log {
		log.Info(fmt.Sprintf("[*] 主机探活完成，存活数量：%d", len(AliveHosts)))
	}
	return AliveHosts
}

func RunIcmp2(hostsList []string, chanHosts chan string) {
	var wg sync.WaitGroup
	for _, host := range hostsList {
		wg.Add(1)
		go func(host string) {
			if icmpAlive(host) {
				livewg.Add(1)
				chanHosts <- host
			}
			wg.Done()
		}(host)
	}
	wg.Wait()
}

func RunIcmp1(hostsList []string, conn *icmp.PacketConn, chanHosts chan string) {
	endflag := false
	go func() {
		for {
			if endflag == true {
				return
			}
			msg := make([]byte, 100)
			_, sourceIP, _ := conn.ReadFrom(msg)
			if sourceIP != nil {
				livewg.Add(1)
				chanHosts <- sourceIP.String()
			}
		}
	}()

	for _, host := range hostsList {
		dst, _ := net.ResolveIPAddr("ip", host)
		IcmpByte := makemsg(host)
		conn.WriteTo(IcmpByte, dst)
	}
	if common.Timeout == 0 {
		if len(hostsList) > 256 {
			common.Timeout = 6
		} else {
			common.Timeout = 3
		}
	}
	//睡眠是否会影响探测？考虑使用for阻塞？
	time.Sleep(time.Duration(common.Timeout))

	endflag = true
	conn.Close()
}

func RunPing(hostsList []string, chanHosts chan string) {
	var wg sync.WaitGroup
	//limiter := make(chan struct{}, 50)
	for _, host := range hostsList {
		wg.Add(1)
		//limiter <- struct{}{}
		go func(host string) {
			if ExecCommandPing(host) {
				livewg.Add(1)
				chanHosts <- host
			}
			//<-limiter
			wg.Done()
		}(host)
	}
	wg.Wait()
}
func ExecCommandPing(ip string) bool {
	var command *exec.Cmd
	if OS == "windows" {
		command = exec.Command("cmd", "/c", "ping -n 1 -w 1 "+ip+" && echo true || echo false")
	} else if OS == "linux" {
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -w 1 "+ip+" >/dev/null && echo true || echo false")
	} else if OS == "darwin" {
		//apple
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -W 1 "+ip+" >/dev/null && echo true || echo false")
	}
	outInfo := bytes.Buffer{}
	command.Stdout = &outInfo
	err := command.Start()
	if err != nil {
		if common.Log {
			log.Error("[-] COMMAND 启动失败！")
		}
		return false
	}
	if err = command.Wait(); err != nil {
		if common.Log {
			log.Error("[-] PING 命令执行失败！")
		}
		return false
	} else {
		if strings.Contains(outInfo.String(), "true") {
			return true
		} else {
			return false
		}
	}

}

func icmpAlive(host string) bool {
	startTime := time.Now()
	conn, err := net.DialTimeout("ip4:icmp", host, 6*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	if err != nil {
		return false
	}
	//带有超时属性的icmp 这相当于同时调用 SetReadDeadline和SetWriteDeadline。
	if err := conn.SetDeadline(startTime.Add(6 * time.Second)); err != nil {
		return false
	}

	msg := makemsg(host)
	if _, err := conn.Write(msg); err != nil {
		return false
	}

	receive := make([]byte, 60)
	if _, err := conn.Read(receive); err != nil {
		return false
	}

	return true
}

func makemsg(host string) []byte {
	msg := make([]byte, 40)
	id0, id1 := genIdentifier(host)
	msg[0] = 8
	msg[1] = 0
	msg[2] = 0
	msg[3] = 0
	msg[4], msg[5] = id0, id1
	msg[6], msg[7] = genSequence(1)
	check := checkSum(msg[0:40])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)
	return msg
}

func checkSum(msg []byte) uint16 {
	sum := 0
	length := len(msg)
	for i := 0; i < length-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if length%2 == 1 {
		sum += int(msg[length-1]) * 256
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)
	answer := uint16(^sum)
	return answer
}

func genSequence(v int16) (byte, byte) {
	ret1 := byte(v >> 8)
	ret2 := byte(v & 255)
	return ret1, ret2
}

func genIdentifier(host string) (byte, byte) {
	return host[0], host[1]
}
