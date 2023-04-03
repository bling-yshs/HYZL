package main

import (
    "bufio"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
)

//↓工具函数

func downloadFile(url string, outputPath string) error {
    if outputPath == "" {
        outputPath = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Temp", filepath.Base(url))
    }

    // 创建一个文件来存储下载的内容
    output, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("Error creating file: %v", err)
    }
    defer output.Close()

    // 发送 GET 请求并将响应写入文件
    client := http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return fmt.Errorf("Error creating request: %v", err)
    }

    // 使用 go 协程异步执行请求
    respChan := make(chan *http.Response)
    errChan := make(chan error)
    go func() {
        resp, err := client.Do(req)
        if err != nil {
            errChan <- fmt.Errorf("Error executing request: %v", err)
            return
        }
        respChan <- resp
    }()

    // 通过 select 语句等待下载完成或出错
    for {
        select {
        case resp := <-respChan:
            defer resp.Body.Close()
            _, err := io.Copy(output, resp.Body)
            if err != nil {
                return fmt.Errorf("Error writing to file: %v", err)
            } else {
                fmt.Printf("Downloaded to %s\n", outputPath)
                return nil
            }
        case err := <-errChan:
            return fmt.Errorf("Error downloading file: %v", err)
        }
    }
}

func downloadFileSync(dir string, fileUrl string) {
    resp, err := http.Get(fileUrl)
    if err != nil {
        printErr(err)
    }
    defer resp.Body.Close()

    fileName, _ := url.PathUnescape(filepath.Base(fileUrl))
    filePath := filepath.Join(dir, fileName)

    out, err := os.Create(filePath)
    if err != nil {
        printErr(err)
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        printErr(err)
    }
}

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

func ReadChoice(allowedValues ...string) string {
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
func checkRedis() {
    _, err := os.Stat("./redis-windows-7.0.4")
    if err == nil {
        return
    }
    printWithEmptyLine("检测到当前目录下不存在 redis-windows-7.0.4 ，请问是否需要自动下载 Redis ？(是:y 退出程序:n)")
    //读取用户输入y或者n
    userChoice := ReadChoice("y", "n")
    if userChoice == "y" {
        executeCmd("git clone --depth 1 https://gitee.com/bling_yshs/redis-windows-7.0.4", "开始下载 Redis ...", "下载 Redis 成功！")
    }
    if userChoice == "n" {
        printWithEmptyLine("退出程序")
        os.Exit(0)
    }
}
