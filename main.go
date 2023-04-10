// 编译： go build
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

type Config struct {
    GitInstalled    bool `json:"git_installed"`
    NodeJSInstalled bool `json:"nodejs_installed"`
    NpmInstalled    bool `json:"npm_installed"`
}

func checkFirstRun() {
    //检查当前目录下是否存在config文件夹
    _, err := os.Stat("./config")
    //如果不存在就创建
    if err != nil {
        err = os.Mkdir("./config", 0777)
        if err != nil {
            printErr(err)
            return
        }
        //再创建config.json
        file, err := os.Create("./config/config.json")
        if err != nil {
            printErr(err)
            return
        }
        defer file.Close()
        var config Config
        config.GitInstalled = false
        config.NodeJSInstalled = false
        config.NpmInstalled = false
        //写入文件
        data, err := json.MarshalIndent(config, "", "    ")
        if err != nil {
            printErr(err)
            return
        }
        _, err = file.Write(data)
        if err != nil {
            printErr(err)
            return
        }
    }
}

func checkEnv() bool {
    var willWrite = false
    var config Config
    file, err := os.Open("./config/config.json")
    if err == nil {
        defer file.Close()
        decoder := json.NewDecoder(file)
        err = decoder.Decode(&config)
        if err != nil {
            printErr(err)
        }
    }
    if !config.GitInstalled {
        if !checkCommand("git -v") {
            printWithEmptyLine("检测到未安装 Git ，请安装后继续")
            return false
        } else {
            config.GitInstalled = true
            willWrite = true
        }
    }
    if !config.NodeJSInstalled {
        if !checkCommand("node -v") {
            printWithEmptyLine("检测到未安装 Node.js ，请安装后继续")
            return false
        } else {
            config.NodeJSInstalled = true
            willWrite = true

        }
    }
    if !config.NpmInstalled {
        if !checkCommand("npm -v") {
            fmt.Print("检测到未安装 npm ，请手动安装Node.js，具体请看：https://note.youdao.com/s/ImCA210l")
        } else {
            config.NpmInstalled = true
            willWrite = true
        }
    }
    if willWrite {
        //写入到文件
        data, err := json.MarshalIndent(config, "", "    ")
        if err != nil {
            printErr(err)
            return false
        }
        err = os.WriteFile("./config/config.json", data, 0777)
        if err != nil {
            printErr(err)
            return false
        }
    }
    return true
}

func startRedis() *exec.Cmd {
    printWithEmptyLine("正在启动 Redis ...")
    os.Chdir("./redis-windows-7.0.4")
    dir, _ := os.Getwd()
    dir += "\\redis-server.exe"
    printWithEmptyLine(dir)
    cmd := exec.Command("cmd", "/c", "start", "", dir)
    err := cmd.Start()
    if err != nil {
        printErr(err)
    }
    println("Redis 启动成功！")
    os.Chdir("..")
    return cmd
}

func isRedisRunning() bool {
    // 执行 tasklist 命令并获取输出结果
    cmd := exec.Command("tasklist")
    output, err := cmd.Output()
    if err != nil {
        panic(err)
    }

    // 检查输出结果中是否包含 redis-server.exe 进程
    if strings.Contains(string(output), "redis-server.exe") {
        return true
    } else {
        return false
    }
}

//↓菜单函数
func mainMenu() {
    for {
        fmt.Println("===主菜单===")
        fmt.Println("1. 安装云崽")
        fmt.Println("2. 云崽管理")
        fmt.Println("3. BUG修复")
        fmt.Println("0. 退出程序")
        fmt.Print("\n请选择操作：")

        var choice int
        _, err := fmt.Scanln(&choice)
        if err != nil {
            printWithEmptyLine("输入错误，请重新选择")
            continue
        }

        switch choice {
        case 0:
            printWithEmptyLine("退出程序")
            return
        case 1:
            clearLog()
            downloadYunzaiFromGitee()
        case 2:
            clearLog()
            manageYunzaiMenu()
        case 3:
            clearLog()
            bugsFixMenu()
        default:
            printWithEmptyLine("选择不正确，请重新选择")
        }
    }
}

var globalRepositoryLink = `https://gitee.com/bling_yshs/YzLauncher-windows`

const version = "v0.0.7"

var programRunPath = ""

func main() {
    getAppInfo(&programRunPath)
    autoUpdate()
    checkFirstRun()
    if !checkEnv() {
        shutdownApp()
    }
    checkRedis()
    println("当前版本:", version)
    mainMenu()
}
