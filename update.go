package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func createRemind() {
	_, err := os.Stat("./config/remindTime.txt")
	if err != nil {
		file, err := os.Create("./config/remindTime.txt")
		if err != nil {
			printErr(err)
			return
		}
		defer file.Close()

		//写入当前时间
		_, err = file.WriteString(time.Now().Format("2006-01-02"))
		if err != nil {
			printErr(err)
			return
		}
	}
}

func readRemind() string {
	//读取remindTime.txt
	_, err2 := os.Stat("./config/remindTime.txt")
	if err2 != nil {
		return ""
	}
	file, err := os.Open("./config/remindTime.txt")
	if err != nil {
		return ""
	}
	defer file.Close()
	//返回文件中的内容
	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		return ""
	}
	return string(buf[:n])
}

// true说明需要更新，false说明不需要更新
func compareRemind(lastRemindTime string) bool {
	if lastRemindTime != "" {
		t, err := time.Parse("2006-01-02", lastRemindTime)
		if err != nil {
			printWithEmptyLine("解析时间出错：" + err.Error())
			return false
		}
		t = t.Add(24 * time.Hour)
		// 获取当前时间
		now := time.Now()

		// 比较时间
		if now.Before(t) {
			return false
		} else {
			_ = os.Remove("./config/remindTime.txt")
		}
	}
	return true
}

func isNewYunzai() bool {
	//得到md5的值
	md5DownloadedPath := filepath.Join(config.SystemTempPath, "yzMD5.txt")
	_, err2 := os.Stat(md5DownloadedPath)
	if err2 != nil {
		return false
	}
	correctMD5, err := getFileContent(md5DownloadedPath)
	if err != nil {
		return false
	}
	//得到下载的文件的md5值
	YzDownloadedPath := filepath.Join(config.SystemTempPath, "YzLauncher-windows.exe")
	_, err = os.Stat(YzDownloadedPath)
	if err != nil {
		return false
	}
	downloadYunzaiMD5 := getFileMD5(YzDownloadedPath)
	if !strings.EqualFold(correctMD5, downloadYunzaiMD5) {
		return false
	}
	//最后检验一下是否是最新的版本
	latestVersion := giteeAPI.getLatestTag()
	if !compareVersion(version, latestVersion) {
		os.Remove(md5DownloadedPath)
		os.Remove(YzDownloadedPath)
		return false
	}
	return true
}

func update() {
	latestVersion := giteeAPI.getLatestTag()
	if compareVersion(updatedVersion, latestVersion) {
		downloadLauncher(latestVersion)
	}
}

func updateThisProgram() {
	_, err := os.Stat("./update.bat")
	if err == nil {
		//删除update.bat
		_ = os.Remove("./update.bat")
		//显示更新日志
		_, err := os.Stat("./config/changelog.txt")
		if err == nil {
			content, err := getFileContent("./config/changelog.txt")
			if err != nil {
				printErr(err)
				return
			}
			printWithEmptyLine("新版本更新内容：\n" + content)
			//删除changelog.txt
			_ = os.Remove("./config/changelog.txt")
		}
	}

	lastRemindTime := readRemind()
	if lastRemindTime != "" {
		//如果lastRemindTime的时间加上一天大于当前时间，就不提醒
		if !compareRemind(lastRemindTime) {
			return
		}
	}
	if isNewYunzai() {
		printWithEmptyLine("新版本启动器已下载，是否立即更新？(是:y 一天内不提醒:n)")
		userChoice := ReadChoice("y", "n")
		if userChoice == "y" {
			createUpdateBat()
			time.Sleep(1 * time.Second)
			cmd := exec.Command("cmd", "/c", "start", "", filepath.Join(programRunPath, "update.bat"))
			_ = cmd.Start()
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}
		if userChoice == "n" {
			createRemind()
			return
		}
	}
}

func downloadLauncher(latestVersion string) {
	updating = true
	if compareVersion(updatedVersion, latestVersion) {
		md5downloadLink, _ := url.JoinPath(globalRepositoryLink, "releases", "download", latestVersion, "yzMD5.txt")
		downloadFile(md5downloadLink, "")
		createChangelog()
		downloadLink, _ := url.JoinPath(globalRepositoryLink, "releases", "download", latestVersion, "YzLauncher-windows.exe")
		downloadFile(downloadLink, "")
		updatedVersion = latestVersion
	}
}

func createChangelog() {
	changelog := giteeAPI.getBody()
	_ = os.WriteFile("./config/changelog.txt", []byte(changelog), os.ModePerm)
}

func createUpdateBat() {
	batchContent := fmt.Sprintf(`
@echo off
echo 正在更新...
ping 127.0.0.1 -n 4 > nul
set launcher=YzLauncher-windows.exe
set md5=yzMD5.txt
set source=%%TEMP%%
set destination=%%CD%%

if exist "%%source%%\%%launcher%%" (
copy /Y "%%source%%\%%launcher%%" "%%destination%%\%s" > nul
del /Q "%%source%%\%%launcher%%"
del /Q "%%source%%\%%md5%%"
) 
start %%destination%%\%s
exit`, programName, programName)

	data1, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(batchContent))
	_ = os.WriteFile(`temp.bat`, data1, os.ModePerm)
	executeCmd(`type temp.bat | find "" /V > update.bat`)
	_ = os.RemoveAll(`temp.bat`)
}

func updateLauncherRightNow() {
	latestVersion := giteeAPI.getLatestTag()
	if !compareVersion(version, latestVersion) {
		printWithEmptyLine("当前已是最新版本")
		return
	}
	printWithEmptyLine("正在下载更新...")
	if updating == false {
		downloadLauncher(latestVersion)
	}
	for {
		if isNewYunzai() {
			break
		}
		time.Sleep(5 * time.Second)
	}
	createUpdateBat()
	time.Sleep(1 * time.Second)
	cmd := exec.Command("cmd", "/c", "start", "", filepath.Join(programRunPath, "update.bat"))
	_ = cmd.Start()
	time.Sleep(1 * time.Second)
	os.Exit(0)
}
