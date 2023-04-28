package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func manageYunzaiMenu() {
	if !yunzaiExists() {
		printWithEmptyLine("当前目录下不存在云崽，请先下载云崽")
		return
	}
	for {
		fmt.Println("===云崽管理===")
		fmt.Println("1. 启动云崽")
		fmt.Println("2. 强制关闭云崽")
		fmt.Println("3. 切换账号")
		fmt.Println("4. 安装插件")
		fmt.Println("5. 安装js插件")
		fmt.Println("6. 自定义终端命令")
		fmt.Println("7. 更新云崽")
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
		case 7:
			clearLog()
			updateYunzaiToLatest()
		default:
			printWithEmptyLine("选择不正确，请重新选择")
		}
	}
}

func updateYunzaiToLatest() {
	os.Chdir("./Yunzai-Bot")
	err := executeCmd("git pull", "正在更新云崽...")
	if err != nil {
		executeCmd("git reset --hard origin/main")
	}
	os.Chdir("..")
}

// 检查云崽是否存在，存在返回true，不存在返回false
func yunzaiExists() bool {
	if _, err := os.Stat("./Yunzai-Bot"); err != nil {
		return false
	}
	if _, err := os.Stat("./Yunzai-Bot/package.json"); err != nil {
		return false
	}
	if _, err := os.Stat("./Yunzai-Bot/plugins"); err != nil {
		return false
	}
	return true
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

func changeMasterQQ() {
	// 读取 YAML 配置文件
	file, err := os.Open("./Yunzai-Bot/config/config/other.yaml")
	if err != nil {
		printErr(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var content string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		content += line
	}

	// 让用户输入新的 masterQQ 值
	fmt.Print("请输入新的主人QQ(直接回车将不改变主人QQ)：")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	newMasterQQ := scanner.Text()

	// 如果用户没有输入新值，就不修改文件
	if newMasterQQ == "" {
		return
	}
	// 修改 masterQQ 值
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "masterQQ:") {
			lines[i+1] = "  - " + newMasterQQ // 在下一行加入新的 masterQQ 值
			break
		}
	}

	newContent := strings.Join(lines, "\n")

	// 将修改后的内容写回文件
	err = os.WriteFile("./Yunzai-Bot/config/config/other.yaml", []byte(newContent), 0644)
	if err != nil {
		printErr(err)
	}

	printWithEmptyLine("主人QQ已修改为" + newMasterQQ)
}

func changeAccount() {

	fmt.Print("请输入 QQ 账号：")
	qq := readInt()
	fmt.Print("请输入密码：")
	pwd := readString()
	fmt.Print("请输入登录方式（1:安卓手机、2:aPad、3:安卓手表、4:MacOS、5:iPad）2023年4月24日：推荐使用5:iPad登录：")
	platform := readInt()
	changeMasterQQ()
	fileContent := fmt.Sprintf("# qq账号\nqq: %d\n# 密码，为空则用扫码登录,扫码登录现在仅能在同一ip下进行\npwd: '%s'\n# 1:安卓手机、 2:aPad 、 3:安卓手表、 4:MacOS 、 5:iPad\nplatform: %d", qq, pwd, platform)
	//覆盖掉./Yunzai-Bot/config/config/qq.yaml
	os.WriteFile("./Yunzai-Bot/config/config/qq.yaml", []byte(fileContent), fs.FileMode(0777))
	printWithEmptyLine("切换账号成功！")
}

func installJsPlugin() {
	//得到下载目录
	jsPluginDir := filepath.Join(programRunPath, "Yunzai-bot/plugins/example")
	//输入js插件的地址
	fmt.Print("请输入需要下载的js插件的地址：")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	jsPluginUrl := scanner.Text()
	//检查url是否为https://开头，并且以js结尾
	if !strings.HasPrefix(jsPluginUrl, "https://") || !strings.HasSuffix(jsPluginUrl, ".js") {
		executeCmd("输入的js插件地址不正确，请重新输入")
		return
	}
	//如果输入格式是https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/blob/master/%E5%96%9C%E6%8A%A5.js则自动转换为https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/raw/master/%E5%96%9C%E6%8A%A5.js
	jsPluginUrl = strings.Replace(jsPluginUrl, "blob", "raw", 1)
	downloadFile(jsPluginUrl, jsPluginDir)
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
