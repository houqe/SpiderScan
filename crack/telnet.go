package crack

import (
	"context"
	"fmt"
	"github.com/ziutek/telnet"
	"io"
	"strings"
	"time"
)

func telnetcon(cancel context.CancelFunc, host, user, passwd string, port, timeout int) {
	conn, err := telnet.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Duration(timeout))
	if err == nil {
		return
	}
	err = login(conn, user, passwd)
	if err == nil {
		if readOutput(conn) {
			defer conn.Close()
			end(host, user, passwd, port, "Telnet")
			cancel()
		}
	}
	return
}

func readOutput(conn *telnet.Conn) bool {
	for {
		data, err := conn.ReadUntil("\n")
		if len(data) > 0 {
			dataStr := strings.TrimSpace(string(data))
			if dataStr != "" {
				if strings.Count(dataStr, "incorrect") > 0 {
					return false
				}
				if strings.Count(dataStr, "Last") > 0 {
					return true
				}
			}
		}
		if err != nil {
			return false
		}
	}
}

func login(conn *telnet.Conn, user string, passwd string) error {
	conn.SetUnixWriteMode(true)
	conn.SkipUntil("login: ")
	if _, err := io.WriteString(conn, fmt.Sprintf("%s\n", user)); err != nil {
		return err
	}
	conn.SkipUntil("Password: ")
	if _, err := io.WriteString(conn, fmt.Sprintf("%s\n", passwd)); err != nil {
		return err
	}
	return nil
}
