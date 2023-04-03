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

func update() {
    batPath := filepath.Join(programRunPath, "update.bat")
    cmd := exec.Command("cmd", "/c", "start", "", batPath)
    cmd.Start()
    time.Sleep(1 * time.Second)
    os.Exit(0)
}

func autoUpdate() {
    _, err := os.Stat("./update.bat")
    if err == nil {
        //删除update.bat
        os.Remove("./update.bat")
    }
    _, latestVersion := getLatestVerion()
    if compareVersion(version, latestVersion) {
        fmt.Println("发现新版本：", latestVersion, "，3 秒后开始更新...")
        time.Sleep(3 * time.Second)
        batPath := filepath.Join(programRunPath, "update.bat")
        downloadLink := `https://gitee.com/bling_yshs/YzLauncher-windows/releases/download/` + latestVersion + `/YzLauncher-windows.exe`
        printWithEmptyLine(downloadLink)
        createUpdateBat(downloadLink, batPath)
        update()
        shutdownApp()
    }
}

func createUpdateBat(latestUrl string, batPath string) {
    batchContent := `@echo off
setlocal enabledelayedexpansion

set "url=` + latestUrl + `"
set "filename=YzLauncher-windows.exe"

curl -L -o "%filename%" "%url%"

if exist "%filename%" (
    move /y "%filename%" ".\%filename%"
    start "" ".\%filename%"
) else (
    echo Failed to download %filename%
)`

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
