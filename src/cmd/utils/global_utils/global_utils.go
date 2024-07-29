package global_utils

import (
	"fmt"
	"github.com/bling-yshs/HYZL/src/cmd/utils/cmd_utils"
	"os"
)

func ShutDownProgram() {
	fmt.Println("按回车键退出...")
	var input string
	_, _ = fmt.Scanln(&input)
	os.Exit(0)
}

func ShutDownYunzai() {
	// 关闭所有node进程
	cmd_utils.ExecuteCmd("taskkill /f /im node.exe", "", "正在关闭云崽...", "云崽关闭成功！")
}

func ShutDownChrome() {
	// 关闭所有chrome进程
	cmd_utils.ExecuteCmd("taskkill /f /im chrome.exe", "", "正在关闭chrome...", "chrome关闭成功！")
}
