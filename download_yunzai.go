package main

import (
	ct "github.com/daviddengcn/go-colortext"
	"os"
)

func downloadYunzaiFromGitee() {
	_, err := os.Stat("./Yunzai-bot")
	if err == nil {
		printWithEmptyLine("检测到当前目录下已存在 Yunzai-bot ，请问是否需要重新下载？(是:y 返回菜单:n)")
		userChoice := ReadChoice("y", "n")
		if userChoice == "y" {
			ct.Foreground(ct.Red, true)
			printWithEmptyLine("重新下载云崽会移除当前目录下的 Yunzai-bot 文件夹，云崽的数据将会被全部删除，且不可恢复，请再次确认是否继续？(是:y 返回菜单:n)")
			ct.ResetColor()
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
	printWithEmptyLine("请选择云崽类别：\n1.官方云崽\n2.喵喵云崽")
	userChoice := ReadChoice("1", "2")
	if userChoice == "1" {
		executeCmd("git clone --depth 1 -b main https://gitee.com/yoimiya-kokomi/Yunzai-Bot.git", "开始下载官方云崽...", "下载官方云崽成功！")
	}
	var installMiao = false
	if userChoice == "2" {
		executeCmd("git clone --depth 1 -b master https://gitee.com/yoimiya-kokomi/Miao-Yunzai.git", "开始下载喵喵云崽...", "下载喵喵云崽成功！")
		//将Miao-Yunzai文件夹重命名为Yunzai-Bot
		os.Rename("./Miao-Yunzai", "./Yunzai-Bot")
		installMiao = true
	}
	//进入Yunzai-Bot文件夹
	os.Chdir("./Yunzai-Bot")
	b2 := checkCommand("pnpm -v")
	if !b2 {
		executeCmd("npm install pnpm -g --registry=https://registry.npmmirror.com", "开始安装 pnpm ...", "安装 pnpm 成功！")
	}
	executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...")
	executeCmd("pnpm config set PUPPETEER_DOWNLOAD_HOST=https://npmmirror.com/mirrors", "设置 pnpm 镜像源成功！")
	executeCmd("pnpm install", "开始安装云崽依赖", "安装云崽依赖成功！")
	if installMiao {
		installMiaoPlugin()
	}
	wd.changeToRoot()
}
