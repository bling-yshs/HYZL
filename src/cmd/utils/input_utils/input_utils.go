package input_utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 持续读取用户输入，直到用户输入的是一个有效的整数
func ReadUint32() uint32 {
	var input uint32
	for {
		_, err := fmt.Scanln(&input)
		if err == nil {
			break
		} else {
			fmt.Println("输入无效，请重新输入：")
		}
	}
	return input
}

func ReadString() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err == nil {
			return strings.TrimSpace(input)
		} else {
			fmt.Println("输入无效，请重新输入：")
		}
	}
}

func ReadChoice(choices []string) string {
	var input string
	for {
		_, err := fmt.Scanln(&input)
		if err == nil {
			for _, choice := range choices {
				if input == choice {
					return input
				}
			}
		}
		fmt.Println("输入无效，请重新输入：")
	}
}
