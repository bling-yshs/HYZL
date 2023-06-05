package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func manageYunzaiMenu() {
	if !yunzaiExists() {
		printWithEmptyLine("当前目录下不存在云崽，请先下载云崽")
		return
	}
	for {
		wd.changeToRoot()
		options := []string{
			"启动云崽",
			"强制关闭云崽",
			"修改云崽账号密码或者修改主人QQ",
			"安装插件",
			"安装js插件",
			"自定义终端命令",
			"强制更新云崽",
		}

		choice := showMenu("云崽管理", options, false)

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
	wd.changeToYunzai()
	err := executeCmd("git pull", "正在更新云崽...")
	if err != nil {
		executeCmd("git reset --hard origin/HEAD")
	}
}

// 检查云崽是否存在，存在返回true，不存在返回false
func yunzaiExists() bool {
	wd.changeToRoot()
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
	wd.changeToYunzai()
	printWithEmptyLine("正在启动云崽...")
	dir, _ := os.Getwd()
	cmd := exec.Command("cmd", "/C", "start", "/d", dir, "cmd", "/k", "node app")
	cmd.Start()
	printWithEmptyLine("云崽启动成功！")
}

func closeYunzai() {
	exec.Command("taskkill", "/FI", "WINDOWTITLE eq Yunzai-bot", "/T", "/F").Run()
	executeCmd("taskkill /f /im node.exe", "正在关闭云崽...", "云崽关闭成功！")
}

func changeMasterQQ() {
	var isOtherYamlExists = true
	// 读取 YAML 配置文件
	stat, err := os.Stat("./Yunzai-Bot/config/config/other.yaml")
	if err != nil || stat.Size() == 0 {
		isOtherYamlExists = false
		printWithEmptyLine("警告：检测到other.yaml配置文件内容为空，请问是否还原默认配置？(是:y 退出修改:n)")
		choice := ReadChoice("y", "n")
		if choice == "y" {
			stat, err := os.Stat("./Yunzai-Bot/config/default_config/other.yaml")
			if err != nil || stat.Size() == 0 {
				downloadFile("https://gitee.com/yoimiya-kokomi/Yunzai-Bot/raw/main/config/default_config/other.yaml", "./Yunzai-Bot/config/config/other.yaml")
			} else {
				copyFile("./Yunzai-Bot/config/default_config/other.yaml", "./Yunzai-Bot/config/config/other.yaml")
			}
		}
		if choice == "n" {
			return
		}
	}

	content, err := os.ReadFile("./Yunzai-Bot/config/config/other.yaml")

	var newMasterQQ int64

	if isOtherYamlExists == true {
		fmt.Print("请输入新的主人QQ(直接回车将不改变主人QQ)：")
		newMasterQQ = readInt(true)

		// 如果用户没有输入新值，就不修改文件
		if newMasterQQ == 0 {
			return
		}
	}
	if isOtherYamlExists == false {
		for {
			fmt.Print("请输入新的主人QQ：")
			newMasterQQ = readInt()
			if newMasterQQ == 0 {
				printWithEmptyLine("主人QQ不能为空！")
				continue
			} else {
				break
			}
		}

	}

	// 修改 masterQQ 值
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "masterQQ:") {
			lines[i+1] = "  - " + strconv.FormatInt(newMasterQQ, 10) // 在下一行加入新的 masterQQ 值
			break
		}
	}

	newContent := strings.Join(lines, "\n")

	// 将修改后的内容写回文件
	err = os.WriteFile("./Yunzai-Bot/config/config/other.yaml", []byte(newContent), os.ModePerm)
	if err != nil {
		printErr(err)
	}

	printWithEmptyLine("主人QQ已修改为" + strconv.FormatInt(newMasterQQ, 10))
}

func changeAccount() {
	fmt.Print("请输入 QQ 账号(直接回车将不改变QQ账号和密码)：")
	qq := readInt(true)
	if qq != 0 {
		fmt.Print("请输入密码：")
		pwd := readString()
		fmt.Print("请输入登录方式（1:安卓手机、2:aPad、3:安卓手表、4:MacOS、5:iPad、6:old_Android）2023年4月24日：推荐使用6:old_Android登录：")
		platform := readInt()
		fileContent := fmt.Sprintf("# qq账号\nqq: %d\n# 密码，为空则用扫码登录,扫码登录现在仅能在同一ip下进行\npwd: '%s'\n# 1:安卓手机、 2:aPad 、 3:安卓手表、 4:MacOS 、 5:iPad 、 6:old_Android\nplatform: %d", qq, pwd, platform)
		//覆盖掉./Yunzai-Bot/config/config/qq.yaml
		os.WriteFile("./Yunzai-Bot/config/config/qq.yaml", []byte(fileContent), os.ModePerm)
	}
	changeMasterQQ()
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
		printWithEmptyLine("输入的js插件地址不正确，请重新输入")
		return
	}
	//如果输入格式是https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/blob/master/%E5%96%9C%E6%8A%A5.js则自动转换为https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/raw/master/%E5%96%9C%E6%8A%A5.js
	jsPluginUrl = strings.Replace(jsPluginUrl, "blob", "raw", 1)
	err := downloadFile(jsPluginUrl, jsPluginDir)
	if err == nil {
		printWithEmptyLine("下载成功！")
	}
}

func customCommand() {
	wd.changeToYunzai()
	for {
		fmt.Println()
		fmt.Print("请输入命令(输入0退出)：")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		command := scanner.Text()
		printWithEmptyLine(command)
		if "0" == command {
			break
		}
		executeCmd(command)
	}
	// 读取用户输入的一串字符串
}
