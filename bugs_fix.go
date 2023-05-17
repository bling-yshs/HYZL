package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func bugsFixMenu() {
	if !yunzaiExists() {
		printWithEmptyLine("当前目录下不存在云崽，请先下载云崽")
		return
	}
	for {
		wd.changeToRoot()
		options := []string{
			"重装依赖",
			"修复 puppeteer 启动失败(Windows Server 2012专用)",
			"修复 puppeteer 弹出cmd窗口 或者 puppeteer启动失败(Windows Server 2012请勿使用)",
			"修复 云崽登录QQ失败",
			"修复 #重启 失败(也就是pnpm start pm2报错)",
			"修复 cookie 总是失效过期(Redis启动参数错误导致)",
		}

		choice := showMenu("BUG修复", options, false)

		switch choice {
		case 0:
			clearLog()
			return
		case 1:
			clearLog()
			reInstallDep()
		case 2:
			clearLog()
			pupCanNotStartFix()
		case 3:
			clearLog()
			pupPopFix()
		case 4:
			clearLog()
			icqqProblemFix()
		case 5:
			clearLog()
			pm2Fix()
		case 6:
			clearLog()
			cookieRedisFix()
		case 7:
			clearLog()
			return
		default:
			printWithEmptyLine("选择不正确，请重新选择")
		}
	}
}

func cookieRedisFix() {
	wd.changeToRedis()
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
	printWithEmptyLine("开始修复 错误码45 错误码238 QQ版本过低...")
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
	executeCmd("pnpm install icqq@latest -w")
	//读取./config/config/qq.yaml
	s, err := getFileContent("./config/config/qq.yaml")
	if err != nil {
		printErr(err)
		return
	}
	regex := regexp.MustCompile(`platform: \d`)
	s = regex.ReplaceAllString(s, `platform: 6`)
	//写入./config/config/qq.yaml
	err = os.WriteFile("./config/config/qq.yaml", []byte(s), os.ModePerm)
	if err != nil {
		printErr(err)
		return
	}
	printWithEmptyLine("修复成功！")
}

func pupPopFix() {
	wd.changeToYunzai()
	executeCmd("git reset --hard origin/HEAD")
	executeCmd("git pull", "正在更新云崽到最新版本...", "更新云崽到最新版本成功！")
	executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
	executeCmd("pnpm config set PUPPETEER_DOWNLOAD_HOST=https://npmmirror.com/mirrors", "开始设置 puppeteer Chromium 镜像源...", "设置 puppeteer Chromium 镜像源成功！")
	executeCmd("pnpm uninstall puppeteer", "正在修复 puppeteer Chromium...")
	executeCmd("pnpm install puppeteer@19.8.3 -w")
	executeCmd("node ./node_modules/puppeteer/install.js")
}

func reInstallDep() {
	wd.changeToYunzai()
	if _, err := os.Stat("./node_modules"); err == nil {
		printWithEmptyLine("检测到当前目录下已存在 node_modules ，请问是否需要重新安装依赖？(是:y 返回菜单:n)")
		userChoice := ReadChoice("y", "n")
		if userChoice == "y" {
			executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...")
			executeCmd("pnpm config set PUPPETEER_DOWNLOAD_HOST=https://npmmirror.com/mirrors", "开始设置 puppeteer Chromium 镜像源...", "设置 puppeteer Chromium 镜像源成功！")
			os.RemoveAll("./node_modules")
			executeCmd("pnpm install", "开始安装云崽依赖...", "安装云崽依赖成功！")
		}
		if userChoice == "n" {
			return
		}
	} else {
		executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...")
		executeCmd("pnpm config set PUPPETEER_DOWNLOAD_HOST=https://npmmirror.com/mirrors", "开始设置 puppeteer Chromium 镜像源...", "设置 puppeteer Chromium 镜像源成功！")
		executeCmd("pnpm install", "", "安装云崽依赖成功！")
	}
}

func pupCanNotStartFix() {
	wd.changeToYunzai()
	executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
	executeCmd("pnpm config set PUPPETEER_DOWNLOAD_HOST=https://npmmirror.com/mirrors", "开始设置 puppeteer Chromium 镜像源...", "设置 puppeteer Chromium 镜像源成功！")
	executeCmd("pnpm uninstall puppeteer", "正在修复 puppeteer Chromium...")
	executeCmd("pnpm install puppeteer@19.7.3 -w")
	executeCmd("node ./node_modules/puppeteer/install.js")
	printWithEmptyLine("正在下载cmd弹窗修复文件...")
	downloadFile("https://gitee.com/bling_yshs/YzLauncher-windows/raw/master/NonProjectRequirements/puppeteer.js", "./lib/puppeteer")
}
