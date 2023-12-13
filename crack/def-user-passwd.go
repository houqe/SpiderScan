package crack

import (
	"SpiderScan/common"
	"embed"
	"strings"
)

var userMap = map[string][]string{
	"ssh":        {"root"},
	"mysql":      {"root"},
	"redis":      {"", "root"},
	"postgresql": {"postgres", "root"},
	"sqlserver":  {"sa", "administrator"},
	"ftp":        {"ftp", "admin"},
	"smb":        {"administrator", "guest"},
	"telnet":     {"admin", "root"},
	"tomcat":     {"tomcat", "manager", "admin"},
	"rdp":        {"administrator"},
	"oracle":     {"orcl", "sys", "system"},
}

//go:embed password.txt
var passwd embed.FS

func Passwdlist() []string {
	var passwdlist []string
	data, _ := passwd.ReadFile("password.txt")
	datastr := strings.ReplaceAll(string(data), "\r\n", "\n")
	for _, u := range strings.Split(datastr, "\n") {
		passwdlist = append(passwdlist, u)
	}
	passwdlist = common.RemoveDuplicates(passwdlist)
	return passwdlist
}

func Userlist(mode string) []string {
	return userMap[mode]
}
