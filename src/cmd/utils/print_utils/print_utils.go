package print_utils

import (
	"fmt"
	ct "github.com/daviddengcn/go-colortext"
)

func PrintError(err error) {
	ct.Foreground(ct.Magenta, true)
	fmt.Println("发生了以下错误，如有疑问，请将此界面截图并反馈给作者：")
	ct.Foreground(ct.Red, true)
	fmt.Printf("%+v\n", err)
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
