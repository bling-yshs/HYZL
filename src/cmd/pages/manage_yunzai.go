package pages

import (
	"bufio"
	"fmt"
	"github.com/bling-yshs/HYZL/src/cmd/structs/menu_option"
	"github.com/bling-yshs/HYZL/src/cmd/structs/yunzai"
	"github.com/bling-yshs/HYZL/src/cmd/utils/cmd_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/global_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/http_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/input_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/io_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/menu_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/redis_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/windows_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/yaml_utils"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/pkg/errors"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func ManageYunzaiMenu() {
	//检查是否存在Global.YunzaiName文件夹
	_, err := os.Stat(yunzai.GetYunzai().Name)
	if os.IsNotExist(err) {
		print_utils.PrintWithEmptyLine("未检测到云崽文件夹，请先下载云崽！")
		return
	}
	options := []menu_option.MenuOption{
		{"启动云崽", startYunzai},
		{"强制关闭云崽(强制关闭node程序和chrome)", closeYunzai},
		{"自定义终端命令", customCommand},
		{"安装插件", installPluginsMenu},
		{"安装js插件", installJsPlugin},
		{"修改云崽账号密码", changeAccount},
		{"强制更新云崽", updateYunzaiToLatest},
		{"设置签名API", setQsignAPI},
		{"手动更新所有插件", updateAllPlugins},
	}

	for {
		menu_utils.PrintMenu("云崽管理", options, false)
		choice := input_utils.ReadUint32()
		if choice == 0 {
			cmd_utils.ClearLog()
			return
		}
		menu_utils.DealChoice(choice, options, false)
	}
}

func updateAllPlugins() {
	// 得到yunzai/plugins文件夹下的所有文件夹，查看每个文件夹内是否存在package.json，如果存在，则在那个文件夹内执行git pull
	pluginsDir := filepath.Join(yunzai.GetYunzai().Path, "plugins")
	// 遍历插件目录
	dir, err := os.ReadDir(pluginsDir)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：读取插件目录失败"))
		return
	}
	for _, d := range dir {
		if d.IsDir() {
			pluginPath := filepath.Join(pluginsDir, d.Name())
			// 判断是否存在package.json
			_, err := os.Stat(filepath.Join(pluginPath, "package.json"))
			if os.IsNotExist(err) {
				continue
			}
			cmd_utils.ExecuteCmd("git pull", pluginPath, fmt.Sprintf("正在更新 %s", d.Name()), "")
		}
	}
}

func startYunzai() {
	if !redis_utils.IsRedisRunning() {
		redis_utils.StartRedis()
		// 等待1秒
		time.Sleep(1 * time.Second)
	}

	// 检查是否有node.exe在运行
	isNodeRunning := windows_utils.IsProcessRunning("node.exe")

	if isNodeRunning {
		print_utils.PrintWithEmptyLine("检测到后台存在 node 程序正在运行，可能为云崽的后台进程，是否关闭云崽并重新启动？(是:y 跳过:n)")
		choice := input_utils.ReadChoice([]string{"y", "n"})
		if choice == "y" {
			global_utils.ShutDownYunzai()
		}
	}
	print_utils.PrintWithEmptyLine("正在启动云崽...")
	cmd := exec.Command("cmd", "/C", "start", "cmd", "/k", "node app")
	cmd.Dir = path.Join(yunzai.GetYunzai().Path)
	cmd.Start()
	print_utils.PrintWithEmptyLine("云崽启动成功！")
}

func closeYunzai() {
	exec.Command("taskkill", "/FI", "WINDOWTITLE eq Yunzai-bot", "/T", "/F").Run()
	global_utils.ShutDownYunzai()
	global_utils.ShutDownChrome()
}

func customCommand() {
	for {
		fmt.Println()
		fmt.Print("请输入命令(输入0退出)：")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		command := scanner.Text()
		print_utils.PrintWithEmptyLine(command)
		if "0" == command {
			cmd_utils.ClearLog()
			break
		}
		cmd_utils.ExecuteCmd(command, yunzai.GetYunzai().Path, "", "")
	}
}

func installJsPlugin() {
	// 得到下载目录
	jsPluginDir := filepath.Join(yunzai.GetYunzai().Path, "plugins", "example")
	// 输入js插件的地址
	fmt.Print("请输入需要下载或复制的js插件的地址：")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	jsPluginUrl := scanner.Text()
	// 检查输入是否为https://开头，并且以js结尾
	if strings.HasPrefix(jsPluginUrl, "https://") && strings.HasSuffix(jsPluginUrl, ".js") {
		// 如果输入格式是https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/blob/master/%E5%96%9C%E6%8A%A5.js则自动转换为https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/raw/master/%E5%96%9C%E6%8A%A5.js
		jsPluginUrl = strings.Replace(jsPluginUrl, "blob", "raw", 1)
		fileName := filepath.Base(jsPluginUrl)
		unescapeFileName, err := url.QueryUnescape(fileName)
		if err != nil {
			print_utils.PrintError(errors.Wrap(err, "错误描述：解析文件名失败！"))
			return
		}
		err = http_utils.DownloadFile(jsPluginUrl, path.Join(jsPluginDir, unescapeFileName), false)
		if err != nil {
			print_utils.PrintError(errors.Wrap(err, "错误描述：下载失败！"))
			return
		}
		print_utils.PrintWithEmptyLine("下载成功！")
	} else if filepath.IsAbs(jsPluginUrl) && strings.HasSuffix(jsPluginUrl, ".js") {
		err := io_utils.CopyFile(jsPluginUrl, filepath.Join(jsPluginDir, filepath.Base(jsPluginUrl)))
		if err != nil {
			print_utils.PrintError(errors.Wrap(err, "错误描述：复制失败！"))
			return
		}
		fmt.Println("复制成功！")
	} else {
		fmt.Println("输入的js插件地址不正确！")
	}
}

func changeAccount() {
	fmt.Print("请输入 QQ 账号(输入0退出)：")
	qq := input_utils.ReadUint32()
	if qq == 0 {
		return
	}
	fmt.Print("请输入 QQ 密码(输入0退出)：")
	pwd := input_utils.ReadString()
	if pwd == "0" {
		return
	}
	err := yaml_utils.UpdateValueYAML(filepath.Join(yunzai.GetYunzai().Path, "config", "config", "qq.yaml"), "qq", qq)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：更新qq失败"))
		return
	}
	err = yaml_utils.UpdateValueYAML(filepath.Join(yunzai.GetYunzai().Path, "config", "config", "qq.yaml"), "pwd", pwd)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：更新密码失败"))
		return
	}
	print_utils.PrintWithEmptyLine("切换账号成功！")
}

func updateYunzaiToLatest() {
	cmd_utils.ExecuteCmd("git pull", yunzai.GetYunzai().Path, "正在更新云崽...", "")
	cmd_utils.ExecuteCmd("git reset --hard origin/HEAD", yunzai.GetYunzai().Path, "", "更新云崽成功")
}

func setQsignAPI() {
	// 先检查是否存在config/config/bot.yaml
	_, err := os.Stat(path.Join(yunzai.GetYunzai().Path, "config/config/bot.yaml"))
	if os.IsNotExist(err) {
		print_utils.PrintWithColor(ct.Red, true, "未检测到 bot.yaml 文件，如果您是第一次下载云崽，请先启动一次再运行此选项！")
		return
	}
	// 再检查是否存在config/config/qq.yaml
	_, err = os.Stat(path.Join(yunzai.GetYunzai().Path, "config/config/qq.yaml"))
	if os.IsNotExist(err) {
		print_utils.PrintWithColor(ct.Red, true, "未检测到 qq.yaml 文件，如果您是第一次下载云崽，请先启动一次再运行此选项！")
		return
	}
	print_utils.PrintWithColor(ct.Yellow, true, "请输入git bash地址(直接回车默认为C:/Program Files/Git/git-bash.exe):")
	reader := bufio.NewReader(os.Stdin)
	gitBashPath, err := reader.ReadString('\n')
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：读取输入失败"))
		return
	}
	gitBashPath = strings.TrimSpace(gitBashPath)
	if gitBashPath == "" {
		gitBashPath = "C:/Program Files/Git/git-bash.exe"
	}

	cmd := exec.Command(gitBashPath, "-c", "curl -L Gitee.com/haanxuan/QSign/raw/main/X | bash")
	cmd.Dir = yunzai.GetYunzai().Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "错误描述：执行命令失败"))
		return
	}
}
