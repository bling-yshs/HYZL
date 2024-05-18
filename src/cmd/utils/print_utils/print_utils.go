package print_utils

import (
	"fmt"
	ct "github.com/daviddengcn/go-colortext"
)

func PrintError(a ...any) {
	ct.Foreground(ct.Red, true)
	fmt.Println("发生了以下错误：")
	fmt.Println(a...)
	ct.ResetColor()
}

func PrintWithEmptyLine(a ...any) {
	fmt.Println()
	fmt.Println(a...)
	fmt.Println()
}
func PrintWithColor(color ct.Color, bright bool, a ...any) {
	ct.Foreground(color, bright)
	fmt.Println(a...)
	ct.ResetColor()
}
