package windows_utils

import (
	"github.com/James-Ye/go-frame/win"
	"os/exec"
	"strings"
)

func GetWindowsVersion() uint32 {
	version, _, _ := win.RtlGetNtVersionNumbers()
	return version
}

// 检查进程列表，是否存在指定进程
func IsProcessRunning(processName string) bool {
	cmd := exec.Command("cmd", "/c", "tasklist")
	output, _ := cmd.Output()
	if strings.Contains(string(output), processName) {
		return true
	} else {
		return false
	}
}
