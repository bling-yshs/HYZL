package main

import (
	"bufio"
	"fmt"
	"github.com/bling-yshs/YzLauncher-windows/tools"
	"github.com/mitchellh/go-ps"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type MenuOption struct {
	Label  string
	Action func()
}

func manageYunzaiMenu() {
	//检查是否存在yunzaiName文件夹
	if !yunzaiExists() {
		printWithEmptyLine("未检测到云崽文件夹，请先下载云崽！")
		return
	}
	for {
		options := []MenuOption{
			{"启动云崽", startYunzai},
			{"强制关闭云崽(强制关闭node程序)", closeYunzai},
			{"自定义终端命令", customCommand},
			{"安装插件", installPluginsMenu},
			{"安装js插件", installJsPlugin},
			{"修改云崽账号密码或者修改主人QQ", changeAccount},
			{"强制更新云崽", updateYunzaiToLatest},
			{"以node apps方式启动", startQQNTYunzai},
			{"设置qsign.icu的签名API", setQsignAPI},
		}

		choice := showMenu("云崽管理", options, false)
		if choice == 0 {
			return
		}
	}
}

func setQsignAPI() {
	wd.changeToYunzai()
	err := tools.UpdateOrAppendToYaml("./config/config/bot.yaml", "sign_api_addr", "https://hlhs-nb.cn/signed/sign?key=114514&ver=9.0.17")
	if err != nil {
		printWithRedColor("设置签名API失败！")
		return
	}
	err = tools.UpdateOrAppendToYaml("./config/config/bot.yaml", "ver", "")
	if err != nil {
		printWithRedColor("设置签名API失败！")
		return
	}
	tools.UpdateValueYAML("./config/config/qq.yaml", "platform", "2")
	printWithEmptyLine("设置签名API成功！")
}

func startQQNTYunzai() {
	if !isRedisRunning() {
		_ = startRedis()
		// 等待1秒
		time.Sleep(1 * time.Second)
	}
	wd.changeToYunzai()
	// 检查是否有node.exe在运行
	processList, err := ps.Processes()
	if err != nil {
		printWithRedColor("无权限获取进程列表!")
		return
	}

	isNodeRunning := false
	for _, process := range processList {
		if strings.ToLower(process.Executable()) == "node.exe" {
			isNodeRunning = true
			break
		}
	}

	if isNodeRunning {
		printWithEmptyLine("检测到后台存在 node 程序正在运行，可能为云崽的后台进程，是否关闭云崽并重新启动？(是:y 跳过:n)")
		choice := ReadChoice("y", "n")
		if choice == "y" {
			closeYunzai()
		}
	}
	printWithEmptyLine("正在启动云崽...")
	dir, _ := os.Getwd()
	cmd := exec.Command("cmd", "/C", "start", "/d", dir, "cmd", "/k", "node apps")
	cmd.Start()
	printWithEmptyLine("云崽启动成功！")
}

func updateYunzaiToLatest() {
	wd.changeToYunzai()
	executeCmd("git pull", "正在更新云崽...")
	executeCmd("git reset --hard origin/HEAD")
}

// 检查云崽是否存在，存在返回true，不存在返回false
func yunzaiExists() bool {
	wd.changeToRoot()
	if _, err := os.Stat(yunzaiName); err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(yunzaiName, "package.json")); err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(yunzaiName, "plugins")); err != nil {
		return false
	}
	return true
}

func startYunzai() {
	if !isRedisRunning() {
		_ = startRedis()
		// 等待1秒
		time.Sleep(1 * time.Second)
	}
	wd.changeToYunzai()
	// 检查是否有node.exe在运行
	processList, err := ps.Processes()
	if err != nil {
		printWithRedColor("无权限获取进程列表!")
		return
	}

	isNodeRunning := false
	for _, process := range processList {
		if strings.ToLower(process.Executable()) == "node.exe" {
			isNodeRunning = true
			break
		}
	}

	if isNodeRunning {
		printWithEmptyLine("检测到后台存在 node 程序正在运行，可能为云崽的后台进程，是否关闭云崽并重新启动？(是:y 跳过:n)")
		choice := ReadChoice("y", "n")
		if choice == "y" {
			closeYunzai()
		}
	}
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
	wd.changeToRoot()
	var isOtherYamlExists = true

	// 读取 YAML 配置文件
	stat, err := os.Stat(filepath.Join(yunzaiName, "config/config/other.yaml"))
	if err != nil || stat.Size() == 0 {
		fmt.Println(err)
		fmt.Println(stat.Size())
		isOtherYamlExists = false
		printWithEmptyLine("警告：检测到other.yaml配置文件内容为空，请问是否还原默认配置？(是:y 退出修改:n)")
		choice := ReadChoice("y", "n")
		if choice == "y" {
			stat, err := os.Stat(filepath.Join(yunzaiName, "config/default_config/other.yaml"))
			if err != nil || stat.Size() == 0 {
				os.RemoveAll(filepath.Join(yunzaiName, "config/default_config/other.yaml"))
				downloadFile("https://gitee.com/yoimiya-kokomi/Yunzai-Bot/raw/main/config/default_config/other.yaml", filepath.Join(yunzaiName, "config/config"))
			} else {
				copyFile(filepath.Join(yunzaiName, "config/default_config/other.yaml"), filepath.Join(yunzaiName, "config/config/other.yaml"))
			}
		}
		if choice == "n" {
			return
		}
	}

	content, err := os.ReadFile(filepath.Join(yunzaiName, "config/config/other.yaml"))

	var newMasterQQ int

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
			lines[i+1] = "  - " + strconv.Itoa(newMasterQQ) // 在下一行加入新的 masterQQ 值
			break
		}
	}

	newContent := strings.Join(lines, "\n")

	// 将修改后的内容写回文件
	err = os.WriteFile(filepath.Join(yunzaiName, "config/config/other.yaml"), []byte(newContent), os.ModePerm)
	if err != nil {
		printErr(err)
	}

	printWithEmptyLine("主人QQ已修改为" + strconv.Itoa(newMasterQQ))
}

func changeAccount() {
	wd.changeToRoot()
	fmt.Print("请输入 QQ 账号(直接回车将不改变QQ账号和密码)：")
	qq := readInt(true)
	if qq != 0 {
		fmt.Print("请输入密码：")
		pwd := readString()
		fmt.Print("请输入登录方式（1:安卓手机、2:aPad、3:安卓手表、4:MacOS、5:iPad、6:TIM）：")
		platform := readInt()
		fileContent := fmt.Sprintf("# qq账号\nqq: %d\n# 密码，为空则用扫码登录,扫码登录现在仅能在同一ip下进行\npwd: '%s'\n# 1:安卓手机、 2:aPad 、 3:安卓手表、 4:MacOS 、 5:iPad 、 6:old_Android\nplatform: %d", qq, pwd, platform)
		// 覆盖掉./Yunzai-Bot/config/config/qq.yaml
		os.WriteFile(filepath.Join(yunzaiName, "config/config/qq.yaml"), []byte(fileContent), os.ModePerm)
	}
	changeMasterQQ()
	printWithEmptyLine("切换账号成功！")
}

func installJsPlugin() {
	// 得到下载目录
	jsPluginDir := filepath.Join(programRunPath, filepath.Join(yunzaiName, "plugins/example"))
	// 输入js插件的地址
	fmt.Print("请输入需要下载或复制的js插件的地址：")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	jsPluginUrl := scanner.Text()
	// 检查输入是否为https://开头，并且以js结尾
	if strings.HasPrefix(jsPluginUrl, "https://") && strings.HasSuffix(jsPluginUrl, ".js") {
		// 如果输入格式是https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/blob/master/%E5%96%9C%E6%8A%A5.js则自动转换为https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/raw/master/%E5%96%9C%E6%8A%A5.js
		jsPluginUrl = strings.Replace(jsPluginUrl, "blob", "raw", 1)
		err := downloadFile(jsPluginUrl, jsPluginDir)
		if err == nil {
			fmt.Println("下载成功！")
		}
	} else if filepath.IsAbs(jsPluginUrl) && strings.HasSuffix(jsPluginUrl, ".js") {
		err := copyFile(jsPluginUrl, filepath.Join(jsPluginDir, filepath.Base(jsPluginUrl)))
		if err == nil {
			fmt.Println("复制成功！")
		}
	} else {
		fmt.Println("输入的js插件地址不正确！")
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
