package main

import (
    "fmt"
    "os"
)

func bugsFixMenu() {
    for {
        fmt.Println("===BUG修复===")
        fmt.Println("1. 重装依赖")
        fmt.Println("2. 修复 puppeteer Chromium 问题")
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
            pupFix()
        default:
            printWithEmptyLine("选择不正确，请重新选择")
        }
    }
}

func reInstallDep() {
    os.Chdir("./Yunzai-Bot")
    executeCmd("pnpm config set puppeteer_download_host=https://registry.npmmirror.com", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
    if _, err := os.Stat("./node_modules"); err == nil {
        fmt.Println("检测到当前目录下已存在 node_modules ，请问是否需要重新安装依赖？(是:y 返回菜单:n)")
        userChoice := ReadChoice("y", "n")
        if userChoice == "y" {
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

func pupFix() {
    os.Chdir("./Yunzai-Bot")
    executeCmd("pnpm install puppeteer@19.7.3 -w", "正在修复 puppeteer...")
    executeCmd("node ./node_modules/puppeteer/install.js", "正在下载 Chromium...")
    os.Chdir("..")
}
