package main

import (
    "fmt"
    "os"
    "regexp"
)

func bugsFixMenu() {
    for {
        fmt.Println("===BUG修复===")
        fmt.Println("1. 重装依赖")
        fmt.Println("2. 修复 puppeteer Chromium 启动失败")
        fmt.Println("3. 修复 puppeteer Chromium 弹出cmd窗口(Windows Server 2012请勿使用)")
        fmt.Println("4. 修复 错误码45 错误码238")
        fmt.Println("0. 返回上一级")
        fmt.Print("\n请选择操作：")
        var choice int
        _, err := fmt.Scanln(&choice)
        if err != nil {
            printWithEmptyLine("输入错误，请重新选择")
            continue
        }

        switch choice {
        case 0:
            clearLog()
            return
        case 1:
            clearLog()
            reInstallDep()
        case 2:
            clearLog()
            pupCanNotStartFix()
        case 3:
            clearLog()
            pupPopFix()
        case 4:
            clearLog()
            errorCodeFix()
        default:
            printWithEmptyLine("选择不正确，请重新选择")
        }
    }
}

func errorCodeFix() {
    printWithEmptyLine("开始修复 错误码45 错误码238...")
    os.Chdir("./Yunzai-Bot")
    //读取./config/config/qq.yaml
    s, err := getFileContent("./config/config/qq.yaml")
    if err != nil {
        printErr(err)
        return
    }
    regex := regexp.MustCompile(`platform: \d`)
    s = regex.ReplaceAllString(s, `platform: 4`)
    //写入./config/config/qq.yaml
    err = os.WriteFile("./config/config/qq.yaml", []byte(s), 0777)
    if err != nil {
        printErr(err)
        return
    }
    printWithEmptyLine("修复成功！")
}

func pupPopFix() {
    os.Chdir("./Yunzai-Bot")
    executeCmd("git reset --hard origin/main")
    executeCmd("git pull", "正在更新云崽到最新版本...", "更新云崽到最新版本成功！")
    executeCmd("pnpm uninstall puppeteer", "正在修复 puppeteer...")
    executeCmd("pnpm install puppeteer@19.8.3 -w")
    executeCmd("node ./node_modules/puppeteer/install.js")
    os.Chdir("..")
}

func reInstallDep() {
    os.Chdir("./Yunzai-Bot")
    if _, err := os.Stat("./node_modules"); err == nil {
        fmt.Println("检测到当前目录下已存在 node_modules ，请问是否需要重新安装依赖？(是:y 返回菜单:n)")
        userChoice := ReadChoice("y", "n")
        if userChoice == "y" {
            executeCmd("pnpm config set puppeteer_download_host=https://registry.npmmirror.com", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
            os.RemoveAll("./node_modules")
            executeCmd("pnpm update", "开始安装云崽依赖...")
            executeCmd("pnpm install -P", "", "安装云崽依赖成功！")
        }
        if userChoice == "n" {
            return
        }
    } else {
        executeCmd("pnpm install -P", "", "安装云崽依赖成功！")
    }
    os.Chdir("..")
}

func pupCanNotStartFix() {
    os.Chdir("./Yunzai-Bot")
    executeCmd("pnpm uninstall puppeteer", "正在修复 puppeteer...")
    executeCmd("pnpm install puppeteer@19.7.3 -w")
    executeCmd("node ./node_modules/puppeteer/install.js", "正在下载 Chromium...")
    os.Chdir("..")
}
