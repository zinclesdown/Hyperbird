package booklibrary

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

func warn(msg string, err error) {
	if err != nil {
		fmt.Printf("\033[33m%s\033[0m", msg) // 使用黄色文本打印msg
		log.Println(err)
	}
}

func assert(msg string, err error) {
	if err != nil {
		color.Red(msg) // 使用红色文本打印msg
		log.Fatal(err)
	}
}
