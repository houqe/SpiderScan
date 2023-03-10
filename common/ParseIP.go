package common

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ParseIP(host string) (hosts []string, err error) {
	hosts = ParseIPs(host)
	if len(hosts) == 0 && host != "" {
		err = errors.New("[!] 主机解析失败，请检查支持的IP格式")
	}
	return
}

func ParseIPs(ip string) (hosts []string) {
	if strings.Contains(ip, ",") {
		ipList := strings.Split(ip, ",")
		for _, ip := range ipList {
			hosts = append(hosts, parseIP(ip)...)
		}
	} else {
		hosts = parseIP(ip)
	}

	return
}

func parseIP(ip string) []string {
	switch {
	case strings.Contains(ip, "/"):
		return parseIP2(ip)
	case strings.Contains(ip, "-"):
		return parseIP1(ip)
	//处理单个ip
	default:
		testIP := net.ParseIP(ip)
		if testIP == nil {
			return nil
		}
		return []string{ip}
	}

}

// 把 192.168.x.x/xx 转换成 192.168.x.x-192.168.x.x
func parseIP2(ip string) (hosts []string) {
	_, ipNet, err := net.ParseCIDR(ip)
	if err != nil {
		return
	}
	hosts = parseIP1(IPRange(ipNet))
	return
}

// 解析ip段: 192.168.111.1-255
//
//	192.168.111.1-192.168.112.255
func parseIP1(ip string) []string {
	IPRange := strings.Split(ip, "-")
	testIP := net.ParseIP(IPRange[0])
	var AllIP []string
	if len(IPRange[1]) < 4 {
		Range, err := strconv.Atoi(IPRange[1])
		if testIP == nil || Range > 255 || err != nil {
			return nil
		}
		SplitIP := strings.Split(IPRange[0], ".")
		ip1, err1 := strconv.Atoi(SplitIP[3])
		ip2, err2 := strconv.Atoi(IPRange[1])
		PrefixIP := strings.Join(SplitIP[0:3], ".")
		if ip1 > ip2 || err1 != nil || err2 != nil {
			return nil
		}
		for i := ip1; i <= ip2; i++ {
			AllIP = append(AllIP, PrefixIP+"."+strconv.Itoa(i))
		}
	} else {
		SplitIP1 := strings.Split(IPRange[0], ".")
		SplitIP2 := strings.Split(IPRange[1], ".")
		if len(SplitIP1) != 4 || len(SplitIP2) != 4 {
			return nil
		}
		start, end := [4]int{}, [4]int{}
		for i := 0; i < 4; i++ {
			ip1, err1 := strconv.Atoi(SplitIP1[i])
			ip2, err2 := strconv.Atoi(SplitIP2[i])
			if ip1 > ip2 || err1 != nil || err2 != nil {
				return nil
			}
			start[i], end[i] = ip1, ip2
		}
		startNum := start[0]<<24 | start[1]<<16 | start[2]<<8 | start[3]
		endNum := end[0]<<24 | end[1]<<16 | end[2]<<8 | end[3]
		for num := startNum; num <= endNum; num++ {
			ip := strconv.Itoa((num>>24)&0xff) + "." + strconv.Itoa((num>>16)&0xff) + "." + strconv.Itoa((num>>8)&0xff) + "." + strconv.Itoa((num)&0xff)
			AllIP = append(AllIP, ip)
		}
	}
	return AllIP
}

// 获取起始IP、结束IP
func IPRange(c *net.IPNet) string {
	start := c.IP.String()
	mask := c.Mask
	bcst := make(net.IP, len(c.IP))
	copy(bcst, c.IP)
	for i := 0; i < len(mask); i++ {
		ipIdx := len(bcst) - i - 1
		bcst[ipIdx] = c.IP[ipIdx] | ^mask[len(mask)-i-1]
	}
	end := bcst.String()
	return fmt.Sprintf("%s-%s", start, end) //返回用-表示的ip段,192.168.1.0-192.168.255.255
}
