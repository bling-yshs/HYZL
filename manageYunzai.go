package main

import (
    "bufio"
    "fmt"
    "io/fs"
    "os"
    "os/exec"
    "strings"
    "time"
)

func manageYunzaiMenu() {

    for {
        fmt.Println("===云崽管理===")
        fmt.Println("1. 启动云崽")
        fmt.Println("2. 强制关闭云崽")
        fmt.Println("3. 切换账号")
        fmt.Println("4. 安装插件")
        fmt.Println("5. 安装js插件")
        fmt.Println("6. 自定义终端命令")
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
            installJsPlugin()
        case 6:
            clearLog()
            customCommand()
        default:
            printWithEmptyLine("选择不正确，请重新选择")
        }
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
func installJsPlugin() {
    jsPluginDir := programRunPath + "/Yunzai-bot/plugins/example"
    //输入js插件的地址，例如https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/raw/master/%E5%96%9C%E6%8A%A5.js
    fmt.Print("请输入js插件的地址：")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    jsPluginUrl := scanner.Text()
    //检查url是否为https://开头，并且以js结尾
    if !strings.HasPrefix(jsPluginUrl, "https://") || !strings.HasSuffix(jsPluginUrl, ".js") {
        fmt.Println("输入的js插件地址不正确，请重新输入")
        return
    }
    //如果输入格式是https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/blob/master/%E5%96%9C%E6%8A%A5.js则自动转换为https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/raw/master/%E5%96%9C%E6%8A%A5.js
    jsPluginUrl = strings.Replace(jsPluginUrl, "blob", "raw", 1)
    //下载js插件，保存到jsPluginDir
    downloadFileSync(jsPluginDir, jsPluginUrl)
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