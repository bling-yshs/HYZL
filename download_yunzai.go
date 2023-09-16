package main

import (
	"os"
)

func downloadYunzaiFromGitee() {
	_, err := os.Stat(yunzaiName)
	if err == nil {
		printWithEmptyLine("检测到当前目录下已存在 " + yunzaiName + "，请问是否需要重新下载？(是:y 返回菜单:n)")
		userChoice := ReadChoice("y", "n")
		if userChoice == "y" {
			printWithRedColor("重新下载云崽会移除当前目录下的 " + yunzaiName + " 文件夹，云崽的数据将会被全部删除，且不可恢复，请再次确认是否继续？(是:y 返回菜单:n)")
			userChoice := ReadChoice("y", "n")
			if userChoice == "n" {
				return
			}
			//删除文件夹
			printWithEmptyLine("正在删除 " + yunzaiName + " 文件夹...")
			os.RemoveAll(yunzaiName)
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
	wd.changeToYunzai()
	b2 := checkCommandExist("pnpm -v")
	if !b2 {
		executeCmd("npm install pnpm -g --registry=https://registry.npmmirror.com", "开始安装 pnpm ...", "安装 pnpm 成功！")
	}
	executeCmd("npm config set registry https://registry.npmmirror.com")
	executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...")
	executeCmd("pnpm install", "开始安装云崽依赖", "安装云崽依赖成功！")
	//检查是否存在"node_modules/icqq/package.json"，如果不存在则报错提示
	_, err = os.Stat("./node_modules/icqq/package.json")
	if err != nil {
		printWithRedColor("检测到当前目录下不存在 node_modules/icqq/package.json ，初步判断为您的云崽依赖没有正常安装，请尝试使用 BUG修复->重装依赖，若还是无法解决，请到云崽仓库将安装依赖时的报错截图发 issue 反馈，地址 https://gitee.com/yoimiya-kokomi/Miao-Yunzai/issues")
		shutdownApp()
	}
	if installMiao {
		printWithEmptyLine("开始下载必须的喵喵插件...")
		installMiaoPlugin()
	}
}
