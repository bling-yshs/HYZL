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
	//检查是否存在yunzaiName文件夹
	if !yunzaiExists() {
		printWithEmptyLine("未检测到云崽文件夹，请先下载云崽！")
		return
	}
	for {
		options := []MenuOption{
			{"启动签名API 并启动云崽", startSignApiAndYunzai},
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
func startSignApiAndYunzai() {
	wd.changeToRoot()
	//检测API文件夹是否存在
	_, err := os.Stat("API")
	if err != nil {
		printWithEmptyLine("当前目录下不存在API文件夹，是否需要自动开始下载?(前提是你能正常下载 gitee 的文件)(是:y 返回:n)")
		c := ReadChoice("y", "n")
		if c == "y" {
			executeCmd("git clone --depth 1 https://gitee.com/bling_yshs/hyzl-sign-api.git")
			_, err := os.Stat("hyzl-sign-api")
			if err != nil {
				printWithEmptyLine("当前无法连接到 Gitee，请从网盘内下载签名API")
			}
			_, err = os.Stat("hyzl-sign-api/start.bat")
			if err != nil {
				printWithEmptyLine("当前无法连接到 Gitee，请从网盘内下载签名API")
			}
			os.Rename("hyzl-sign-api", "API")
		} else {
			return
		}
	}
	//检查是否存在API/start.bat
	_, err = os.Stat("API/start.bat")
	if err != nil {
		//检查是否存在API/API/start.bat
		_, err = os.Stat("API/API/start.bat")
		if err != nil {
			printWithRedColor("请确保 API 文件夹下的 start.bat 文件存在！")
			fmt.Println(err)
		} else {
			printWithRedColor("检查到 API 文件夹下存在 API 文件夹嵌套，请将 API 文件夹下的 API 文件夹往上移动一层！")
		}
		return
	}
	//检查是否存在API/当前 API 版本 1.2.0文件
	_, err = os.Stat("API/当前 API 版本 1.2.0")
	if err != nil {
		printWithEmptyLine("检查到新版API，是否需要自动更新？(如果旧版API可以正常使用那么请勿更新)(是:y 否:n)")
		c := ReadChoice("y", "n")
		if c == "y" {
			os.Chdir("API")
			//检查是否存在.git文件夹
			_, err = os.Stat(".git")
			if err != nil {
				printWithRedColor("检测到当前 API 文件夹不是通过 Gitee 自动下载的，请问是否需要完全删除并重新下载 API ？(是:y 返回:n)")
				c := ReadChoice("y", "n")
				if c == "y" {
					wd.changeToRoot()
					os.RemoveAll("API")
					executeCmd("git clone --depth 1 https://gitee.com/bling_yshs/hyzl-sign-api.git")
					os.Rename("hyzl-sign-api", "API")
				} else {
					return
				}
			} else {
				executeCmd("git pull")
				executeCmd("git reset --hard origin/HEAD")
			}
		}
	}
	wd.changeToRoot()
	//检查platform是否为1或者2
	value, err := tools.GetValueFromYAMLFile(filepath.Join(yunzaiName, "config/config/qq.yaml"), "platform")
	if err == nil {
		if value != 1 {
			printWithRedColor("当前配置文件中的 platform 值不为 1: Android(安卓手机) ，可能会导致登录失败，是否需要修改？(y/n)")
			choice := ReadChoice("y", "n")
			if choice == "y" {
				tools.UpdateValueYAML(filepath.Join(yunzaiName, "config/config/qq.yaml"), "platform", 1)
			}
		}
	} else {
		printWithRedColor("检测到 config/config/qq.yaml 文件不存在，所以您可能是初次使用云崽\n后续初始化时请注意一定要选择登录方式为 Android(安卓手机)。\n签名API地址请填入 http://127.0.0.1:1539/sign?key=191539 \n或者您重新登陆一次，启动器会自动帮您修改")
	}
	//检查node_modules/icqq/package.json里的version是否大于0.6.1
	icqqVersionStr, err := tools.GetValueFromJSONFile(filepath.Join(yunzaiName, "node_modules/icqq/package.json"), "version")
	if err != nil {
		printWithRedColor("读取 node_modules/icqq/package.json 值失败，初步判断为您的云崽依赖没有正常安装，请尝试使用 BUG修复->重装依赖，若还是无法解决，请到云崽仓库将安装依赖时的报错截图发 issue 反馈，地址 https://gitee.com/yoimiya-kokomi/Miao-Yunzai/issues")
		return
	}
	icqqVersion, err := semver.NewVersion(icqqVersionStr.(string))
	minVersion, _ := semver.NewVersion("0.6.1")
	if icqqVersion.LessThan(minVersion) {
		printWithRedColor("当前 icqq 版本小于 0.6.1，可能会导致签名 api 失效，是否需要自动将 icqq 更改到 0.6.1?(是:y 否:n)")
		choice := ReadChoice("y", "n")
		if choice == "y" {
			wd.changeToYunzai()
			executeCmd("pnpm uninstall icqq")
			executeCmd("pnpm install icqq@0.6.1 -w")
		}
	}

	skipEditConfig := false

	wd.changeToRoot()
	if config.IsAPIPortSet == false {
		fmt.Print("请输入签名API端口，直接回车则使用默认端口1539：")
		i := readInt(true)
		if i == 0 {
			apiPort = 1539
			config.IsAPIPortSet = true
			writeConfig(&config)
			//写入到filepath.Join(programRunPath, "API", "txlib", "8.9.80", "config.json")
			err := tools.UpdateValueInJSONFile(filepath.Join(programRunPath, "API", "txlib", "8.9.80", "config.json"), "server", "port", apiPort)
			if err != nil {
				return
			}
		} else {
			apiPort = i
			config.IsAPIPortSet = true
			writeConfig(&config)
			err := tools.UpdateValueInJSONFile(filepath.Join(programRunPath, "API", "txlib", "8.9.80", "config.json"), "server", "port", apiPort)
			if err != nil {
				return
			}
		}
	} else {
		//从filepath.Join(programRunPath, "API", "txlib", "8.9.80", "config.json")读取端口号
		port, err := tools.GetValueFromJSONFile(filepath.Join(programRunPath, "API", "txlib", "8.9.80", "config.json"), "port")
		if err != nil {
			skipEditConfig = true
		} else {
			apiPort = int(port.(float64))
		}

		//从filepath.Join(programRunPath, "API", "txlib", "8.9.80", "config.json")读取key
		key, err := tools.GetValueFromJSONFile(filepath.Join(programRunPath, "API", "txlib", "8.9.80", "config.json"), "key")
		if err != nil {
			skipEditConfig = true
		} else {
			apiKey = key.(string)
		}
	}

	//检测1539端口是否被占用，并提示是否关闭

	port := apiPort

	// 获取当前系统上的所有网络连接
	connections, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		printWithRedColor("当前" + strconv.Itoa(port) + "端口被占用，是否需要关闭占用端口的进程?(是:y 否:n)")
		choice := ReadChoice("y", "n")
		if choice == "y" {
			pid, err := portInUse(port)
			if err != nil {
				printWithRedColor("错误：监听端口失败，却无法正常关闭占用端口")
				return
			}
			killCmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
			if err = killCmd.Run(); err != nil {
				fmt.Println("无法终止进程:", err)
				return
			}
		}
	} else {
		_ = connections.Close()
	}
	//检查是否存在JAVA_HOME环境变量
	_, exists := os.LookupEnv("JAVA_HOME")
	if !exists {
		printWithEmptyLine("当前系统未设置 JAVA_HOME 环境变量，正在自动设置...")
		JavaHome := filepath.Join(programRunPath, "API", "jre-11.0.19")
		setJavaHomeCommand := "setx JAVA_HOME \"" + JavaHome + "\""
		executeCmd(setJavaHomeCommand, "正在设置 JAVA_HOME 环境变量...", "设置 JAVA_HOME 环境变量成功！")
		_ = os.Setenv("JAVA_HOME", JavaHome)
	} else {
		env := os.Getenv("JAVA_HOME")
		_, err := os.Stat(env)
		//如果环境变量不存在，就设置JAVA_HOME环境变量
		if err != nil {
			printWithEmptyLine("当前系统 JAVA_HOME 环境变量所在文件夹不存在，正在自动设置新的环境变量...")
			JavaHome := filepath.Join(programRunPath, "API", "jre-11.0.19")
			setJavaHomeCommand := "setx JAVA_HOME \"" + JavaHome + "\""
			executeCmd(setJavaHomeCommand, "正在设置 JAVA_HOME 环境变量...", "设置 JAVA_HOME 环境变量成功！")
			_ = os.Setenv("JAVA_HOME", JavaHome)
		}
	}

	if !skipEditConfig {
		//修改bot.yaml，添加sign_api_addr: http://127.0.0.1:1539/sign?key=191539
		_ = tools.UpdateOrAppendToYaml(filepath.Join(yunzaiName, "config/config/bot.yaml"), "sign_api_addr", "http://127.0.0.1:"+strconv.Itoa(port)+"/sign?key="+apiKey)
	}
	//运行./API/start.bat
	os.Chdir("./API")
	cmd := exec.Command("cmd", "/c", "start", "start.bat")
	cmd.Start()
	wd.changeToRoot()
	printWithEmptyLine("正在启动签名API，请勿关闭此窗口...")
	//每隔五秒向http://127.0.0.1:1539/sign?key=191539发送一次get请求，直到返回200为止
	for {
		resp, err := http.Get("http://127.0.0.1:" + strconv.Itoa(port) + "/sign?key=" + apiKey)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		if resp.StatusCode == 200 {
			printWithEmptyLine("签名API启动成功！")
			break
		}
	}
	startYunzai()
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
	}
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
