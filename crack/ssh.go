package crack

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

func SSH(cancel context.CancelFunc, host, user, passwd string, port, timeout int) {
	configssh := &ssh.ClientConfig{
		Timeout:         time.Duration(timeout) * time.Second, //ssh连接timeout时间
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	configssh.Auth = []ssh.AuthMethod{ssh.Password(passwd)}

	addr := fmt.Sprintf("%s:%d", host, port)
	sshClient, err := ssh.Dial("tcp", addr, configssh)
	if err != nil {
		return
	}
	defer sshClient.Close()
	end(host, user, passwd, port, "SSH")
	cancel()
}
