package cmd_utils

import (
	"fmt"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
	ct "github.com/daviddengcn/go-colortext"
	"os"
	"os/exec"
	"syscall"
)

func CheckCommandExist(command string) bool {
	cmd := exec.Command("cmd", "/c", command)
	err := cmd.Run()
	if err == nil {
		return true
	} else {
		return false
	}
}

// 检查命令是否存在，如果存在返回cmd的返回值，否则返回错误
func CheckCommand(command string) (string, error) {
	cmd := exec.Command("cmd", "/c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	} else {
		return string(output), nil
	}

}

func ClearLog() {
	ExecuteCmd("cls", "", "", "")
}

func ExecuteCmd(command string, dir string, beforeMsg string, afterMsg string) {
	cmd := exec.Command("cmd.exe")
	cmd.Stdout = os.Stdout // 直接将命令标准输出连接到标准输出流
	cmd.Stderr = os.Stderr // 将错误输出连接到标准错误流
	cmd.Stdin = os.Stdin   // 将标准输入连接到命令的标准输入
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c %s`, command), HideWindow: true}
	if dir != "" {
		cmd.Dir = dir
	}
	if beforeMsg != "" {
		print_utils.PrintWithEmptyLine(beforeMsg)
	}
	print_utils.PrintWithColor(ct.Green, true, "正在执行命令："+command)
	cmd.Run()
	if afterMsg != "" {
		print_utils.PrintWithEmptyLine(afterMsg)
	}
}
