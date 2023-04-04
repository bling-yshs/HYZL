package main

import (
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
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

//true说明需要更新，false说明不需要更新
func compareRemind(lastRemindTime string) bool {
    if lastRemindTime != "" {
        t, err := time.Parse("2006-01-02", lastRemindTime)
        if err != nil {
            fmt.Println("解析时间出错：", err)
            return false
        }
        t = t.Add(24 * time.Hour)
        // 获取当前时间
        now := time.Now()

        // 比较时间
        if now.Before(t) {
            return false
        } else {
            os.Remove("./config/remindTime.txt")
        }
    }
    return true
}

func autoUpdate() {
    _, latestVersion := getLatestVerion()
    if !compareVersion(version, latestVersion) {
        return
    }
    lastRemindTime := readRemind()
    //如果lastRemindTime的时间加上一天大于当前时间，就不提醒
    if !compareRemind(lastRemindTime) {
        return
    }
    _, err := os.Stat("./update.bat")
    if err == nil {
        //删除update.bat
        os.Remove("./update.bat")
    }
    //得到tempPath
    tempPath := os.Getenv("TEMP")
    //得到md5的值
    md5DownloadedPath := filepath.Join(tempPath, "yzMD5.txt")
    correctMD5, err := getFileContent(md5DownloadedPath)
    if err != nil {
        downloadYz(latestVersion)
        return
    }
    //得到下载的文件的md5值
    YzDownloadedPath := filepath.Join(tempPath, "YzLauncher-windows.exe")
    _, err = os.Stat(YzDownloadedPath)
    if err == nil {
        downloadFileMD5 := getFileMD5(YzDownloadedPath)
        //比较correctMD5与downloadFileMD5
        if !strings.EqualFold(correctMD5, downloadFileMD5) {
            //如果不相等，就删除YzLauncher-windows.exe
            os.Remove(YzDownloadedPath)
            os.Remove(md5DownloadedPath)
            downloadYz(latestVersion)
            return
        }
        _, err = os.Stat(YzDownloadedPath)
        if err == nil {
            fmt.Println("新版本启动器已下载，是否立即更新？(是:y 一天内不提醒:n)")
            userChoice := ReadChoice("y", "n")
            //创建update.bat
            if userChoice == "y" {
                createUpdateBat()
                time.Sleep(1 * time.Second)
                cmd := exec.Command("cmd", "/c", "start", "", filepath.Join(programRunPath, "update.bat"))
                cmd.Start()
                time.Sleep(1 * time.Second)
                os.Exit(0)
            }
            if userChoice == "n" {
                createRemind()
                return
            }
        }
    }

}

func downloadYz(latestVersion string) {
    if compareVersion(version, latestVersion) {
        md5downloadLink := `https://gitee.com/bling_yshs/YzLauncher-windows/releases/download/` + latestVersion + `/yzMD5.txt`
        go downloadFile(md5downloadLink, "")
        downloadLink := `https://gitee.com/bling_yshs/YzLauncher-windows/releases/download/` + latestVersion + `/YzLauncher-windows.exe`
        go downloadFile(downloadLink, "")
    }
}

func createUpdateBat() {
    batchContent := `
@echo off
echo Updating...
ping 127.0.0.1 -n 4 > nul
set launcher=YzLauncher-windows.exe
set md5=yzMD5.txt
set source=%TEMP%
set destination=%CD%

if exist "%source%\%launcher%" (
copy /Y "%source%\%launcher%" "%destination%\%launcher%" > nul
del /Q "%source%\%launcher%"
del /Q "%source%\%md5%"
) 
start %destination%\%launcher%
exit
`

    err := os.WriteFile("update.bat", []byte(batchContent), 0777)
    if err != nil {
        fmt.Println(err)
        return
    }
}

//返回最新版本的下载链接和版本号
func getLatestVerion() (string, string) {
    url := "https://gitee.com/bling_yshs/YzLauncher-windows/releases/latest"

    client := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            // Disable automatic redirect following
            return http.ErrUseLastResponse
        },
    }
    resp, err := client.Get(url)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    newLink := resp.Header.Get("Location")
    segments := strings.Split(newLink, "/")

    // Get the last segment
    return newLink, segments[len(segments)-1]
}
