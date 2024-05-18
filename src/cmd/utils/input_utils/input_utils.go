package input_utils

import "fmt"

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
	var input string
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
