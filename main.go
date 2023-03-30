// 编译： go build -o 云崽启动器.exe main.go
package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "time"
)

type Config struct {
    GitInstalled    bool `json:"git_installed"`
    NodeJSInstalled bool `json:"nodejs_installed"`
}

//↓工具函数
func checkCommand(command string) bool {
    cmd := exec.Command("cmd", "/c", command)
    err := cmd.Run()
    if err == nil {
        return true
    } else {
        return false
    }
}
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

func executeCmd(stringArgs ...string) {

    cmd := exec.Command("cmd", "/C", stringArgs[0])
    if len(stringArgs) >= 2 {
        printWithEmptyLine(stringArgs[1])
    }
    cmd.Stdout = os.Stdout // 直接将命令标准输出连接到标准输出流
    cmd.Stderr = os.Stderr // 将错误输出连接到标准错误流
    cmd.Stdin = os.Stdin   // 将标准输入连接到命令的标准输入

    err := cmd.Start()
    if err != nil {
        printErr(err)
    }

    err = cmd.Wait()
    if err != nil {
        printErr(err)
    }
    if len(stringArgs) >= 3 {
        printWithEmptyLine(stringArgs[2])
    }
}

func printErr(err error) {
    fmt.Println("发生错误，请截图并反馈给作者:", err)
}

func printWithEmptyLine(str string) {
    fmt.Println()
    fmt.Println(str)
    fmt.Println()
}

//↑工具函数

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

func checkEnv() bool {
    b := checkCommand("git -v")
    if !b {
        printWithEmptyLine("检测到未安装 Git ，请安装后继续")
        return false
    }
    b2 := checkCommand("node -v")
    if !b2 {
        printWithEmptyLine("检测到未安装 Node.js ，请安装后继续")
        return false
    }

    b3 := checkCommand("npm -v")
    if !b3 {
        fmt.Print("检测到未安装 npm ，请手动安装Node.js，具体请看：https://note.youdao.com/s/ImCA210l")
    }
    return true
}

func checkRedis() {
    _, err := os.Stat("./redis-windows-7.0.4")
    if err == nil {
        return
    }
    printWithEmptyLine("检测到当前目录下不存在 redis-windows-7.0.4 ，请问是否需要自动下载 Redis ？(是:y 退出程序:n)")
    //读取用户输入y或者n
    userChoice := ReadInput("y", "n")
    if userChoice == "y" {
        executeCmd("git clone --depth 1 https://gitee.com/bling_yshs/redis-windows-7.0.4", "开始下载 Redis ...", "下载 Redis 成功！")
    }
    if userChoice == "n" {
        printWithEmptyLine("退出程序")
        os.Exit(0)
    }
}

func downloadYunzai() {
    _, err := os.Stat("./Yunzai-bot")
    if err == nil {
        printWithEmptyLine("检测到当前目录下已存在 Yunzai-bot ，请问是否需要重新下载？(是:y 返回菜单:n)")
        userChoice := ReadInput("y", "n")
        if userChoice == "y" {
            //删除文件夹
            os.RemoveAll("./Yunzai-bot")
        }
        if userChoice == "n" {
            return
        }
    }
    executeCmd("git clone --depth 1 -b main https://gitee.com/yoimiya-kokomi/Yunzai-Bot.git", "开始下载云崽...", "下载云崽成功！")
    //进入Yunzai-Bot文件夹
    os.Chdir("./Yunzai-Bot")
    b2 := checkCommand("pnpm -v")
    if !b2 {
        executeCmd("npm install pnpm -g --registry=https://registry.npmmirror.com", "开始安装 pnpm ...", "安装 pnpm 成功！")
    }
    executeCmd("pnpm config set registry https://registry.npmmirror.com")
    executeCmd("pnpm config set puppeteer_download_host=https://registry.npmmirror.com", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
    executeCmd("pnpm install -P", "开始安装云崽依赖", "安装云崽依赖成功！")
    os.Chdir("..")
}

func startRedis() *exec.Cmd {
    printWithEmptyLine("正在启动 Redis ...")
    os.Chdir("./redis-windows-7.0.4")
    dir, _ := os.Getwd()
    dir += "/redis-server.exe"
    cmd := exec.Command("cmd", "/c", "start", dir)
    err := cmd.Start()
    fmt.Println(err)
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

func startYunzai() {
    if !isRedisRunning() {
        startRedis()
        //等待3秒
        time.Sleep(3 * time.Second)
    }
    os.Chdir("./Yunzai-Bot")
    printWithEmptyLine("正在启动云崽...")
    dir, _ := os.Getwd()
    cmd := exec.Command("cmd", "/C", "start", "/d", dir, "cmd", "/k", "node app")
    cmd.Run()
    printWithEmptyLine("云崽启动成功！")
    os.Chdir("..")
}

func reInstallDep() {
    os.Chdir("./Yunzai-Bot")
    executeCmd("pnpm config set puppeteer_download_host=https://registry.npmmirror.com", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
    if _, err := os.Stat("./node_modules"); err == nil {
        fmt.Println("检测到当前目录下已存在 node_modules ，请问是否需要重新安装依赖？(是:y 返回菜单:n)")
        userChoice := ReadInput("y", "n")
        if userChoice == "y" {
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

func customCommand() {
    // 读取用户输入的一串字符串
    fmt.Print("请输入命令：")
    reader := bufio.NewReader(os.Stdin)
    command, _ := reader.ReadString('\n')
    command = strings.TrimSuffix(command, "\n")

    os.Chdir("./Yunzai-Bot")
    executeCmd(command)
    os.Chdir("..")
}

func closeYunzai() {
    exec.Command("taskkill", "/FI", "WINDOWTITLE eq Yunzai-bot", "/T", "/F").Run()
    executeCmd("taskkill /f /im node.exe", "正在关闭云崽...", "云崽关闭成功！")
}

func changeAccount() {
}

func clearLog() {
    //执行clear指令清除控制台
    cmd := exec.Command("cmd", "/c", "cls")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func manageYunzai() {

    for {
        fmt.Println("===云崽管理===")
        fmt.Println("1. 启动云崽")
        fmt.Println("2. 强制关闭云崽")
        fmt.Println("3. 切换账号")
        fmt.Println("4. 重装依赖")
        fmt.Println("5. 自定义终端命令")
        fmt.Println("0. 返回上一级")
        fmt.Print("请选择操作：")
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
            startYunzai()
        case 2:
            clearLog()
            closeYunzai()
        case 3:
            clearLog()
            changeAccount()
        case 4:
            clearLog()
            reInstallDep()
        case 5:
            clearLog()
            customCommand()
        default:
            printWithEmptyLine("选择不正确，请重新选择")
        }
    }
}

func menu() {
    for {
        fmt.Println("===主菜单===")
        fmt.Println("1. 安装云崽")
        fmt.Println("2. 云崽管理")
        fmt.Println("3. BUG修复")
        fmt.Println("0. 退出程序")
        fmt.Print("请选择操作：")

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
            downloadYunzai()
        case 2:
            clearLog()
            manageYunzai()
        case 3:
            clearLog()

            // TODO: 执行BUG修复的相关代码
        default:
            printWithEmptyLine("选择不正确，请重新选择")
        }
    }
}

func main() {
    checkFirstRun()
    if !checkEnv() {
        //按任意键退出
        fmt.Println("按回车键退出...")
        fmt.Scanln()
        return
    }
    checkRedis()
    menu()
}
