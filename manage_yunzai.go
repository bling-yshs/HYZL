package main

import (
	"bufio"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/bling-yshs/YzLauncher-windows/tools"
	"net"
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
			"从官方云崽切换为喵喵云崽",
			"启动签名API",
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
		case 8:
			clearLog()
			updateOfficialYunzaiToMiaoYunzai()
		case 9:
			clearLog()
			signApi()
		default:
			printWithEmptyLine("选择不正确，请重新选择")
		}
	}
}
func signApi() {
	wd.changeToRoot()
	//检测API文件夹是否存在
	_, err := os.Stat("API")
	if err != nil {
		printWithEmptyLine("当前目录下不存在API文件夹，请下载并解压API文件夹到当前目录")
		return
	}
	//检查platform是否为1或者2
	value, err := tools.GetValueFromYAMLFile(filepath.Join(yunzaiName, "config/config/qq.yaml"), "platform")
	if err != nil {
		printRedInfo("读取 config/config/qq.yaml 失败，请检查配置文件是否存在")
		return
	}
	if value != 1 && value != 2 {
		printRedInfo("当前配置文件中的 platform 值不为 1: Android 或者 2:AndroidPad ，可能会导致登录失败，是否需要修改？(y/n)")
		choice := ReadChoice("y", "n")
		if choice == "y" {
			printWithEmptyLine("请输入 1 或者 2")
			platform := ReadChoice("1", "2")
			tools.UpdateYAMLFile(filepath.Join(yunzaiName, "config/config/qq.yaml"), "platform", platform)
		}
	}
	//检查node_modules/icqq/package.json里的version是否大于0.4.8
	icqqVersionStr, err := tools.GetValueFromJSONFile(filepath.Join(yunzaiName, "node_modules/icqq/package.json"), "version")
	if err != nil {
		printRedInfo("读取 node_modules/icqq/package.json 值失败，请反馈给作者")
		return
	}
	icqqVersion, err := semver.NewVersion(icqqVersionStr.(string))
	minVersion, _ := semver.NewVersion("0.4.8")
	if icqqVersion.LessThan(minVersion) {
		printRedInfo("当前 icqq 版本过低，请更新 icqq 到 0.4.8 以上")
		return
	}
	//检测8080端口是否被占用
	port := "8080"
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		printRedInfo("当前" + port + "端口被占用，请检查签名API是否已经启动，或关闭占用端口的程序后重试")
		return
	} else {
		_ = listener.Close() // 关闭监听器以释放端口
	}
	//检查是否存在JAVA_HOME环境变量
	_, exists := os.LookupEnv("JAVA_HOME")
	if !exists {
		printWithEmptyLine("当前系统未设置JAVA_HOME环境变量，正在自动设置...")
		JavaHome := filepath.Join(programRunPath, "API", "jre-11.0.19")
		var setJavaHomeCommand string = "setx JAVA_HOME \"" + JavaHome + "\""
		executeCmd(setJavaHomeCommand, "正在设置JAVA_HOME环境变量...", "设置JAVA_HOME环境变量成功！")
		_ = os.Setenv("JAVA_HOME", JavaHome)
	}
	//修改bot.yaml，添加sign_api_addr: http://127.0.0.1:8080/sign
	_ = tools.AppendToYaml(filepath.Join(yunzaiName, "config/config/bot.yaml"), "sign_api_addr", "http://127.0.0.1:8080/sign")
	//运行./API/start.bat
	os.Chdir("./API")
	cmd := exec.Command("cmd", "/c", "start", "start.bat")
	cmd.Start()
}

func updateOfficialYunzaiToMiaoYunzai() {
	wd.changeToYunzai()
	printWithEmptyLine("请确认是否要切换为喵喵云崽，此操作不可逆！(y/n)")
	userChoice := ReadChoice("y", "n")
	if userChoice == "n" {
		return
	}
	if userChoice == "y" {
		executeCmd("git branch -m main master")
		executeCmd("git remote rm origin")
		executeCmd("git remote add origin https://gitee.com/yoimiya-kokomi/Miao-Yunzai.git")
		executeCmd("git fetch")
		executeCmd("git branch --set-upstream-to=origin/master master")
		executeCmd("git reset --hard origin/master")
		executeCmd("pnpm update")
		executeCmd("pnpm install")
		puppeteerProblemFix()
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
	stat, err := os.Stat(filepath.Join(yunzaiName, "config/config/other.yaml"))
	if err != nil || stat.Size() == 0 {
		isOtherYamlExists = false
		printWithEmptyLine("警告：检测到other.yaml配置文件内容为空，请问是否还原默认配置？(是:y 退出修改:n)")
		choice := ReadChoice("y", "n")
		if choice == "y" {
			stat, err := os.Stat(filepath.Join(yunzaiName, "config/default_config/other.yaml"))
			if err != nil || stat.Size() == 0 {
				downloadFile("https://gitee.com/yoimiya-kokomi/Yunzai-Bot/raw/main/config/default_config/other.yaml", filepath.Join(yunzaiName, "config/config/other.yaml"))
			} else {
				copyFile(filepath.Join(yunzaiName, "config/default_config/other.yaml"), filepath.Join(yunzaiName, "config/config/other.yaml"))
			}
		}
		if choice == "n" {
			return
		}
	}

	content, err := os.ReadFile(filepath.Join(yunzaiName, "config/config/other.yaml"))

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
	err = os.WriteFile(filepath.Join(yunzaiName, "config/config/other.yaml"), []byte(newContent), os.ModePerm)
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
		os.WriteFile(filepath.Join(yunzaiName, "config/config/qq.yaml"), []byte(fileContent), os.ModePerm)
	}
	changeMasterQQ()
	printWithEmptyLine("切换账号成功！")
}

func installJsPlugin() {
	//得到下载目录
	jsPluginDir := filepath.Join(programRunPath, filepath.Join(yunzaiName, "plugins/example"))
	//输入js插件的地址
	fmt.Print("请输入需要下载或复制的js插件的地址：")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	jsPluginUrl := scanner.Text()
	//检查输入是否为https://开头，并且以js结尾
	if strings.HasPrefix(jsPluginUrl, "https://") && strings.HasSuffix(jsPluginUrl, ".js") {
		//如果输入格式是https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/blob/master/%E5%96%9C%E6%8A%A5.js则自动转换为https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/raw/master/%E5%96%9C%E6%8A%A5.js
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
