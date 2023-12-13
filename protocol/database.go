package protocol

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net"
	"strings"
	"time"
)

func IsPgsqlProtocol(host, port string) bool {
	if !strings.Contains(port, "5432") {
		return false
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable  TimeZone=Asia/Shanghai connect_timeout=%d", host, "postgres", "123456", port, 3)
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), //设置日志级别
	})
	if err != nil {
		if strings.Contains(err.Error(), "server error") || strings.Contains(err.Error(), "SQLSTATE") || strings.Contains(err.Error(), "致命错误") {
			return true
		}
		return false
	} else {
		return true
	}
}

func IsMySqlProtocol(host, port string) (bool, string) {
	// 按照常理来说，端口号小于4位必然不是MySQL
	if len(port) < 4 {
		return false, ""
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return false, ""
	}
	defer conn.Close()

	err = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		return false, ""
	}

	buf := make([]byte, 4)
	_, err = conn.Read(buf)
	if err != nil {
		return false, ""
	}

	buf = make([]byte, 1)
	_, err = conn.Read(buf)
	if err != nil {
		return false, ""
	}

	if buf[0] != 10 {
		//fmt.Println("Not a MySQL server")
		return false, ""
	}

	buf = make([]byte, 1)
	version := ""
	for {
		_, err = conn.Read(buf)
		if err != nil {
			return false, ""
		}
		if buf[0] == 0 {
			break
		}
		version += string(buf)
	}
	return true, version
}
