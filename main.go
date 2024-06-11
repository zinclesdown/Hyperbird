package main

import (
	"fmt"
)

func main() {
	// 清除终端输出
	fmt.Print("\033[H\033[2J")
	testAll()
}
