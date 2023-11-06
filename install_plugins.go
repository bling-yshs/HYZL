package main

import (
	"os"
	"path/filepath"
)

func installPluginsMenu() {
	for {
		options := []MenuOption{
			{"锅巴插件(云崽后台管理)(#锅巴帮助)", installGuobaPlugin},
			{"喵喵插件(#喵喵帮助)", installMiaoPlugin},
			{"逍遥插件(扫码登录)(#图鉴帮助)", installXiaoyaoPlugin},
			{"枫叶插件(小黑子插件)", installFengyePlugin},
			{"原神/星铁图鉴插件", installAtlas},
			{"土块插件(表情包制作)(#土块帮助)", installEarthKPlugin},
			{"闲心插件(#闲心帮助)", installXianxinPlugin},
			{"打开云崽插件库", openYunzaiPluginLibrary},
		}

		choice := showMenu("安装插件", options, false)
		if choice == 0 {
			return
		}
	}
}

// 打开云崽插件库https://gitee.com/yhArcadia/Yunzai-Bot-plugins-index
func openYunzaiPluginLibrary() {
	executeCmd("start https://gitee.com/yhArcadia/Yunzai-Bot-plugins-index")
}

// ↓插件安装函数

// https://gitee.com/xianxincoder/xianxin-plugin
func installXianxinPlugin() {
	installPluginsTemplate("闲心插件 (https://gitee.com/xianxincoder/xianxin-plugin)", "xianxin-plugin", "git clone --depth=1 https://gitee.com/xianxincoder/xianxin-plugin ./plugins/xianxin-plugin/")

}

func installEarthKPlugin() {
	installPluginsTemplate("土块插件 (https://gitee.com/diqiushengwu/earth-k-plugin)", "earth-k-plugin", "git clone --depth=1 https://gitee.com/diqiushengwu/earth-k-plugin ./plugins/earth-k-plugin/")
}

func installAtlas() {
	installPluginsTemplate("Atlas图鉴插件 (https://gitee.com/Nwflower/atlas)", "Atlas", "git clone --depth=1 https://gitee.com/Nwflower/atlas ./plugins/Atlas/")
}

func installGuobaPlugin() {
	installPluginsTemplate("锅巴插件 (https://gitee.com/guoba-yunzai/guoba-plugin)", "Guoba-Plugin", "git clone --depth=1 https://gitee.com/guoba-yunzai/guoba-plugin.git ./plugins/Guoba-Plugin/", "pnpm install --no-lockfile --filter=guoba-plugin -w")
}

func installMiaoPlugin() {
	installPluginsTemplate("喵喵插件 (https://gitee.com/yoimiya-kokomi/miao-plugin)", "miao-plugin", "git clone --depth 1 -b master https://gitee.com/yoimiya-kokomi/miao-plugin.git ./plugins/miao-plugin/", "pnpm add image-size -w")
}

func installXiaoyaoPlugin() {
	installPluginsTemplate("逍遥插件 (https://gitee.com/Ctrlcvs/xiaoyao-cvs-plugin.git)", "xiaoyao-cvs-plugin", "git clone --depth=1 https://gitee.com/Ctrlcvs/xiaoyao-cvs-plugin.git ./plugins/xiaoyao-cvs-plugin/", "pnpm add promise-retry -w", "pnpm add superagent -w")
}

func installFengyePlugin() {
	installPluginsTemplate("枫叶插件 (https://gitee.com/kesally/hs-qiqi-cv-plugin.git)", "hs-qiqi-plugin", "git clone --depth=1 https://gitee.com/kesally/hs-qiqi-cv-plugin.git ./plugins/hs-qiqi-plugin/")
}

func installPluginsTemplate(pluginChineseName string, dirName string, command ...string) {
	wd.changeToYunzai()
	pluginDir := filepath.Join(programRunPath, yunzaiName, "plugins", dirName)
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
