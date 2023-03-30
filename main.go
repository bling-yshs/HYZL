package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

type Config struct {
    GitInstalled    bool `json:"git_installed"`
    NodeJSInstalled bool `json:"nodejs_installed"`
}

//工具函数
func ReadInput(allowedValues ...string) string {
    allowedSet := make(map[string]bool)
    for _, val := range allowedValues {
        allowedSet[val] = true
    }
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("请输入你的选择: ")
        scanner.Scan()
        userInput := strings.TrimSpace(scanner.Text())
        if allowedSet[userInput] {
            return userInput
        }
        fmt.Println("输入有误，请输入：", strings.Join(allowedValues, " 或者 "))
    }
}

func RunCommand(commandStr string) error {
    cmd := exec.Command("cmd", "/C", commandStr)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return err
    }
    stderr, err := cmd.StderrPipe()
    if err != nil {
        return err
    }

    if err := cmd.Start(); err != nil {
        return err
    }

    scannerOut := bufio.NewScanner(stdout)
    scannerErr := bufio.NewScanner(stderr)

    go func() {
        for scannerOut.Scan() {
            fmt.Println(scannerOut.Text())
        }
    }()
    go func() {
        for scannerErr.Scan() {
            fmt.Println(scannerErr.Text())
        }
    }()

    if err := cmd.Wait(); err != nil {
        return err
    }

    return nil
}

func printErr(err error) {
    fmt.Println("发生错误，请截图并反馈给作者:", err)
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
    }

}

func checkEnv() {

    _, err := exec.LookPath("git")
    if err != nil {
        fmt.Println("检测到未安装 Git ，请安装后继续")
        return
    }

    _, err = exec.LookPath("node")
    if err != nil {
        fmt.Println("检测到未安装 Node.js ，请安装后继续")
        return
    }
}

func checkRedis() {
    _, err := os.Stat("./redis-windows-7.0.4")
    if err == nil {
        return
    }
    fmt.Println("检测到当前目录下不存在 redis-windows-7.0.4 ，请问是否需要自动下载 Redis ？(是:y 退出程序:n)")
    //读取用户输入y或者n
    userChoice := ReadInput("y", "n")
    if userChoice == "y" {
        cmd := exec.Command("git", "clone", "--depth", "1", "https://gitee.com/bling_yshs/redis-windows-7.0.4")
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        err := cmd.Run()
        if err != nil {
            fmt.Println("下载 Redis 失败:", err)
            return
        }
        fmt.Println("下载 Redis 成功！")
    }
    if userChoice == "n" {
        fmt.Println("退出程序")
    }

}
func downloadYunzai() {

}

func menu() {
    for {
        fmt.Println("请选择操作：")
        fmt.Println("1. 安装云崽")
        fmt.Println("2. 云崽管理")
        fmt.Println("3. BUG修复")
        fmt.Println("0. 退出程序")

        var choice int
        _, err := fmt.Scanln(&choice)
        if err != nil {
            fmt.Println("输入错误，请重新选择。")
            continue
        }

        switch choice {
        case 0:
            fmt.Println("退出程序")
            return
        case 1:
            fmt.Println("您选择了安装云崽。")
            // TODO: 执行安装云崽的相关代码
        case 2:
            fmt.Println("您选择了云崽管理。")
            // TODO: 执行云崽管理的相关代码
        case 3:
            fmt.Println("您选择了BUG修复。")
            // TODO: 执行BUG修复的相关代码
        default:
            fmt.Println("选择不正确，请重新选择。")
        }
    }
}

func main() {
    checkFirstRun()
    checkEnv()
    checkRedis()
    menu()
}
