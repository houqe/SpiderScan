package protocol

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func IsRedisProtocol(conn net.Conn) bool {
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, err := fmt.Fprintf(conn, "*1\n$4\nPING\n")
	if err != nil {
		return false
	}
	response, err := readResponse(conn)
	if err == nil {
		if strings.Contains(response, "NOAUTH") || strings.Contains(response, "PONG") {
			return true
		}
	}
	return false
}

func readResponse(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return line, nil
}
