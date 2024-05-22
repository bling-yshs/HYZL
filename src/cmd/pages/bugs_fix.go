package pages

// import (
//
//	"bufio"
//	"fmt"
//	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
//	"github.com/bling-yshs/HYZL/src/cmd/utils/wd_utils"
//	ct "github.com/daviddengcn/go-colortext"
//	"os"
//	"path/filepath"
//	"strings"
//
// )
func BugsFixMenu() {

}

//	if !yunzaiExists() {
//		print_utils.PrintWithEmptyLine("当前目录下不存在云崽，请先下载云崽")
//		return
//	}
//	for {
//		options := []MenuOption{
//			{"重装依赖", reInstallDep},
//			{"修复 puppeteer 的各种问题", puppeteerProblemFix},
//			{"修复 云崽登录QQ失败(显示被风控发不出消息也可以尝试此选项)", icqqProblemFix},
//			{"修复 #重启 失败(也就是pnpm start pm2报错)", pm2Fix},
//			{"修复 cookie 总是失效过期(Redis启动参数错误导致)", cookieRedisFix},
//			{"修复 喵喵云崽监听报错(也就是sqlite3问题)", listenFix},
//		}
//
//		choice := showMenu("BUG修复", options, false)
//		if choice == 0 {
//			return
//		}
//	}
//}
//
//func listenFix() {
//
//	_, err := os.Stat("./plugins/miao-plugin")
//	if err != nil {
//		print_utils.PrintWithEmptyLine("检测到未安装喵喵插件，是否安装?(是:y 否:n)")
//		choice := input_utils.ReadChoice([]string{"y", "n"})
//		if choice == "y" {
//			installMiaoPlugin()
//		}
//		return
//	}
//	_, err = os.Stat("./plugins/miao-plugin/index.js")
//	if err != nil {
//		print_utils.PrintWithEmptyLine("检测到喵喵插件不完整，是否重新安装?(是:y 否:n)")
//		choice := input_utils.ReadChoice([]string{"y", "n"})
//		if choice == "y" {
//			installMiaoPlugin()
//		}
//		return
//	}
//	cmd_utils.ExecuteCmd("pnpm install sqlite3@5.1.5 -w", "正在修复sqlite3缺失...")
//}
//func cookieRedisFix() {
//	err := wd.changeToRedis()
//	if err != nil {
//		print_utils.PrintWithColor(ct.Red, true,
//			"Redis目录不存在！")
//		return
//	}
//	downloadFile("https://gitee.com/bling_yshs/redis-windows-7.0.4/raw/master/redis.conf", "./")
//	print_utils.PrintWithEmptyLine("修复成功！")
//}
//func pm2Fix() {
//
//	cmd_utils.ExecuteCmd("pnpm uninstall pm2", "正在修复...")
//	cmd_utils.ExecuteCmd("pnpm install pm2@latest -w", "", "修复成功！")
//}
//
//func icqqProblemFix() {
//
//	print_utils.PrintWithEmptyLine("开始修复云崽登录失败...")
//	_, err2 := os.Stat("./data")
//	//如果data文件夹存在
//	if err2 == nil {
//		dataDir, _ := filepath.Abs("./data")
//		files, err := os.ReadDir(dataDir)
//		if err == nil {
//			print_utils.PrintWithEmptyLine("正在删除 token 以及 device.json 缓存...")
//			for _, file := range files {
//				name := file.Name()
//				// 如果文件名以_token结尾,删除该文件
//				if strings.HasSuffix(name, "_token") {
//					os.Remove(filepath.Join(dataDir, name))
//				}
//				if name == "icqq" {
//					os.RemoveAll(filepath.Join(dataDir, name))
//				}
//				// 如果文件名为device.json,删除该文件
//				if name == "device.json" {
//					os.Remove(filepath.Join(dataDir, name))
//				}
//			}
//		}
//	}
//	cmd_utils.ExecuteCmd("pnpm uninstall icqq")
//	cmd_utils.ExecuteCmd("pnpm install icqq@0.6.1 -w")
//	//读取./config/config/qq.yaml
//	tools.UpdateValueYAML("./config/config/qq.yaml", "platform", 1)
//	print_utils.PrintWithEmptyLine("修复成功！")
//}
//
//func puppeteerProblemFix() {
//
//	print_utils.PrintWithEmptyLine("1.正常修复 2.通过edge修复(推荐)")
//	choice := ReadChoice("1", "2")
//	if choice == "1" {
//		cmd_utils.ExecuteCmd("pnpm install puppeteer -w")
//		print_utils.PrintWithColor(ct.Red, true,
//			"如果你非常长时间未下载成功，请尝试通过edge修复")
//		cmd_utils.ExecuteCmd("node ./node_modules/puppeteer/install.js")
//		print_utils.PrintWithEmptyLine("修复成功，大概")
//	}
//	if choice == "2" {
//		print_utils.PrintWithEmptyLine("请按以下步骤操作：\n1.下载并安装最新版的 edge 地址: https://yshs.lanzouj.com/iBe9O1clqsdc\n官网地址:https://www.microsoft.com/zh-cn/edge/download\n安装完成后回车进入下一步")
//		//等待用户回车
//		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
//		//检查是否存在C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe
//		_, err := os.Stat(`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`)
//		edgePath := ""
//		if err != nil {
//			print_utils.PrintWithColor(ct.Red, true,
//				"未识别到edge路径，请手动输入，例如C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe")
//			edgePath = readString()
//			edgePath = fmt.Sprintf(`%s/msedge.exe`, strings.ReplaceAll(edgePath, `\`, `/`))
//			_ = tools.UpdateOrAppendToYaml("./config/config/bot.yaml", "chromium_path", edgePath)
//		} else {
//			edgePath = `C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe`
//		}
//		_ = tools.UpdateOrAppendToYaml("./config/config/bot.yaml", "chromium_path", edgePath)
//		cmd_utils.ExecuteCmd("pnpm install puppeteer@21.1.1 -w --ignore-scripts")
//		print_utils.PrintWithEmptyLine("修复成功")
//	}
//
//}
//
//func reInstallDep() {
//
//	if _, err := os.Stat("./node_modules"); err == nil {
//		print_utils.PrintWithEmptyLine("检测到当前目录下已存在 node_modules ，请问是否需要重新安装依赖？(是:y 返回菜单:n)")
//		userChoice := input_utils.ReadChoice([]string{"y", "n"})
//		if userChoice == "y" {
//			cmd_utils.ExecuteCmd("npm config set registry https://registry.npmmirror.com")
//			cmd_utils.ExecuteCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...")
//			cmd_utils.ExecuteCmd("pnpm install", "开始安装云崽依赖...", "安装云崽依赖成功！")
//		}
//		if userChoice == "n" {
//			return
//		}
//	} else {
//		cmd_utils.ExecuteCmd("npm config set registry https://registry.npmmirror.com")
//		cmd_utils.ExecuteCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...")
//		cmd_utils.ExecuteCmd("pnpm install", "", "安装云崽依赖成功！")
//	}
//}
