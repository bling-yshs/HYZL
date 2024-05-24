package print_utils

import (
	"fmt"
	ct "github.com/daviddengcn/go-colortext"
	"runtime/debug"
)

func PrintError(err error) {
	ct.Foreground(ct.Red, true)
	fmt.Printf("发生了以下错误，请将此界面截图并反馈给作者：")
	stack := debug.Stack()
	fmt.Println(string(stack))
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
