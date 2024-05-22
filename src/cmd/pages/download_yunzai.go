package pages

import (
	"github.com/bling-yshs/HYZL/src/cmd/structs/global"
	"github.com/bling-yshs/HYZL/src/cmd/utils/cmd_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/git_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/global_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/input_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
	ct "github.com/daviddengcn/go-colortext"
	"os"
)

func DownloadYunzaiMenu() {
	_, err := os.Stat(global.Global.YunzaiName)
	if err == nil {
		print_utils.PrintWithEmptyLine("检测到当前目录下已存在 " + global.Global.YunzaiName + "，请问是否需要重新下载？(是:y 返回菜单:n)")
		userChoice := input_utils.ReadChoice([]string{"y", "n"})
		if userChoice == "y" {
			print_utils.PrintWithColor(ct.Red, true,
				"重新下载云崽会移除当前目录下的 "+global.Global.YunzaiName+" 文件夹，云崽的数据将会被全部删除，且不可恢复，请再次确认是否继续？(是:y 返回菜单:n)")
			userChoice := input_utils.ReadChoice([]string{"y", "n"})
			if userChoice == "n" {
				return
			}
			//删除文件夹
			print_utils.PrintWithEmptyLine("正在删除 " + global.Global.YunzaiName + " 文件夹...")
			os.RemoveAll(global.Global.YunzaiName)
		}
		if userChoice == "n" {
			return
		}
	}
	print_utils.PrintWithEmptyLine("目前已删除官方云崽的下载选项，现自动为您下载喵喵云崽")
	git_utils.Clone("https://gitee.com/yoimiya-kokomi/Miao-Yunzai.git", "", "开始下载喵喵云崽...", "下载喵喵云崽成功！")
	//将Miao-Yunzai文件夹重命名为Yunzai-Bot
	os.Rename("./Miao-Yunzai", "./Yunzai-Bot")
	//进入Yunzai-Bot文件夹

	b2 := cmd_utils.CheckCommandExist("pnpm -v")
	if !b2 {
		cmd_utils.ExecuteCmd("npm install pnpm -g --registry=https://registry.npmmirror.com", "", "开始安装 pnpm ...", "安装 pnpm 成功！")
	}
	cmd_utils.ExecuteCmd("npm config set registry https://registry.npmmirror.com", "", "开始设置 npm 镜像源...", "设置 npm 镜像源成功！")
	cmd_utils.ExecuteCmd("pnpm config set registry https://registry.npmmirror.com", "", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
	cmd_utils.ExecuteCmd("pnpm install", "./Yunzai-Bot", "开始安装云崽依赖", "安装云崽依赖成功！")
	//检查是否存在"node_modules/icqq/package.json"，如果不存在则报错提示
	_, err = os.Stat("./node_modules/icqq/package.json")
	if err != nil {
		print_utils.PrintWithColor(ct.Red, true,
			"检测到当前目录下不存在 node_modules/icqq/package.json ，初步判断为您的云崽依赖没有正常安装，请尝试使用 BUG修复->重装依赖，若还是无法解决，请到云崽仓库将安装依赖时的报错截图发 issue 反馈，地址 https://gitee.com/yoimiya-kokomi/Miao-Yunzai/issues")
		global_utils.ShutDownProgram()
	}
	print_utils.PrintWithEmptyLine("开始下载必须的喵喵插件...")
	InstallMiaoPlugin()
}
