// 编译： go build -o 云崽启动器.exe main.go
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

func executeCmd(command string, startMsg string, returnMsg string) {
    fmt.Println(startMsg)
    fmt.Println()
    cmd := exec.Command("cmd", "/C", command)
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
    fmt.Println(returnMsg)
    fmt.Println()
}

func printErr(err error) {
    fmt.Println("发生错误，请截图并反馈给作者:", err)
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
        executeCmd("git clone --depth 1 https://gitee.com/bling_yshs/redis-windows-7.0.4", "开始下载 Redis ...", "下载 Redis 成功！")
    }
    if userChoice == "n" {
        fmt.Println("退出程序")
        os.Exit(0)
    }

}
func downloadYunzai() {
    _, err := os.Stat("./Yunzai-bot")
    if err == nil {
        fmt.Println("检测到当前目录下已存在 Yunzai-bot ，请问是否需要重新下载？(是:y 返回菜单:n)")
        userChoice := ReadInput("y", "n")
        if userChoice == "y" {
            //删除文件夹
            err = os.RemoveAll("./Yunzai-bot")
        }
        if userChoice == "n" {
            return
        }
    }
    b := checkCommand("npm -v")
    if b == false {
        fmt.Print("无法使用npm命令，请手动安装Node.js，具体请看：https://note.youdao.com/s/ImCA210l")
    }
    executeCmd("git clone --depth 1 -b main https://gitee.com/yoimiya-kokomi/Yunzai-Bot.git", "开始下载云崽...", "下载云崽成功！")
    //进入Yunzai-Bot文件夹
    os.Chdir("./Yunzai-Bot")
    b2 := checkCommand("pnpm -v")
    if b2 == false {
        executeCmd("npm install pnpm -g --registry=https://registry.npmmirror.com", "开始安装 pnpm ...", "安装 pnpm 成功！")
    }
    executeCmd("pnpm config set registry https://registry.npmmirror.com", "", "")
    executeCmd("pnpm config set puppeteer_download_host=https://registry.npmmirror.com", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
    executeCmd("pnpm install -P", "开始安装云崽依赖", "安装云崽依赖成功！")
    os.Chdir("..")
}
func startRedis() *exec.Cmd {
    fmt.Println("正在启动 Redis ...")
    fmt.Println()

    // 进入 redis-windows-7.0.4 目录
    if err := os.Chdir("./redis-windows-7.0.4"); err != nil {
        panic(err)
    }

    comm, _ := os.Getwd()
    comm += "./redis-server.exe"
    // 启动 Redis 服务器
    cmd := exec.Command(comm)
    if err := cmd.Start(); err != nil {
        panic(err)
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

func startYunzai() {
    if !isRedisRunning() {
        startRedis()
    }
    os.Chdir("./Yunzai-Bot")
    fmt.Println("正在启动云崽...")
    dir, _ := os.Getwd()
    cmd := exec.Command("cmd.exe", "/c", "start /d", dir, "node app")

    fmt.Println(dir)
    if err := cmd.Start(); err != nil {
        panic(err)
    }
    fmt.Println("云崽启动成功！")
    os.Chdir("..")
}

func reInstallDep() {
    os.Chdir("./Yunzai-Bot")
    executeCmd("pnpm config set puppeteer_download_host=https://registry.npmmirror.com", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
    if _, err := os.Stat("./node_modules"); err == nil {
        fmt.Println("检测到当前目录下已存在 node_modules ，请问是否需要重新安装依赖？(是:y 返回菜单:n)")
        userChoice := ReadInput("y", "n")
        if userChoice == "y" {
            executeCmd("pnpm update", "", "")
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
    //读取用户输入的一串字符串
    fmt.Print("请输入命令：")
    var command string
    _, err := fmt.Scanln(&command)
    if err != nil {
        fmt.Println("请输入字符串")
        return
    }
    os.Chdir("./Yunzai-Bot")
    executeCmd(command, "", "")
    os.Chdir("..")
}

func manageYunzai() {
    dir, _ := os.Getwd()
    fmt.Println(dir)
    for {
        fmt.Println("===云崽管理===")
        fmt.Println("1. 启动云崽")
        fmt.Println("2. 切换账号")
        fmt.Println("3. 重装依赖")
        fmt.Println("4. 输入命令")
        fmt.Println("0. 返回上一级")
        fmt.Print("请选择操作：")
        var choice int
        _, err := fmt.Scanln(&choice)
        if err != nil {
            fmt.Println("输入错误，请重新选择")
            continue
        }

        switch choice {
        case 0:
            fmt.Println("退出程序")
            return
        case 1:
            fmt.Println("您选择了启动云崽")
            fmt.Println()
            startYunzai()
        case 2:
            fmt.Println("您选择了切换账号")
            fmt.Println()
        case 3:
            fmt.Println("您选择了重装依赖")
            fmt.Println()
            reInstallDep()
        case 4:
            fmt.Println("您选择了输入命令")
            fmt.Println()
            customCommand()
        default:
            fmt.Println("选择不正确，请重新选择")
        }
    }
}

func menu() {
    dir, _ := os.Getwd()
    fmt.Println(dir)
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
            fmt.Println("输入错误，请重新选择")
            continue
        }

        switch choice {
        case 0:
            fmt.Println("退出程序")
            return
        case 1:
            fmt.Println("您选择了安装云崽")
            fmt.Println()
            downloadYunzai()
        case 2:
            fmt.Println("您选择了云崽管理")
            fmt.Println()
            manageYunzai()
        case 3:
            fmt.Println("您选择了BUG修复")
            fmt.Println()
            // TODO: 执行BUG修复的相关代码
        default:
            fmt.Println("选择不正确，请重新选择")
        }
    }
}

func main() {
    checkFirstRun()
    checkEnv()
    checkRedis()
    menu()
}
