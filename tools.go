package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// ↓工具函数
func readInt() int64 {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		s := scanner.Text()
		if s == "" {
			continue
		}
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			fmt.Println("输入错误，请重新输入")
			continue
		}
		return i
	}
}

func readString() string {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		s := scanner.Text()
		if s == "" {
			fmt.Println("输入错误，请重新输入")
			continue
		}
		return s
	}
}

func getFileContent(filePath string) (string, error) {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read file error: %v", err)
	}

	// 返回字符串格式的文件内容
	return string(content), nil
}

func getFileMD5(fPath string) string {
	// 读取文件内容
	data, err := os.ReadFile(fPath)
	if err != nil {
		printErr(err)
	}

	// 计算MD5值
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)

	// 将MD5值转换为字符串格式
	return hex.EncodeToString(cipherStr)
}

func downloadFile(downloadURL string, downloadFilePath string) {
	res, err := http.Get(downloadURL)
	if err != nil {
		printWithEmptyLine("下载失败，错误信息为：" + err.Error())
		return
	}
	var savePath = downloadFilePath

	if downloadFilePath == "" {
		downloadFilePath := os.Getenv("TEMP")
		savePath = filepath.Join(downloadFilePath, filepath.Base(downloadURL))
		if downloadFilePath == "" {
			printWithEmptyLine("无法获取到用户目录")
		}
	}

	f, err := os.Create(savePath)
	if err != nil {
		printWithEmptyLine("创建文件失败，错误信息为：" + err.Error())
		return
	}
	defer f.Close()

	_, _ = io.Copy(f, res.Body)
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
	_ = cmd.Run()
}

// 比较版本号，如果需要更新返回true，否则返回false
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
	cmd := exec.Command("cmd.exe")
	cmd.Stdout = os.Stdout // 直接将命令标准输出连接到标准输出流
	cmd.Stderr = os.Stderr // 将错误输出连接到标准错误流
	cmd.Stdin = os.Stdin   // 将标准输入连接到命令的标准输入
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c %s`, stringArgs[0]), HideWindow: true}
	if len(stringArgs) >= 2 {
		if len(stringArgs[1]) > 0 {
			printWithEmptyLine(stringArgs[1])
		}
	}
	printWithEmptyLine("\x1b[1m\x1b[32m" + "正在执行命令：" + stringArgs[0] + "\x1b[0m")
	err := cmd.Run()
	if err != nil {
		printErr(err)
	}
	if len(stringArgs) >= 3 {
		if len(stringArgs[2]) > 0 {
			printWithEmptyLine(stringArgs[2])
		}
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
	_, _ = fmt.Scanln()
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
