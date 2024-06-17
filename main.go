package main

import (
	"fmt"
	"hyperbird/servers"
)

func main() {
	// 清除终端输出
	fmt.Print("\033[H\033[2J")
	testAll()
	servers.StartAllListen("0.0.0.0", "8080")
}
