package common

import (
	"strconv"
	"strings"
)

func ParsePort(port string) (ports []int) {
	if port == "" {
		return
	}
	slices := strings.Split(port, ",")
	for _, port := range slices {
		port = strings.TrimSpace(port)
		if port == "" {
			continue
		}
		upper := port
		if strings.Contains(port, "-") {
			ranges := strings.Split(port, "-")
			if len(ranges) < 2 {
				continue
			}
			startPort, _ := strconv.Atoi(ranges[0])
			endPort, _ := strconv.Atoi(ranges[1])
			if startPort < endPort {
				port = ranges[0]
				upper = ranges[1]
			} else {
				port = ranges[1]
				upper = ranges[0]
			}
		}
		start, _ := strconv.Atoi(port)
		end, _ := strconv.Atoi(upper)
		for i := start; i <= end; i++ {
			ports = append(ports, i)
		}
	}
	return removeDuplicatePort(ports)
}

func removeDuplicatePort(old []int) []int {
	result := []int{}
	temp := map[int]struct{}{}
	for _, item := range old {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// 一些默认端口匹配的服务
var PortProtocols = map[string]string{
	"25":    "SMTP",
	"53":    "DNS",
	"110":   "POP3",
	"135":   "RPC服务",
	"137":   "NetBIOS名称服务",
	"138":   "NetBIOS数据报服务",
	"139":   "NetBIOS会话服务",
	"161":   "SNMP",
	"162":   "SNMP-trap",
	"143":   "IMAP",
	"445":   "SMB",
	"465":   "SMTPS",
	"514":   "syslog",
	"993":   "IMAPS",
	"995":   "POP3S",
	"1433":  "数据库|SqlServer",
	"1521":  "数据库|Oracle",
	"1723":  "PPTP",
	"2049":  "NFS",
	"3389":  "RDP",
	"5900":  "VNC",
	"5901":  "VNC",
	"5672":  "RabbitMq",
	"27017": "数据库|MongoDB",
	"2181":  "ZooKeeper",
}
