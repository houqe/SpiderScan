package protocol

import "strings"

func IsRsyncProtocol(line string) bool {
	return strings.HasPrefix(line, "@RSYNCD")
}
