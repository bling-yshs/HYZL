package main

import (
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
			{"修复 喵喵云崽监听报错(sqlite3缺失，注意可能对icqq版本也有要求)", sqliteFix},
		}

		choice := showMenu("BUG修复", options, false)
		if choice == 0 {
			return
		}
	}
}

func sqliteFix() {
	wd.changeToYunzai()
	executeCmd("pnpm install sqlite3@5.1.6 -w", "正在修复sqlite3缺失...")
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
	executeCmd("pnpm install icqq@0.4.8 -w")
	//读取./config/config/qq.yaml
	tools.UpdateYAMLFile("./config/config/qq.yaml", "platform", 1)
	printWithEmptyLine("修复成功！")
}

func puppeteerProblemFix() {
	wd.changeToYunzai()
	executeCmd("npm config set registry https://registry.npmmirror.com")
	executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
	executeCmd("npm config set PUPPETEER_DOWNLOAD_HOST=https://cdn.npmmirror.com/binaries", "开始设置 puppeteer Chromium 镜像源...", "设置 puppeteer Chromium 镜像源成功！")
	executeCmd("pnpm uninstall puppeteer -w", "正在修复 puppeteer Chromium...")
	if windowsVersion < 10 {
		//WinServer2012
		executeCmd("pnpm install puppeteer@19.7.3 -w")
		executeCmd("node ./node_modules/puppeteer/install.js")
		printWithEmptyLine("正在下载修复文件...")
		_, err := os.Stat("./renderers")
		if err != nil {
			downloadFile("https://gitee.com/bling_yshs/YzLauncher-windows/raw/master/NonProjectRequirements/WinServer2012YunzaiFix/Official/puppeteer.js", "./lib/puppeteer")
		} else {
			downloadFile("https://gitee.com/bling_yshs/YzLauncher-windows/raw/master/NonProjectRequirements/WinServer2012YunzaiFix/Miao/config_default.yaml", "./renderers/puppeteer")
			downloadFile("https://gitee.com/bling_yshs/YzLauncher-windows/raw/master/NonProjectRequirements/WinServer2012YunzaiFix/Miao/puppeteer.js", "./renderers/puppeteer/lib")
		}
		printWithEmptyLine("下载修复文件完成!")
	} else {
		//Win10+
		executeCmd("pnpm install puppeteer@19.8.3 -w")
		executeCmd("node ./node_modules/puppeteer/install.js")
		printWithEmptyLine("正在下载修复文件...")
		_, err := os.Stat("./renderers")
		if err != nil {
			downloadFile("https://gitee.com/bling_yshs/YzLauncher-windows/raw/master/NonProjectRequirements/Win10YunzaiFix/Official/puppeteer.js", "./lib/puppeteer")
		} else {
			downloadFile("https://gitee.com/bling_yshs/YzLauncher-windows/raw/master/NonProjectRequirements/Win10YunzaiFix/Miao/config_default.yaml", "./renderers/puppeteer")
			downloadFile("https://gitee.com/bling_yshs/YzLauncher-windows/raw/master/NonProjectRequirements/Win10YunzaiFix/Miao/puppeteer.js", "./renderers/puppeteer/lib")
		}
		printWithEmptyLine("下载修复文件完成!")
	}
}

func reInstallDep() {
	wd.changeToYunzai()
	if _, err := os.Stat("./node_modules"); err == nil {
		printWithEmptyLine("检测到当前目录下已存在 node_modules ，请问是否需要重新安装依赖？(是:y 返回菜单:n)")
		userChoice := ReadChoice("y", "n")
		if userChoice == "y" {
			if _, err = tools.CheckKeyInJSONFile("./package.json", "puppeteer"); err == nil {
				_ = tools.UpdateValueInJSONFile("./package.json", "dependencies", "puppeteer", "19.8.3")
			}
			executeCmd("npm config set registry https://registry.npmmirror.com")
			executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...")
			executeCmd("npm config set PUPPETEER_DOWNLOAD_HOST=https://cdn.npmmirror.com/binaries", "开始设置 puppeteer Chromium 镜像源...", "设置 puppeteer Chromium 镜像源成功！")
			os.RemoveAll("./node_modules")
			executeCmd("pnpm install", "开始安装云崽依赖...", "安装云崽依赖成功！")
		}
		if userChoice == "n" {
			return
		}
	} else {
		if _, err = tools.CheckKeyInJSONFile("./package.json", "puppeteer"); err == nil {
			_ = tools.UpdateValueInJSONFile("./package.json", "dependencies", "puppeteer", "19.8.3")
		}
		executeCmd("npm config set registry https://registry.npmmirror.com")
		executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...")
		executeCmd("npm config set PUPPETEER_DOWNLOAD_HOST=https://cdn.npmmirror.com/binaries", "开始设置 puppeteer Chromium 镜像源...", "设置 puppeteer Chromium 镜像源成功！")
		executeCmd("pnpm install", "", "安装云崽依赖成功！")
	}
}
