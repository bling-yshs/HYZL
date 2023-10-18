package main

import (
	"bufio"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/bling-yshs/YzLauncher-windows/tools"
	"os"
	"path/filepath"
	"strings"
)

func bugsFixMenu() {
	if !yunzaiExists() {
		printWithEmptyLine("当前目录下不存在云崽，请先下载云崽")
		return
	}
	for {
		options := []MenuOption{
			{"重装依赖", reInstallDep},
			{"修复 puppeteer 的各种问题", puppeteerProblemFix},
			{"修复 云崽登录QQ失败(显示被风控发不出消息也可以尝试此选项)", icqqProblemFix},
			{"修复 #重启 失败(也就是pnpm start pm2报错)", pm2Fix},
			{"修复 cookie 总是失效过期(Redis启动参数错误导致)", cookieRedisFix},
			{"修复 喵喵云崽监听报错(也就是sqlite3问题)", listenFix},
		}

		choice := showMenu("BUG修复", options, false)
		if choice == 0 {
			return
		}
	}
}

func listenFix() {
	wd.changeToYunzai()
	_, err := os.Stat("./plugins/miao-plugin")
	if err != nil {
		printWithEmptyLine("检测到未安装喵喵插件，是否安装?(是:y 否:n)")
		choice := ReadChoice("y", "n")
		if choice == "y" {
			installMiaoPlugin()
		}
		return
	}
	_, err = os.Stat("./plugins/miao-plugin/index.js")
	if err != nil {
		printWithEmptyLine("检测到喵喵插件不完整，是否重新安装?(是:y 否:n)")
		choice := ReadChoice("y", "n")
		if choice == "y" {
			installMiaoPlugin()
		}
		return
	}
	executeCmd("pnpm install sqlite3@5.1.5 -w", "正在修复sqlite3缺失...")
}
func cookieRedisFix() {
	err := wd.changeToRedis()
	if err != nil {
		printWithRedColor("Redis目录不存在！")
		return
	}
	downloadFile("https://gitee.com/bling_yshs/redis-windows-7.0.4/raw/master/redis.conf", "./")
	printWithEmptyLine("修复成功！")
}
func pm2Fix() {
	wd.changeToYunzai()
	executeCmd("pnpm uninstall pm2", "正在修复...")
	executeCmd("pnpm install pm2@latest -w", "", "修复成功！")
}

func icqqProblemFix() {
	wd.changeToYunzai()
	printWithEmptyLine("开始修复云崽登录失败...")
	_, err2 := os.Stat("./data")
	//如果data文件夹存在
	if err2 == nil {
		dataDir, _ := filepath.Abs("./data")
		files, err := os.ReadDir(dataDir)
		if err == nil {
			printWithEmptyLine("正在删除 token 以及 device.json 缓存...")
			for _, file := range files {
				name := file.Name()
				// 如果文件名以_token结尾,删除该文件
				if strings.HasSuffix(name, "_token") {
					os.Remove(filepath.Join(dataDir, name))
				}
				if name == "icqq" {
					os.RemoveAll(filepath.Join(dataDir, name))
				}
				// 如果文件名为device.json,删除该文件
				if name == "device.json" {
					os.Remove(filepath.Join(dataDir, name))
				}
			}
		}
	}
	executeCmd("pnpm uninstall icqq")
	executeCmd("pnpm install icqq@0.5.4 -w")
	//读取./config/config/qq.yaml
	tools.UpdateValueYAML("./config/config/qq.yaml", "platform", 1)
	printWithEmptyLine("修复成功！")
}

func puppeteerProblemFix() {
	wd.changeToYunzai()
	printWithEmptyLine("1.正常修复(推荐) 2.通过edge修复")
	choice := ReadChoice("1", "2")
	if choice == "1" {
		executeCmd("pnpm install puppeteer -w")
		printWithRedColor("如果你非常长时间未下载成功，请尝试通过edge修复")
		executeCmd("node ./node_modules/puppeteer/install.js")
		printWithEmptyLine("修复成功，大概")
	}
	if choice == "2" {
		printWithEmptyLine("请按以下步骤操作：\n1.下载并安装最新版的edge 地址: https://www.microsoft.com/zh-cn/edge/download 安装完成后回车进入下一步")
		//等待用户回车
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		//检查是否存在C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe
		_, err := os.Stat(`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`)
		edgePath := ""
		if err != nil {
			printWithRedColor("未识别到edge路径，请手动输入，例如C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe")
			edgePath = readString()
			edgePath = fmt.Sprintf(`%s/msedge.exe`, strings.ReplaceAll(edgePath, `\`, `/`))
			_ = tools.UpdateOrAppendToYaml("./config/config/bot.yaml", "chromium_path", edgePath)
		} else {
			edgePath = `C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe`
		}
		_ = tools.UpdateOrAppendToYaml("./config/config/bot.yaml", "chromium_path", edgePath)
		//判断yunzaiName\node_modules\puppeteer\package.json内version的值是否小于21.1.1，是的话就执行pnpm install puppeteer -w
		version, err := tools.GetValueFromJSONFile("./node_modules/puppeteer/package.json", "version")
		if err != nil {
			printWithRedColor("获取 puppeteer 版本号失败，是否需要重新安装？(是:y 返回菜单:n)")
			choice := ReadChoice("y", "n")
			if choice == "y" {
				executeCmd("pnpm install puppeteer -w")
			}
			return
		}
		pupVersion, _ := semver.NewVersion(version.(string))
		minVersion, _ := semver.NewVersion("21.1.1")
		if pupVersion.LessThan(minVersion) {
			executeCmd("pnpm install puppeteer -w")
		}
		printWithEmptyLine("修复成功")
	}

}

func reInstallDep() {
	wd.changeToYunzai()
	if _, err := os.Stat("./node_modules"); err == nil {
		printWithEmptyLine("检测到当前目录下已存在 node_modules ，请问是否需要重新安装依赖？(是:y 返回菜单:n)")
		userChoice := ReadChoice("y", "n")
		if userChoice == "y" {
			executeCmd("npm config set registry https://registry.npmmirror.com")
			executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...")
			executeCmd("pnpm install", "开始安装云崽依赖...", "安装云崽依赖成功！")
		}
		if userChoice == "n" {
			return
		}
	} else {
		executeCmd("npm config set registry https://registry.npmmirror.com")
		executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...")
		executeCmd("pnpm install", "", "安装云崽依赖成功！")
	}
}
