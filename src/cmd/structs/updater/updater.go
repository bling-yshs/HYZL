package updater

import (
	"encoding/json"
	"fmt"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/structs/global"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/http_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/input_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/print_utils"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/hashicorp/go-version"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type updater struct {
	Version    string `json:"version"`
	Url        string `json:"url"`
	Timestamp  int64  `json:"timestamp"`
	Changelog  string `json:"changelog"`
	Deprecated bool   `json:"deprecated"`
}

var Updater = updater{}

var updaters = []updater{}

const url = "https://mirror.ghproxy.com/https://raw.githubusercontent.com/bling-yshs/YzLauncher-windows-updater/main/updater.json"

// 如果有更新，返回true，否则返回false
func CheckUpdate() bool {
	// 当前版本
	current, err := version.NewVersion(global.Global.ProgramVersion)
	if err != nil {
		print_utils.PrintError(err)
		return false
	}
	// 获取最新版本
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(url)
	if err != nil {
		print_utils.PrintError(err)
		return false
	}
	defer response.Body.Close()
	// 解析json
	err = json.NewDecoder(response.Body).Decode(&updaters)
	if err != nil {
		print_utils.PrintError(err)
		return false
	}
	var latestVersionStr string
	// 得到第一个没有废弃的版本
	for index, item := range updaters {
		if !item.Deprecated {
			latestVersionStr = item.Version
			break
		}
		if index == len(updaters)-1 {
			// 如果所有版本都被废弃，返回false
			return false
		}
	}
	latest, err := version.NewVersion(latestVersionStr)
	if err != nil {
		print_utils.PrintError(err)
		return false
	}
	// 如果第一个版本大于当前版本，说明有更新
	if latest.GreaterThan(current) {
		Updater = updaters[0]
		return true
	}
	return false
}

func AskUpdate() bool {
	// 询问是否更新
	fmt.Printf("检测到启动器有更新，是否更新？(y/n)：")
	choice := input_utils.ReadChoice([]string{"y", "n"})
	if choice == "n" {
		return false
	}
	return true
}

func UpdateRightNow() {
	if Updater.Url == "" {
		// 重新获取更新
		CheckUpdate()
	}
	// 下载更新
	print_utils.PrintWithColor(ct.Yellow, true, "正在下载更新...")
	err := http_utils.DownloadFile(Updater.Url, "YzLauncher-windows-new.exe")
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	generateUpdateBat()
	runUpdateBat()
}

func generateUpdateBat() {
	// 生成更新脚本
	str := fmt.Sprintf(`
@echo off
echo 正在更新启动器
REM 延迟一下，等待启动器关闭
ping 127.0.0.1 -n 4 > nul
REM 检查当前目录下是否存在 YzLauncher-windows-new.exe
IF EXIST YzLauncher-windows-new.exe (
	REM 如果存在，替换掉旧的启动器
	RENAME "%s" YzLauncher-windows-old.exe
	RENAME YzLauncher-windows-new.exe "%s"
	REM 删除旧的启动器
	DEL YzLauncher-windows-old.exe
) ELSE (
	REM 如果不存在，打印错误信息
	echo 未找到更新文件，请重新下载
)
`, global.Global.ProgramName)
	os.WriteFile("update.bat", []byte(str), os.ModePerm)
}

func runUpdateBat() {
	// 检查是否存在更新脚本
	if _, err := os.Stat("update.bat"); os.IsNotExist(err) {
		fmt.Println("未找到更新脚本，请重新下载")
		return
	}
	// 执行更新脚本
	print_utils.PrintWithColor(ct.Yellow, true, "正在执行更新脚本...")
	exec.Command("cmd", "/c", "start", "", "update.bat").Start()
	os.Exit(0)
}