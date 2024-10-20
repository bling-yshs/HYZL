package pages

import (
	"fmt"
	"github.com/bling-yshs/HYZL/src/cmd/structs/app"
	"github.com/bling-yshs/HYZL/src/cmd/structs/menu_option"
	"github.com/bling-yshs/HYZL/src/cmd/structs/yunzai"
	"github.com/bling-yshs/HYZL/src/cmd/utils/cmd_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/http_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/input_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/io_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/menu_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func BugsFixMenu() {
	_, err := os.Stat(yunzai.GetYunzai().Name)
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
		{"修复 无法正确获取图片链接", fixGetImageLink},
		{"修复 桌面出现白色窗口(治标不治本)", fixPuppeteerWhiteWindow},
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

func fixPuppeteerWhiteWindow() {
	// 原理，复制Yunzai-Bot/renderers/puppeteer/config_default.yaml为config.yaml，然后用yaml解析，修改里面的args数组，添加--window-position=-10000,-10000
	defaultConfigPath := path.Join(yunzai.GetYunzai().Path, "renderers", "puppeteer", "config_default.yaml")
	configPath := path.Join(yunzai.GetYunzai().Path, "renderers", "puppeteer", "config.yaml")
	err := io_utils.CopyFile(defaultConfigPath, configPath)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：复制 config_default.yaml 失败！"))
		return
	}
	// 开始修改 config.yaml
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：读取 config.yaml 失败！"))
		return
	}
	var node yaml.Node
	if err := yaml.Unmarshal(bytes, &node); err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：解析 config.yaml 失败！"))
		return
	}
	list := node.Content[0]
	for index, each := range list.Content {
		if each.Value == "args" {
			// 找到 "args" 节点
			argsNode := list.Content[index+1]
			// 创建新的节点
			newNode := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Style: 0,
				Tag:   "!!str",
				Value: "--window-position=-10000,-10000",
			}
			// 追加到 "args" 节点的内容中
			argsNode.Content = append(argsNode.Content, newNode)
			break
		}
	}
	// 将修改后的 Node 写回 config.yaml 文件
	newData, err := yaml.Marshal(&node)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：序列化 config.yaml 失败！"))
		return
	}
	err = os.WriteFile(configPath, newData, 0644)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：写回 config.yaml 失败！"))
		return
	}
	print_utils.PrintWithEmptyLine("修复成功！")
	// 找到 args 数组并追加 --window-position=-10000,-10000

	//var config map[string]interface{}
	//err = yaml.Unmarshal(bytes, &config)
	//if err != nil {
	//	print_utils.PrintError(errors.Wrap(err, "错误描述：解析 config.yaml 失败！"))
	//	return
	//}
	//// 找到 args 数组并追加 --window-position=-10000,-10000
	//if renderers, ok := config["args"].([]interface{}); ok {
	//	renderers = append(renderers, "--window-position=-10000,-10000")
	//	config["args"] = renderers
	//} else {
	//	print_utils.PrintError(errors.New("未找到 args 数组或类型错误"))
	//	return
	//}
	//
	//// 将修改后的配置重新序列化为 YAML 格式
	//newData, err := yaml.Marshal(&config)
	//if err != nil {
	//	print_utils.PrintError(errors.Wrap(err, "错误描述：序列化 config.yaml 失败！"))
	//	return
	//}
	//
	//// 写回 config.yaml 文件
	//err = ioutil.WriteFile(configPath, newData, 0644)
	//if err != nil {
	//	print_utils.PrintError(errors.Wrap(err, "错误描述：写回 config.yaml 失败！"))
	//	return
	//}

}

func fixGetImageLink() {
	fmt.Println("正在下载修复文件，请确保能正常访问 gitee")
	// https://gitee.com/bling_yshs/resources/raw/master/HYZL/parser.js 下载到 Yunzai-Bot/node_modules/icqq/lib/message/parser.js
	parserPath := path.Join(yunzai.GetYunzai().Path, "node_modules", "icqq", "lib", "message", "parser.js")
	fmt.Println(parserPath)
	url := "https://gitee.com/bling_yshs/resources/raw/master/HYZL/parser.js"
	err := http_utils.DownloadFile(url, parserPath, true)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：下载 parser.js 失败！"))
	}
}

func ffmpegFix() {
	// https://hyzl.r2.yshs.fun/resources/ffmpeg.exe 下载
	ffmpegPath := path.Join(app.GetApp().Workdir, "ffmpeg", "ffmpeg.exe")
	err := http_utils.DownloadFile("https://hyzl.r2.yshs.fun/resources/ffmpeg.exe", ffmpegPath, true)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：下载 ffmpeg.exe 失败！"))
		return
	}
	cmd := fmt.Sprintf(`setx PATH "%%PATH%%;%s"`, strings.ReplaceAll(filepath.Dir(ffmpegPath), "/", `\`))
	// setx PATH "%PATH%;ffmpeg的路径"
	cmd_utils.ExecuteCmd(cmd, "", "正在设置环境变量...", "设置环境变量成功！")
	print_utils.PrintWithColor(ct.Red, true, "请务必重新启动一次云崽！")
}

func sqliteFix() {
	_, err := os.Stat(path.Join(yunzai.GetYunzai().Path, "plugins/miao-plugin/index.js"))
	if os.IsNotExist(err) {
		print_utils.PrintWithEmptyLine("检测到喵喵插件不完整，是否重新安装?(是:y 否:n)")
		choice := input_utils.ReadChoice([]string{"y", "n"})
		if choice == "y" {
			InstallMiaoPlugin()
		}
	}
	cmd_utils.ExecuteCmd("pnpm install sqlite3@5.1.5 -w", yunzai.GetYunzai().Path, "正在安装 sqlite3...", "安装 sqlite3 成功！")
}
func cookieRedisFix() {
	_, err := os.Stat("redis-windows-7.0.4")
	if os.IsNotExist(err) {
		print_utils.PrintWithColor(ct.Red, true, "未检测到 redis-windows-7.0.4 文件夹，请先下载 redis！")
		return
	}
	err = http_utils.DownloadFile("https://gitee.com/bling_yshs/redis-windows-7.0.4/raw/master/redis.conf", "redis-windows-7.0.4/redis.conf", true)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：下载 redis.conf 失败！"))
		return
	}
	print_utils.PrintWithEmptyLine("修复成功！")
}

func pm2Fix() {
	cmd_utils.ExecuteCmd("pnpm uninstall pm2", yunzai.GetYunzai().Path, "正在卸载 pm2...", "卸载 pm2 成功！")
	cmd_utils.ExecuteCmd("pnpm install pm2@latest -w", yunzai.GetYunzai().Path, "正在安装 pm2...", "安装 pm2 成功！")
}

func reInstallDep() {
	_, err := os.Stat(path.Join(yunzai.GetYunzai().Path, "node_modules"))
	if os.IsNotExist(err) {
		print_utils.PrintWithEmptyLine("未检测到 node_modules 文件夹，是否重新安装依赖?(是:y 否:n)")
		choice := input_utils.ReadChoice([]string{"y", "n"})
		if choice == "n" {
			return
		}
	}
	cmd_utils.ExecuteCmd("npm config set registry https://registry.npmmirror.com", "", "开始设置 npm 镜像源...", "设置 npm 镜像源成功！")
	cmd_utils.ExecuteCmd("pnpm config set registry https://registry.npmmirror.com", "", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
	cmd_utils.ExecuteCmd("pnpm install", yunzai.GetYunzai().Path, "开始安装云崽依赖...", "安装云崽依赖成功！")
}
