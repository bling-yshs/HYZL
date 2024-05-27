package updater

import (
	"encoding/json"
	"fmt"
	"github.com/bling-yshs/HYZL/src/cmd/structs/global"
	"github.com/bling-yshs/HYZL/src/cmd/utils/http_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/input_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/io_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/hashicorp/go-version"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type updater struct {
	Version    string `json:"version"`
	Url        string `json:"url"`
	MD5        string `json:"md5"`
	Timestamp  int64  `json:"timestamp"`
	Changelog  string `json:"changelog"`
	Deprecated bool   `json:"deprecated"`
}

var updaterInstance = updater{}

var updaters = []updater{}

const url = "https://mirror.ghproxy.com/https://raw.githubusercontent.com/bling-yshs/HYZL-updater/main/updater.json"

// 判断缓存的更新文件是否是最新的
func IsUpdateTempNew() bool {
	// 从json中读取更新文件信息
	bytes, err := os.ReadFile("./config/updater.json")
	if err != nil {
		print_utils.PrintError(err)
		return false
	}
	var updaterJson updater
	err = json.Unmarshal(bytes, &updaterJson)
	if err != nil {
		print_utils.PrintError(err)
		return false
	}
	tempStr := updaterJson.Version
	// 获取最新版本号
	latestStr := updaterInstance.Version
	// 比较版本号
	tempVersion, err := version.NewVersion(tempStr)
	if err != nil {
		print_utils.PrintError(err)
		return false
	}
	latestVersion, err := version.NewVersion(latestStr)
	if err != nil {
		print_utils.PrintError(err)
		return false
	}
	return tempVersion.GreaterThanOrEqual(latestVersion)
}

// 写入config/updater.json
func WriteUpdaterJson() {
	// 检查如果不存在updater.json，则创建
	if _, err := os.Stat("./config/updater.json"); os.IsNotExist(err) {
		_, err := os.Create("./config/updater.json")
		if err != nil {
			print_utils.PrintError(err)
			return
		}
	}
	// 将updaterInstance写入updater.json
	jsonBytes, err := json.Marshal(updaterInstance)
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	err = os.WriteFile("./config/updater.json", jsonBytes, os.ModePerm)
	if err != nil {
		print_utils.PrintError(err)
		return
	}
}

// 检查config文件夹下是否存在更新文件，有则返回true，否则返回false
func UpdateTempExist() bool {
	_, err := os.Stat("./config/HYZL-new.exe")
	// 如果存在
	if err == nil {
		return true
	}
	// 如果不存在
	return false
}

func CleanUpdateTemp() {
	// 判断config文件夹下是否存在更新文件，有则删除
	if _, err := os.Stat("./config/HYZL-new.exe"); err == nil {
		os.Remove("./config/HYZL-new.exe")
	}
}

func CleanUpdater() {
	// 如果当前目录下存在更新脚本，删除
	if _, err := os.Stat("update.bat"); err == nil {
		os.Remove("update.bat")
	}
}

// 如果有更新，返回true，否则返回false
func CheckForUpdate() bool {
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
	// 得到第一个没有废弃的版本
	var latestVersionStr string
	var latestVersionIndex int
	for index, item := range updaters {
		if !item.Deprecated {
			latestVersionStr = item.Version
			latestVersionIndex = index
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
	updaterInstance = updaters[latestVersionIndex]
	// 如果第一个版本大于当前版本，说明有更新
	if latest.GreaterThan(current) {
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

func DownloadUpdate() {
	err := http_utils.DownloadFile(updaterInstance.Url, "./config/HYZL-new.exe", false)
	if err != nil {
		print_utils.PrintError(err)
		return
	}
}

func UpdateRightNow() {
	if !UpdateTempExist() {
		if CheckForUpdate() {
			DownloadUpdate()
		} else {
			fmt.Println("当前已经是最新版本")
			return
		}
	}
	// 检查文件MD5是否一致
	bytes, err := os.ReadFile("./config/updater.json")
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	var configUpdater updater
	err = json.Unmarshal(bytes, &configUpdater)
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	configMD5 := configUpdater.MD5
	fileMd5, err := io_utils.CalcMD5("./config/HYZL-new.exe")
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	if configMD5 != fileMd5 {
		fmt.Println("文件MD5不一致，正在重新下载...")
		if CheckForUpdate() {
			DownloadUpdate()
		} else {
			fmt.Println("当前已经是最新版本")
			return
		}
	}
	// 将config里的启动器复制到当前目录
	err = io_utils.MoveFile("./config/HYZL-new.exe", "HYZL-new.exe")
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	global.Config.JustFinishedUpdating = true
	global.WriteConfig()
	generateUpdateBat()
	runUpdateBat()
}

func generateUpdateBat() {
	// 生成更新脚本
	batchContent := fmt.Sprintf(`
@echo off
echo 正在更新启动器
REM 延迟一下，等待启动器关闭
ping 127.0.0.1 -n 3 > nul
REM 检查当前目录下是否存在 HYZL-new.exe
echo 正在检查是否存在更新文件
IF EXIST HYZL-new.exe (
	REM 如果存在，替换掉旧的启动器
	echo 更新文件存在，正在替换启动器
	RENAME "%s" HYZL-old.exe
	RENAME HYZL-new.exe "%s"
	REM 删除旧的启动器
	echo 正在删除旧的启动器
	DEL HYZL-old.exe
	echo 更新完成，正在重新启动启动器
	START "" "%s"
) ELSE (
	REM 如果不存在，打印错误信息
	echo 未找到更新文件，请重新下载
)
`, global.Global.ProgramName, global.Global.ProgramName, global.Global.ProgramName)
	batchContent = strings.ReplaceAll(batchContent, "\n", "\r\n")
	data, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(batchContent))
	os.WriteFile("update.bat", data, os.ModePerm)
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

func ShowChangelog() {
	// 检查是不是刚刚更新完
	if !global.Config.JustFinishedUpdating {
		return
	}
	global.Config.JustFinishedUpdating = false
	global.WriteConfig()
	// 显示更新日志
	print_utils.PrintWithColor(ct.Magenta, true, "更新日志：")
	// 从json中读取更新文件信息，写入到updaterInstance
	bytes, err := os.ReadFile("./config/updater.json")
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	err = json.Unmarshal(bytes, &updaterInstance)
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	fmt.Println(updaterInstance.Changelog)
}
