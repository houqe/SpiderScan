package crack

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/net/context"
	"strconv"
	"sync"
	"time"
)

var ListCrackHost []SussCrack

type SussCrack struct {
	Host   string
	User   string
	Passwd string
	Port   int
	Mode   string
}

var MapCrackHost = make(map[HostPort]SussCrack) //使用 Host和Port作为键，SussCrack结构体作为值。对于每个 Host 和 Port 的组合，只会存储一次弱口令信息，避免并发遇到匿名用户输出。

type HostPort struct {
	Host string
	Port int
}

// ConnectionFunc 定义一个函数类型
type ConnetionFunc func(cancel context.CancelFunc, host, user, passwd string, newport, timeout int)

var connectionFuncs = map[string]ConnetionFunc{
	"ssh":   SSH,
	"mysql": mySql,
}

func Run(host, port, mode string, timeout, chanCount int) {
	ch := make(chan struct{}, chanCount)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel() //确保所有的goroutine都已经退出
	newport, _ := strconv.Atoi(port)
	fmt.Printf(color.YellowString("[*] 弱口令爆破 - %s  %s - %s\n"), mode, host, port)
	for _, user := range Userlist(mode) {
		for _, passwd := range Passwdlist() {
			ch <- struct{}{}
			wg.Add(1)
			if connFunc, ok := connectionFuncs[mode]; ok {
				go crackOnce(ctx, cancel, host, user, passwd, newport, timeout, ch, &wg, connFunc, mode)
			} else {
				wg.Done()
				<-ch
			}
		}
	}
	wg.Wait()
}

func end(host, user, passwd string, port int, mode string) {
	fmt.Printf(color.RedString("[!] 弱口令爆破成功  %s %s:%s %s - %s\n"), mode, host, strconv.Itoa(port), user, passwd)

	//MapCrackHost[HostPort{Host: host, Port: port}] = SussCrack{Host: host, Port: port, Passwd: passwd, User: user, Mode: mode}
}

func crackOnce(ctx context.Context, cancel context.CancelFunc, host, user, passwd string, newport, timeout int,
	ch <-chan struct{}, wg *sync.WaitGroup, connFunc ConnetionFunc, key string) {
	defer done(ch, wg)

	hasDone := make(chan struct{}, 1)
	go func() {
		connFunc(cancel, host, user, passwd, newport, timeout)
		hasDone <- struct{}{}
	}()
}

func done(ch <-chan struct{}, wg *sync.WaitGroup) {
	<-ch
	wg.Done()
}
