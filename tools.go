package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/James-Ye/go-frame/win"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	ct "github.com/daviddengcn/go-colortext"
)

type WorkingDirectory struct{}

func (dir *WorkingDirectory) changeToRoot() {
	os.Chdir(programRunPath)
}

func (dir *WorkingDirectory) changeToYunzai() {
	os.Chdir(filepath.Join(programRunPath, yunzaiName))
}

func (dir *WorkingDirectory) changeToRedis() {
	os.Chdir(filepath.Join(programRunPath, "redis-windows-7.0.4"))
}

// ↓工具函数

// 参数canBeEmpty如果为true，表示输入可以为空，返回0
func readInt(canBeEmpty ...bool) int64 {

	allowEmpty := false
	if len(canBeEmpty) > 0 {
		allowEmpty = canBeEmpty[0]
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		s := scanner.Text()
		if s == "" {
			if allowEmpty {
				return 0
			}
			fmt.Print("输入不能为空，请重新输入：")
			continue
		}
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			fmt.Print("输入错误，请重新输入：")
			continue
		}
		return i
	}
}

// 参数canBeEmpty如果为true，表示输入可以为空，返回空字符串
func readString(canBeEmpty ...bool) string {
	allowEmpty := false
	if len(canBeEmpty) > 0 {
		allowEmpty = canBeEmpty[0]
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		s := scanner.Text()
		if s == "" {
			if allowEmpty {
				return ""
			}
			fmt.Print("输入不能为空，请重新输入：")
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

func copyFile(srcPath string, dstPath string) error {
	// 打开要复制的源文件
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 复制源文件到目标文件
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func downloadFile(downloadURL string, downloadFilePath string) error {
	res, err := http.Get(downloadURL)
	if res.StatusCode != http.StatusOK {
		return errors.New("网页返回错误状态码")
	}
	if err != nil {
		printWithEmptyLine("下载文件失败，请检查网络连接，错误信息为：" + err.Error())
		return err
	}
	fileName, _ := url.QueryUnescape(filepath.Base(downloadURL))

	var savePath = downloadFilePath
	if downloadFilePath == "" {
		savePath = config.SystemTempPath
	}
	filePath := filepath.Join(savePath, fileName)
	_, err = os.Stat(filePath)
	if err != nil {
		os.MkdirAll(savePath, os.ModePerm)
	}
	f, err := os.Create(filePath)
	if err != nil {
		printWithEmptyLine("创建文件失败，错误信息为：" + err.Error())
		return err
	}
	defer f.Close()

	_, _ = io.Copy(f, res.Body)
	return nil
}

func showMenu(title string, options []string, isMainMenu bool) int {
	for {
		fmt.Println("===" + title + "===")
		for i, option := range options {
			fmt.Printf("%d. %s\n", i+1, option)
		}
		if isMainMenu {
			fmt.Println("0. 退出程序")
		} else {
			fmt.Println("0. 返回上一级")
		}

		fmt.Print("\n请选择操作：")

		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("输入错误，请重新选择")
			continue
		}

		if choice > len(options) {
			fmt.Println("选择不正确，请重新选择")
			continue
		}

		return choice
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

	//获取程序名字
	*args[1] = filepath.Base(execPath)
	*args[2] = filepath.Join(currentDir, "config")
	_, err = os.Stat("./Miao-Yunzai")
	if err == nil {
		*args[3] = "Miao-Yunzai"
	}

}
func getAppInfoInt(args ...*int64) {
	majorVersion, _, _ := win.RtlGetNtVersionNumbers()
	*args[0] = int64(majorVersion)
}

func printRedInfo(str any) {
	ct.Foreground(ct.Red, true)
	printWithEmptyLine(str)
	ct.ResetColor()
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

func executeCmd(stringArgs ...string) error {
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
	ct.Foreground(ct.Green, true)
	printWithEmptyLine("正在执行命令：" + stringArgs[0])
	ct.ResetColor() // 重置颜色
	err := cmd.Run()
	if err != nil {
		return err
	}
	if len(stringArgs) >= 3 {
		if len(stringArgs[2]) > 0 {
			printWithEmptyLine(stringArgs[2])
		}
	}
	return nil
}

func printErr(err error) {
	fmt.Println("发生错误，请截图并反馈给作者:", err)
}

func printWithEmptyLine(str any) {
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
