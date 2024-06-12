package user

import (
	"fmt"
	"log"
)

// func assert(msg string, err error) {
// 	if err != nil {
// 		fmt.Printf("\033[31m%s\033[0m", msg) // 使用红色文本打印msg
// 		log.Fatal(err)
// 	}
// }

func warn(msg string, err error) {
	if err != nil {
		fmt.Printf("\033[33m%s\033[0m", msg) // 使用黄色文本打印msg
		log.Println(err)
	}
}
