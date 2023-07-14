package main

import (
	"bufio"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/bling-yshs/YzLauncher-windows/tools"
	"github.com/mitchellh/go-ps"
	"net"
	"net/http"
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
	for {
		options := []MenuOption{
			{"启动签名API 并启动云崽", signApi},
			{"强制关闭云崽(强制关闭node程序)", closeYunzai},
			{"自定义终端命令", customCommand},
			{"安装插件", installPluginsMenu},
			{"安装js插件", installJsPlugin},
			{"修改云崽账号密码或者修改主人QQ", changeAccount},
			{"强制更新云崽", updateYunzaiToLatest},
			{"从官方云崽切换为喵喵云崽", updateOfficialYunzaiToMiaoYunzai},
			{"启动云崽", startYunzai},
		}

		choice := showMenu("云崽管理", options, false)
		if choice == 0 {
			return
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
	//检查文件目录是否正确
	//检查是否存在API/start.bat
	_, err = os.Stat("API/start.bat")
	if err != nil {
		//检查是否存在API/API/start.bat
		_, err = os.Stat("API/API/start.bat")
		if err != nil {
			printRedInfo("请确保 API 文件夹下的 start.bat 文件存在！")
			fmt.Println(err)
		} else {
			printRedInfo("检查到 API 文件夹下存在 API 文件夹嵌套，请将 API 文件夹下的 API 文件夹往上移动一层！")
		}
		return
	}
	//检查platform是否为1或者2
	value, err := tools.GetValueFromYAMLFile(filepath.Join(yunzaiName, "config/config/qq.yaml"), "platform")
	if err == nil {
		if value != 1 && value != 2 {
			printRedInfo("当前配置文件中的 platform 值不为 1: Android 或者 2:AndroidPad ，可能会导致登录失败，是否需要修改？(y/n)")
			choice := ReadChoice("y", "n")
			if choice == "y" {
				printWithEmptyLine("请输入 1 或者 2")
				platform := ReadChoice("1", "2")
				tools.UpdateYAMLFile(filepath.Join(yunzaiName, "config/config/qq.yaml"), "platform", platform)
			}
		}
	} else {
		printWithEmptyLine("检测到 config/config/qq.yaml 文件不存在，所以您可能是初次使用云崽，后续初始化时请注意选择登录方式为 1：Android，否则可能会导致登录失败")
	}
	//检查node_modules/icqq/package.json里的version是否大于0.4.10
	icqqVersionStr, err := tools.GetValueFromJSONFile(filepath.Join(yunzaiName, "node_modules/icqq/package.json"), "version")
	if err != nil {
		printRedInfo("读取 node_modules/icqq/package.json 值失败，请反馈给作者")
		return
	}
	icqqVersion, err := semver.NewVersion(icqqVersionStr.(string))
	minVersion, _ := semver.NewVersion("0.4.10")
	if !icqqVersion.Equal(minVersion) {
		printRedInfo("当前 icqq 版本不为 0.4.10，可能会导致签名 api 失效，是否需要自动将 icqq 更改到 0.4.10?(是:y 否:n)")
		printRedInfo("因为新版喵喵修改了参数，所以 icqq 版本也跟着调整了，如果你在使用 旧版喵喵+icqq-0.3.8，请选择否")
		choice := ReadChoice("y", "n")
		if choice == "y" {
			wd.changeToYunzai()
			executeCmd("pnpm uninstall icqq")
			executeCmd("pnpm install icqq@0.4.10 -w")
		}
	}
	wd.changeToRoot()
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
		printWithEmptyLine("当前系统未设置 JAVA_HOME 环境变量，正在自动设置...")
		JavaHome := filepath.Join(programRunPath, "API", "jre-11.0.19")
		var setJavaHomeCommand string = "setx JAVA_HOME \"" + JavaHome + "\""
		executeCmd(setJavaHomeCommand, "正在设置JAVA_HOME环境变量...", "设置JAVA_HOME环境变量成功！")
		_ = os.Setenv("JAVA_HOME", JavaHome)
	} else {
		env := os.Getenv("JAVA_HOME")
		_, err := os.Stat(env)
		//如果环境变量不存在，就设置JAVA_HOME环境变量
		if err != nil {
			printWithEmptyLine("当前系统 JAVA_HOME 环境变量所在文件夹不存在，正在自动设置新的环境变量...")
			JavaHome := filepath.Join(programRunPath, "API", "jre-11.0.19")
			var setJavaHomeCommand string = "setx JAVA_HOME \"" + JavaHome + "\""
			executeCmd(setJavaHomeCommand, "正在设置JAVA_HOME环境变量...", "设置JAVA_HOME环境变量成功！")
			_ = os.Setenv("JAVA_HOME", JavaHome)
		}
	}
	//修改bot.yaml，添加sign_api_addr: http://127.0.0.1:8080/sign
	_ = tools.AppendToYaml(filepath.Join(yunzaiName, "config/config/bot.yaml"), "sign_api_addr", "http://127.0.0.1:8080/sign")
	//运行./API/start.bat
	os.Chdir("./API")
	cmd := exec.Command("cmd", "/c", "start", "start.bat")
	cmd.Start()
	printWithEmptyLine("正在启动签名API，请等待弹出的新窗口内显示 [FEKit_]info: task_handle.h:74 TaskSystem not allow 后方可正常启动云崽")
	//每隔两秒向http://127.0.0.1:8080/sign发送一次get请求，直到返回200为止
	for {
		resp, err := http.Get("http://127.0.0.1:8080/sign")
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode == 200 {
			printWithEmptyLine("签名API启动成功！")
			break
		}
	}
	printWithEmptyLine("是否需要立即启动云崽？(y/n)")
	readChoice := ReadChoice("y", "n")
	if readChoice == "y" {
		startYunzai()
	}
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
		_ = startRedis()
		//等待1秒
		time.Sleep(1 * time.Second)
	}
	wd.changeToYunzai()
	//检查是否有node.exe在运行
	processList, err := ps.Processes()
	if err != nil {
		printRedInfo("无权限获取进程列表!")
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
	//
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
		fmt.Print("请输入登录方式（1:安卓手机、2:aPad、3:安卓手表、4:MacOS、5:iPad、6:TIM）：")
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
