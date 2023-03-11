package core

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var (
	AliveHosts []string
	ExitHosts  = make(map[string]struct{})
	OS         = runtime.GOOS
	livewg     sync.WaitGroup
)

func CheckLive(hostlist []string, ping bool) []string {
	chanHosts := make(chan string, len(hostlist))
	go func() {
		for ip := range chanHosts {
			if _, ok := ExitHosts[ip]; !ok {
				ExitHosts[ip] = struct{}{}
				if ping {
					fmt.Printf("[+] Target %-15s is alive for ping\n", ip)
				}
				AliveHosts = append(AliveHosts, ip)
			}
			livewg.Done()
		}
	}()

	if ping {
		RunPing(hostlist, chanHosts)
	}
	livewg.Wait()
	close(chanHosts)
	fmt.Printf("[*] Host alive number %d\n", len(AliveHosts))
	return AliveHosts
}

func RunPing(hostlist []string, chanHosts chan string) {
	var wg sync.WaitGroup
	//limiter := make(chan struct{}, 50)
	for _, host := range hostlist {
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
	outinfo := bytes.Buffer{}
	command.Stdout = &outinfo
	err := command.Start()
	if err != nil {
		return false
	}
	if err = command.Wait(); err != nil {
		return false
	} else {
		if strings.Contains(outinfo.String(), "true") {
			return true
		} else {
			return false
		}
	}

}
