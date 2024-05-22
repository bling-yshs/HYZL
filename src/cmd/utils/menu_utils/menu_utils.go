package menu_utils

import (
	"fmt"
	"github.com/bling-yshs/HYZL/src/cmd/structs/menu_option"
	"github.com/bling-yshs/HYZL/src/cmd/utils/cmd_utils"
	"os"
)

func PrintMenu(title string, options []menu_option.MenuOption, isMainMenu bool) {
	fmt.Println("===" + title + "===")
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option.Label)
	}
	if isMainMenu {
		fmt.Println("0. 退出程序")
	} else {
		fmt.Println("0. 返回上一级")
	}

	fmt.Print("\n请选择操作：")
}

func DealChoice(choice uint32, options []menu_option.MenuOption, isMainMenu bool) {
	cmd_utils.ClearLog()
	// 如果是主菜单并且选择了0就退出程序
	if isMainMenu && choice == 0 {
		os.Exit(0)
		return
	}
	// 如果是子菜单并且选择了0就返回上一级
	if !isMainMenu && choice == 0 {
		return
	}
	// 如果选择了不存在的选项
	if choice > uint32(len(options)) {
		fmt.Println("输入有误，请重新输入")
		return
	}
	options[choice-1].Action()
}
