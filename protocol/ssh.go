package protocol

import (
	"fmt"
	"strings"
)

func IsSSHProtocol(line string) bool {
	return strings.HasPrefix(line, "SSH-")
}

func IsSSHProtocolAPP(line string) string {
	str := strings.ReplaceAll(strings.ReplaceAll(line, "\r", ""), "\n", "")
	if strings.Contains(str, "Comware") {
		return fmt.Sprintf("%-5s|%s", "H3C", str)
	}
	if strings.Contains(str, "Cisco") {
		return fmt.Sprintf("%-5s|%s", "Cisco", str)
	}
	return str
}
