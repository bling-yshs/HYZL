package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func installPluginsMenu() {
	_ = os.Chdir("./Yunzai-Bot")
	for {
		fmt.Println("===安装插件===")
		fmt.Println("1. 锅巴插件")
		fmt.Println("2. 喵喵插件")
		fmt.Println("3. 逍遥插件")
		fmt.Println("4. 枫叶插件")
		fmt.Println("5. 星穹铁道插件")
		fmt.Println("0. 返回上一级")
		fmt.Print("\n请选择操作：")
		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			printWithEmptyLine("输入错误，请重新选择")
			continue
		}

		switch choice {
		case 0:
			clearLog()
			_ = os.Chdir("..")
			return
		case 1:
			clearLog()
			installGuobaPlugin()
		case 2:
			clearLog()
			installMiaoPlugin()
		case 3:
			clearLog()
			installXiaoyaoPlugin()
		case 4:
			clearLog()
			installFengyePlugin()
		case 5:
			clearLog()
			installStarRailPlugin()
		default:
			printWithEmptyLine("选择不正确，请重新选择")
		}
	}
}

//↓插件安装函数

func installStarRailPlugin() {
	installPluginsTemplate("星穹铁道插件", "StarRail-plugin", "git clone --depth=1 https://gitee.com/hewang1an/StarRail-plugin.git ./plugins/StarRail-plugin/")
}

func installGuobaPlugin() {
	installPluginsTemplate("锅巴插件", "Guoba-Plugin", "git clone --depth=1 https://gitee.com/guoba-yunzai/guoba-plugin.git ./plugins/Guoba-Plugin/", "pnpm install --no-lockfile --filter=guoba-plugin -w")
}

func installMiaoPlugin() {
	installPluginsTemplate("喵喵插件", "miao-plugin", "git clone --depth 1 -b master https://gitee.com/yoimiya-kokomi/miao-plugin.git ./plugins/miao-plugin/", "pnpm add image-size -w")
}

func installXiaoyaoPlugin() {
	installPluginsTemplate("逍遥插件", "xiaoyao-cvs-plugin", "git clone --depth=1 https://gitee.com/Ctrlcvs/xiaoyao-cvs-plugin.git ./plugins/xiaoyao-cvs-plugin/", "pnpm add promise-retry -w", "pnpm add superagent -w")
}

func installFengyePlugin() {
	installPluginsTemplate("枫叶插件", "hs-qiqi-plugin", "git clone --depth=1  https://gitee.com/kesally/hs-qiqi-cv-plugin.git  ./plugins/hs-qiqi-plugin/")
}

func installPluginsTemplate(pluginChineseName string, dirName string, command ...string) {
	wd.changeToYunzai()
	pluginDir := filepath.Join(programRunPath, "Yunzai-Bot", "plugins", dirName)
	_, err := os.Stat(pluginDir)
	if err == nil {
		printWithEmptyLine("当前已安装 " + pluginChineseName + "，请问是否需要重新安装？(是:y 返回菜单:n)")
		userChoice := ReadChoice("y", "n")
		if userChoice == "n" {
			return
		}
	}
	_ = os.RemoveAll(pluginDir)
	for _, cmd := range command {
		executeCmd(cmd)
	}
}
