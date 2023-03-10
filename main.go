package main

import (
	"EnScan/common"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	var args common.Args
	common.Flag(&args)
	fmt.Println(start)
}
