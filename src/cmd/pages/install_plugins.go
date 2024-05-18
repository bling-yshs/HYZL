package pages

import (
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/structs/global"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/structs/menu_option"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/cmd_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/input_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/menu_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/print_utils"
	"os"
	"path/filepath"
)

func installPluginsMenu() {
	options := []menu_option.MenuOption{
		{"锅巴插件(云崽后台管理)(#锅巴帮助)", installGuobaPlugin},
		{"喵喵插件(#喵喵帮助)", InstallMiaoPlugin},
		{"逍遥插件(扫码登录)(#图鉴帮助)", installXiaoyaoPlugin},
		{"枫叶插件(小黑子插件)", installFengyePlugin},
		{"Atlas图鉴插件", installAtlasPlugin},
		{"土块插件(表情包制作)(#土块帮助)", installEarthKPlugin},
		{"闲心插件(#闲心帮助)", installXianxinPlugin},
		{"ap插件(#ap帮助)", installAPPlugin},
		{"flower插件(#百连)", installFlowerPlugin},
		{"梁氏插件(更好的伤害计算)", installLiangShiPlugin},
		{"星铁插件", installStarRailPlugin},
		{"打开云崽插件库", openYunzaiPluginLibrary},
	}

	for {
		menu_utils.PrintMenu("安装插件", options, false)
		choice := input_utils.ReadUint32()
		if choice == 0 {
			cmd_utils.ClearLog()
			return
		}
		menu_utils.DealChoice(choice, options, false)
	}

}

// https://gitee.com/hewang1an/StarRail-plugin.git
func installStarRailPlugin() {
	installPluginsTemplate("星铁插件 (https://gitee.com/hewang1an/StarRail-plugin.git)", "StarRail-plugin", "git clone --depth=1 https://gitee.com/hewang1an/StarRail-plugin.git ./plugins/StarRail-plugin/")
}

// https://gitee.com/liangshi233/liangshi-calc.git
func installLiangShiPlugin() {
	installPluginsTemplate("梁氏插件 (https://gitee.com/liangshi233/liangshi-calc.git)", "liangshi-calc", "git clone --depth=1 https://gitee.com/liangshi233/liangshi-calc.git ./plugins/liangshi-calc/")

}

// https://gitee.com/Nwflower/flower-plugin.git
func installFlowerPlugin() {
	installPluginsTemplate("flower插件 (https://gitee.com/Nwflower/flower-plugin.git)", "flower-plugin", "git clone --depth=1 https://gitee.com/Nwflower/flower-plugin.git ./plugins/flower-plugin/")

}

// https://gitee.com/Nwflower/atlas
func installAtlasPlugin() {
	installPluginsTemplate("Atlas图鉴插件 (https://gitee.com/Nwflower/atlas)", "Atlas", "git clone --depth=1 https://gitee.com/Nwflower/atlas ./plugins/Atlas/")

}

// https://gitee.com/yhArcadia/ap-plugin
func installAPPlugin() {
	installPluginsTemplate("AP插件 (https://gitee.com/yhArcadia/ap-plugin)", "ap-plugin", "git clone --depth=1 https://gitee.com/yhArcadia/ap-plugin ./plugins/ap-plugin/")
}

// 打开云崽插件库https://gitee.com/yhArcadia/Yunzai-Bot-plugins-index

func openYunzaiPluginLibrary() {
	cmd_utils.ExecuteCmd("start https://gitee.com/yhArcadia/Yunzai-Bot-plugins-index", "", "", "")
}

// ↓插件安装函数

// https://gitee.com/xianxincoder/xianxin-plugin
func installXianxinPlugin() {
	installPluginsTemplate("闲心插件 (https://gitee.com/xianxincoder/xianxin-plugin)", "xianxin-plugin", "git clone --depth=1 https://gitee.com/xianxincoder/xianxin-plugin ./plugins/xianxin-plugin/")

}

func installEarthKPlugin() {
	installPluginsTemplate("土块插件 (https://gitee.com/SmallK111407/earth-k-plugin)", "earth-k-plugin", "git clone --depth=1 https://gitee.com/SmallK111407/earth-k-plugin ./plugins/earth-k-plugin/")
}

func installGuobaPlugin() {
	installPluginsTemplate("锅巴插件 (https://gitee.com/guoba-yunzai/guoba-plugin)", "Guoba-Plugin", "git clone --depth=1 https://gitee.com/guoba-yunzai/guoba-plugin.git ./plugins/Guoba-Plugin/", "pnpm install --no-lockfile --filter=guoba-plugin -w")
}

func InstallMiaoPlugin() {
	installPluginsTemplate("喵喵插件 (https://gitee.com/yoimiya-kokomi/miao-plugin)", "miao-plugin", "git clone --depth 1 -b master https://gitee.com/yoimiya-kokomi/miao-plugin.git ./plugins/miao-plugin/", "pnpm add image-size -w")
}

func installXiaoyaoPlugin() {
	installPluginsTemplate("逍遥插件 (https://gitee.com/Ctrlcvs/xiaoyao-cvs-plugin.git)", "xiaoyao-cvs-plugin", "git clone --depth=1 https://gitee.com/Ctrlcvs/xiaoyao-cvs-plugin.git ./plugins/xiaoyao-cvs-plugin/", "pnpm add promise-retry -w", "pnpm add superagent -w")
}

func installFengyePlugin() {
	installPluginsTemplate("枫叶插件 (https://gitee.com/kesally/hs-qiqi-cv-plugin.git)", "hs-qiqi-plugin", "git clone --depth=1 https://gitee.com/kesally/hs-qiqi-cv-plugin.git ./plugins/hs-qiqi-plugin/")
}

func installPluginsTemplate(pluginChineseName string, dirName string, command ...string) {

	pluginDir := filepath.Join(global.Global.ProgramRunPath, global.Global.YunzaiName, "plugins", dirName)
	_, err := os.Stat(pluginDir)
	if err == nil {
		print_utils.PrintWithEmptyLine("当前已安装 " + pluginChineseName + "，请问是否需要重新安装？(是:y 返回菜单:n)")
		userChoice := input_utils.ReadChoice([]string{"y", "n"})
		if userChoice == "n" {
			return
		}
	}
	_ = os.RemoveAll(pluginDir)
	for _, cmd := range command {
		cmd_utils.ExecuteCmd(cmd, filepath.Join(global.Global.ProgramRunPath, global.Global.YunzaiName), "", "")
	}
}
