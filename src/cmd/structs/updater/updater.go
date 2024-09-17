package updater

import (
	"encoding/json"
	"fmt"
	"github.com/bling-yshs/HYZL/src/cmd/structs/app"
	"github.com/bling-yshs/HYZL/src/cmd/structs/config"
	"github.com/bling-yshs/HYZL/src/cmd/utils/http_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/io_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
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

var url = app.GetApp().UpdateUrl

// 得到最后一个没有废弃的版本的实例
func GetLatestUpdater() (updater, error) {
	var updaterList []updater
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(url)
	if err != nil {
		return updater{}, errors.Wrap(err, "原因：获取更新文件失败")
	}
	defer response.Body.Close()
	// 解析json
	err = json.NewDecoder(response.Body).Decode(&updaterList)
	if err != nil {
		return updater{}, errors.Wrap(err, "原因：解析更新文件失败")
	}
	// 得到最后一个没有废弃的版本，从前往后遍历
	for _, item := range updaterList {
		if !item.Deprecated {
			return item, nil
		}
	}
	return updater{}, nil
}

func CleanUpdater() {
	// 如果当前目录下存在更新脚本，删除
	if _, err := os.Stat("update.bat"); err == nil {
		os.Remove("update.bat")
	}
}

func DownloadUpdate(url string, showProgress bool) {
	err := http_utils.DownloadFile(url, "./config/HYZL-new.exe", showProgress)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "原因：下载更新文件失败"))
		return
	}
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
`, app.GetApp().Name, app.GetApp().Name, app.GetApp().Name)
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
	if !config.GetConfig().JustFinishedUpdating {
		return
	}
	config.GetConfig().JustFinishedUpdating = false
	config.SaveConfig()
	// 显示更新日志
	print_utils.PrintWithColor(ct.Magenta, true, "更新日志：")
	instance, err := readConfig()
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "原因：读取更新文件失败"))
		return
	}
	fmt.Println(instance.Changelog)
}

// 立即更新启动器
func MenuUpdateRightNow() {
	latestUpdater, err := GetLatestUpdater()
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	latestVersion, err := version.NewVersion(latestUpdater.Version)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "原因：解析最新版本失败"))
		return
	}
	currentVersion, err := version.NewVersion(app.GetApp().Version)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "原因：解析当前版本失败"))
		return
	}
	if !latestVersion.GreaterThan(currentVersion) {
		print_utils.PrintWithEmptyLine("当前版本已经是最新版本")
		return
	}
	DownloadUpdate(latestUpdater.Url, true)
	err = writeConfig(latestUpdater)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "原因：写入更新文件失败"))
		return
	}
	// 将config里的启动器复制到当前目录
	err = io_utils.MoveFile("./config/HYZL-new.exe", "HYZL-new.exe")
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "原因：移动文件失败"))
		return
	}
	config.GetConfig().JustFinishedUpdating = true
	config.GetConfig().HaveUpdate = false
	config.SaveConfig()
	generateUpdateBat()
	runUpdateBat()
}

// 从json中读取更新文件信息
func readConfig() (updater, error) {
	// 读取本地配置
	bytes, err := os.ReadFile("./config/updater.json")
	if err != nil {
		return updater{}, err
	}
	var instance updater
	err = json.Unmarshal(bytes, &instance)
	if err != nil {
		return updater{}, err
	}
	return instance, nil
}

func writeConfig(instance updater) error {
	// 解析并写入到本地配置
	bytes, err := json.Marshal(instance)
	if err != nil {
		return err
	}
	err = os.WriteFile("./config/updater.json", bytes, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
