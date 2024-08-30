package pages

import (
	"github.com/bling-yshs/HYZL/src/cmd/structs/app"
	"github.com/bling-yshs/HYZL/src/cmd/structs/config"
	"github.com/bling-yshs/HYZL/src/cmd/structs/menu_option"
	"github.com/bling-yshs/HYZL/src/cmd/structs/updater"
	"github.com/bling-yshs/HYZL/src/cmd/utils/cmd_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/input_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/menu_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
)

func IndexMenu() {
	print_utils.PrintWithEmptyLine("当前版本:", app.GetApp().Version)
	var updateText string = "立即更新启动器"
	if config.GetConfig().HaveUpdate {
		updateText = "发现新版本，立即更新启动器"
	}
	options := []menu_option.MenuOption{
		{"安装云崽", DownloadYunzaiMenu},
		{"云崽管理", ManageYunzaiMenu},
		{"BUG修复", BugsFixMenu},
		{updateText, updater.MenuUpdateRightNow},
	}
	for {
		menu_utils.PrintMenu("主菜单", options, true)
		choice := input_utils.ReadUint32()
		if choice == 0 {
			cmd_utils.ClearLog()
			return
		}
		menu_utils.DealChoice(choice, options, true)
	}
}
