package pages

import (
	"fmt"
	"github.com/bling-yshs/HYZL/src/cmd/structs/global"
	"github.com/bling-yshs/HYZL/src/cmd/structs/menu_option"
	"github.com/bling-yshs/HYZL/src/cmd/utils/cmd_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/http_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/input_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/menu_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func BugsFixMenu() {
	_, err := os.Stat(global.Global.YunzaiName)
	if os.IsNotExist(err) {
		print_utils.PrintWithEmptyLine("未检测到云崽文件夹，请先下载云崽！")
		return
	}
	options := []menu_option.MenuOption{
		{"重装依赖", reInstallDep},
		{"修复 #重启 失败(也就是pnpm start pm2报错)", pm2Fix},
		{"修复 cookie 总是失效过期(Redis启动参数错误导致)", cookieRedisFix},
		{"修复 喵喵云崽监听报错(也就是sqlite3问题)", sqliteFix},
		{"修复 ffmpeg 未安装", ffmpegFix},
	}

	for {
		menu_utils.PrintMenu("BUG修复", options, false)
		choice := input_utils.ReadUint32()
		if choice == 0 {
			cmd_utils.ClearLog()
			return
		}
		menu_utils.DealChoice(choice, options, false)
	}

}

func ffmpegFix() {
	// https://hyzl.r2.yshs.fun/resources/ffmpeg.exe 下载
	ffmpegPath := path.Join(global.Global.ProgramRunPath, "ffmpeg", "ffmpeg.exe")
	err := http_utils.DownloadFile("https://hyzl.r2.yshs.fun/resources/ffmpeg.exe", ffmpegPath, true)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "原因：下载 ffmpeg.exe 失败！"))
		return
	}
	cmd := fmt.Sprintf(`setx PATH "%%PATH%%;%s"`, strings.ReplaceAll(filepath.Dir(ffmpegPath), "/", `\`))
	// setx PATH "%PATH%;ffmpeg的路径"
	cmd_utils.ExecuteCmd(cmd, "", "正在设置环境变量...", "设置环境变量成功！")
	print_utils.PrintWithColor(ct.Red, true, "请务必重新启动一次云崽！")
}

func sqliteFix() {
	_, err := os.Stat(path.Join(global.Global.YunzaiName, "plugins/miao-plugin/index.js"))
	if os.IsNotExist(err) {
		print_utils.PrintWithEmptyLine("检测到喵喵插件不完整，是否重新安装?(是:y 否:n)")
		choice := input_utils.ReadChoice([]string{"y", "n"})
		if choice == "y" {
			InstallMiaoPlugin()
		}
	}
	cmd_utils.ExecuteCmd("pnpm install sqlite3@5.1.5 -w", global.Global.YunzaiName, "正在安装 sqlite3...", "安装 sqlite3 成功！")
}
func cookieRedisFix() {
	_, err := os.Stat("redis-windows-7.0.4")
	if os.IsNotExist(err) {
		print_utils.PrintWithColor(ct.Red, true, "未检测到 redis-windows-7.0.4 文件夹，请先下载 redis！")
		return
	}
	err = http_utils.DownloadFile("https://gitee.com/bling_yshs/redis-windows-7.0.4/raw/master/redis.conf", "redis-windows-7.0.4", true)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "原因：下载 redis.conf 失败！"))
		return
	}
	print_utils.PrintWithEmptyLine("修复成功！")
}

func pm2Fix() {
	cmd_utils.ExecuteCmd("pnpm uninstall pm2", global.Global.YunzaiName, "正在卸载 pm2...", "卸载 pm2 成功！")
	cmd_utils.ExecuteCmd("pnpm install pm2@latest -w", global.Global.YunzaiName, "正在安装 pm2...", "安装 pm2 成功！")
}

func reInstallDep() {
	_, err := os.Stat(path.Join(global.Global.YunzaiName, "node_modules"))
	if os.IsNotExist(err) {
		print_utils.PrintWithEmptyLine("未检测到 node_modules 文件夹，是否重新安装依赖?(是:y 否:n)")
		choice := input_utils.ReadChoice([]string{"y", "n"})
		if choice == "n" {
			return
		}
	}
	cmd_utils.ExecuteCmd("npm config set registry https://registry.npmmirror.com", "", "开始设置 npm 镜像源...", "设置 npm 镜像源成功！")
	cmd_utils.ExecuteCmd("pnpm config set registry https://registry.npmmirror.com", "", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
	cmd_utils.ExecuteCmd("pnpm install", global.Global.YunzaiName, "开始安装云崽依赖...", "安装云崽依赖成功！")
}
