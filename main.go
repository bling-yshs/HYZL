// 编译： go build
package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/fs"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

type Config struct {
    GitInstalled    bool `json:"git_installed"`
    NodeJSInstalled bool `json:"nodejs_installed"`
    NpmInstalled    bool `json:"npm_installed"`
}

//↓工具函数

func getAppInfo(args ...*string) {
    //获取程序运行路径
    execPath, err := os.Executable()
    if err != nil {
        printErr(err)
    }
    currentDir := filepath.Dir(execPath)
    *args[0] = currentDir
}

func clearLog() {
    //执行clear指令清除控制台
    cmd := exec.Command("cmd", "/c", "cls")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func compareVersion(version string, latestVersion string) bool {
    version = version[1:]                   // 去除前面的v
    v1 := strings.Split(version, ".")       // 将版本号按 "." 分割成数组
    v2 := strings.Split(latestVersion, ".") // 同上

    for i := 0; i < len(v1) || i < len(v2); i++ {
        n1, n2 := 0, 0 // 初始化数字变量

        // 如果第一个版本号没有到头，就将其转换为数字
        if i < len(v1) {
            n1, _ = strconv.Atoi(v1[i])
        }

        // 如果第二个版本号没有到头，就将其转换为数字
        if i < len(v2) {
            n2, _ = strconv.Atoi(v2[i])
        }

        // 比较数字大小，输出结果
        if n2 > n1 {
            return true
        }
    }
    return false
}

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
        if len(stringArgs[1]) > 0 {
            printWithEmptyLine(stringArgs[1])
        }
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

func shutdownApp() {
    fmt.Println("按回车键退出...")
    fmt.Scanln()
    os.Exit(0)
}

//↑工具函数

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

func update() {
    batPath := filepath.Join(programRunPath, "update.bat")
    cmd := exec.Command("cmd", "/c", "start", "", batPath)
    cmd.Start()
    time.Sleep(1 * time.Second)
    os.Exit(0)
}

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
    dir += "\\redis-server.exe"
    printWithEmptyLine(dir)
    cmd := exec.Command("cmd", "/c", "start", "", dir)
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
        //等待1秒
        time.Sleep(1 * time.Second)
    }
    os.Chdir("./Yunzai-Bot")
    printWithEmptyLine("正在启动云崽...")
    dir, _ := os.Getwd()
    cmd := exec.Command("cmd", "/C", "start", "/d", dir, "cmd", "/k", "node app")
    cmd.Start()
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
    // 读取文件内容
    content, err := os.ReadFile("./Yunzai-Bot/config/config/qq.yaml")
    if err != nil {
        panic(err)
    }

    // 将文件内容转换为字符串
    strContent := string(content)

    // 读取用户输入的 qq、pwd 和 platform
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("请输入 QQ 账号：")
    scanner.Scan()
    qq := scanner.Text()
    fmt.Print("请输入密码：")
    scanner.Scan()
    pwd := scanner.Text()
    fmt.Print("请输入登录方式（1:安卓手机、2:aPad、3:安卓手表、4:MacOS、5:iPad）2023年3月31日：推荐使用MacOS登录：")
    scanner.Scan()
    platform := scanner.Text()

    // 替换文件中的 qq、pwd 和 platform 字段
    lines := strings.Split(strContent, "\n")
    for i, line := range lines {
        if strings.HasPrefix(line, "qq:") {
            lines[i] = fmt.Sprintf("qq: %s", qq)
        } else if strings.HasPrefix(line, "pwd:") {
            lines[i] = fmt.Sprintf("pwd: '%s'", pwd)
        } else if strings.HasPrefix(line, "platform:") {
            lines[i] = fmt.Sprintf("platform: %s", platform)
        }
    }
    newContent := strings.Join(lines, "\n")

    // 将更新后的配置写回文件
    err = os.WriteFile("./Yunzai-Bot/config/config/qq.yaml", []byte(newContent), fs.FileMode(0777))
    if err != nil {
        panic(err)
    }

    fmt.Println("云崽账号更新成功！")
}

func pupFix() {
    os.Chdir("./Yunzai-Bot")
    executeCmd("pnpm install puppeteer@19.7.3 -w", "正在修复 puppeteer...")
    executeCmd("node ./node_modules/puppeteer/install.js", "正在下载 Chromium...")
    os.Chdir("..")
}

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

func installGuobaPlugin() {
    installPluginsTemplate("锅巴插件", "Guoba-Plugin", "git clone --depth=1 https://gitee.com/guoba-yunzai/guoba-plugin.git ./plugins/Guoba-Plugin/", "pnpm install --no-lockfile --filter=guoba-plugin -w")
}

func installMiaoPlugin() {
    installPluginsTemplate("喵喵插件", "miao-plugin", "git clone --depth 1 -b master https://gitee.com/yoimiya-kokomi/miao-plugin.git ./plugins/miao-plugin/", "pnpm add image-size -w")
}

func installXiaoyaoPlugin() {
    installPluginsTemplate("逍遥插件", "miao-plugin", "git clone --depth=1 https://gitee.com/Ctrlcvs/xiaoyao-cvs-plugin.git ./plugins/xiaoyao-cvs-plugin/ ./plugins/miao-plugin/", "pnpm add promise-retry -w", "pnpm add superagent -w")
}

func installPluginsTemplate(pluginChineseName string, dirName string, command ...string) {
    pluginDir := "./plugins/" + dirName
    _, err := os.Stat(pluginDir)
    if err == nil {
        fmt.Println("当前已安装 ", pluginChineseName, "，请问是否需要重新安装？(是:y 返回菜单:n)")
        userChoice := ReadInput("y", "n")
        if userChoice == "n" {
            return
        }
    }
    for _, cmd := range command {
        executeCmd(cmd)
    }
}

func installPluginsMenu() {
    os.Chdir("./Yunzai-Bot")
    for {
        fmt.Println("===安装插件===")
        fmt.Println("1. 锅巴插件")
        fmt.Println("2. 喵喵插件")
        fmt.Println("3. 逍遥插件")
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
            os.Chdir("..")
            return
        case 1:
            clearLog()
            installGuobaPlugin()
        case 2:
            clearLog()
            installMiaoPlugin()
        case 3:
            clearLog()
            installXiaoyaoPlugin()
        default:
            printWithEmptyLine("选择不正确，请重新选择")
        }
    }
}

func manageYunzaiMenu() {

    for {
        fmt.Println("===云崽管理===")
        fmt.Println("1. 启动云崽")
        fmt.Println("2. 强制关闭云崽")
        fmt.Println("3. 切换账号")
        fmt.Println("4. 安装插件")
        fmt.Println("5. 自定义终端命令")
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
            startYunzai()
        case 2:
            clearLog()
            closeYunzai()
        case 3:
            clearLog()
            changeAccount()
        case 4:
            clearLog()
            installPluginsMenu()
        case 5:
            clearLog()
            customCommand()
        default:
            printWithEmptyLine("选择不正确，请重新选择")
        }
    }
}

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
            downloadYunzai()
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

func checkUpdate() {
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

const version = "v0.0.4"

var programRunPath = ""

func main() {
    getAppInfo(&programRunPath)
    checkUpdate()
    checkFirstRun()
    if !checkEnv() {
        shutdownApp()
    }
    checkRedis()
    mainMenu()
}
