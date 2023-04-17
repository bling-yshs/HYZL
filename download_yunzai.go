package main

import (
	"os"
)

func downloadYunzaiFromGitee() {
	_, err := os.Stat("./Yunzai-bot")
	if err == nil {
		printWithEmptyLine("检测到当前目录下已存在 Yunzai-bot ，请问是否需要重新下载？(是:y 返回菜单:n)")
		userChoice := ReadChoice("y", "n")
		if userChoice == "y" {
			printWithEmptyLine("\x1b[1m\x1b[31m重新下载云崽会移除当前目录下的 Yunzai-bot 文件夹，云崽的数据将会被全部删除，且不可恢复，请再次确认是否继续？(是:y 返回菜单:n)\x1b[0m")
			userChoice := ReadChoice("y", "n")
			if userChoice == "n" {
				return
			}
			//删除文件夹
			printWithEmptyLine("正在删除 Yunzai-bot 文件夹...")
			os.RemoveAll("./Yunzai-bot")
		}
		if userChoice == "n" {
			return
		}
	}
	executeCmd("git clone --depth 1 -b main https://gitee.com/yoimiya-kokomi/Yunzai-Bot.git", "开始下载云崽...", "下载云崽成功！")
	//进入Yunzai-Bot文件夹
	os.Chdir("./Yunzai-Bot")
	b2 := checkCommand("pnpm -v")
	if !b2 {
		executeCmd("npm install pnpm -g --registry=https://registry.npmmirror.com", "开始安装 pnpm ...", "安装 pnpm 成功！")
	}
	executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...")
	executeCmd("pnpm config set PUPPETEER_DOWNLOAD_HOST=https://npmmirror.com/mirrors", "设置 pnpm 镜像源成功！")
	executeCmd("pnpm install -P", "开始安装云崽依赖", "安装云崽依赖成功！")
	os.Chdir("..")
}
