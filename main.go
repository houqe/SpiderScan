package main

import (
	"EnScan/common"
	"EnScan/core"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	var args common.Args
	common.Flag(&args)
	core.Scan(args)
	t := time.Now().Sub(start)
	fmt.Printf("[*] 扫描结束,耗时: %s\n", t)
}
