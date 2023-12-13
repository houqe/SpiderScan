package core

import (
	"SpiderScan/common"
	"SpiderScan/protocol"
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func parseProtocol(conn net.Conn, host, port string) string {
	// 首先根据默认端口去对应服务
	if protocol, ok := common.PortProtocols[port]; ok {
		return protocol
	}
	// 设置读取超时时间
	if err := conn.SetReadDeadline(time.Now().Add(time.Duration(common.Timeout) * time.Second)); err != nil {
		log.Error("[!] 设置读取超时时间失败")
		return ""
	}

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		line = ""
	}

	switch {
	case protocol.IsSSHProtocol(line):
		return protocol.IsSSHProtocolAPP(line)
	case strings.HasPrefix(line, "220"):
		return "FTP"
	case protocol.IsTelnet(conn):
		return "Telnet"
	case protocol.IsRedisProtocol(conn):
		return "数据库|Redis"
	case protocol.IsPgsqlProtocol(host, port):
		return "数据库|PostgreSQL"
	case protocol.IsRsyncProtocol(line):
		return "rsync|" + line
	default:
		isWeb := protocol.IsWeb(host, port, common.Timeout)
		if isWeb != "" {
			return fmt.Sprintf("%-5s| %s", "WEB应用", isWeb)
		}
	}
	isMySQL, version := protocol.IsMySqlProtocol(host, port)
	if isMySQL {
		return fmt.Sprintf("数据库|MySQL:%s", version)
	}
	return defaultPort(port)
}

func defaultPort(port string) string {
	defMap := map[string]string{
		"3306":  "数据库|MySQL",
		"23":    "Telnet",
		"21":    "FTP",
		"80":    "WEB应用",
		"443":   "WEB应用",
		"61616": "ActiveMQ",
	}
	value, exists := defMap[port]
	if exists {
		return value
	}
	return ""
}
