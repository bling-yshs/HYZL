package pages

import (
	"bufio"
	"fmt"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/structs/global"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/structs/menu_option"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/cmd_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/global_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/http_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/input_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/io_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/menu_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/print_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/redis_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/windows_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/yaml_utils"
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
	_, err := os.Stat(global.Global.YunzaiName)
	if os.IsNotExist(err) {
		print_utils.PrintWithEmptyLine("未检测到云崽文件夹，请先下载云崽！")
		return
	}
	options := []menu_option.MenuOption{
		{"启动云崽", startYunzai},
		{"强制关闭云崽(强制关闭node程序)", closeYunzai},
		{"自定义终端命令", customCommand},
		{"安装插件", installPluginsMenu},
		{"安装js插件", installJsPlugin},
		{"修改云崽账号密码", changeAccount},
		{"强制更新云崽", updateYunzaiToLatest},
		{"设置qsign.icu的签名API", setQsignAPI},
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
	cmd.Dir = path.Join(global.Global.ProgramRunPath, global.Global.YunzaiName)
	cmd.Start()
	print_utils.PrintWithEmptyLine("云崽启动成功！")
}

func closeYunzai() {
	exec.Command("taskkill", "/FI", "WINDOWTITLE eq Yunzai-bot", "/T", "/F").Run()
	global_utils.ShutDownYunzai()
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
		cmd_utils.ExecuteCmd(command, global.Global.YunzaiName, "", "")
	}
}

func installJsPlugin() {
	// 得到下载目录
	jsPluginDir := filepath.Join(global.Global.ProgramRunPath, global.Global.YunzaiName, "plugins", "example")
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
			print_utils.PrintError(err)
			return
		}
		err = http_utils.DownloadFile(jsPluginUrl, path.Join(jsPluginDir, unescapeFileName))
		if err != nil {
			print_utils.PrintError(err)
			return
		}
		print_utils.PrintWithEmptyLine("下载成功！")
	} else if filepath.IsAbs(jsPluginUrl) && strings.HasSuffix(jsPluginUrl, ".js") {
		err := io_utils.CopyFile(jsPluginUrl, filepath.Join(jsPluginDir, filepath.Base(jsPluginUrl)))
		if err != nil {
			print_utils.PrintError(err)
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
	err := yaml_utils.UpdateValueYAML(filepath.Join(global.Global.ProgramRunPath, global.Global.YunzaiName, "config", "config", "qq.yaml"), "qq", qq)
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	err = yaml_utils.UpdateValueYAML(filepath.Join(global.Global.ProgramRunPath, global.Global.YunzaiName, "config", "config", "qq.yaml"), "pwd", pwd)
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	print_utils.PrintWithEmptyLine("切换账号成功！")
}

func updateYunzaiToLatest() {
	cmd_utils.ExecuteCmd("git pull", global.Global.YunzaiName, "正在更新云崽...", "")
	cmd_utils.ExecuteCmd("git reset --hard origin/HEAD", global.Global.YunzaiName, "", "更新云崽成功")
}

func setQsignAPI() {
	err := yaml_utils.UpdateOrAppendToYaml(path.Join(global.Global.ProgramRunPath, global.Global.YunzaiName, "config/config/bot.yaml"), "sign_api_addr", "https://hlhs-nb.cn/signed/sign?key=114514&ver=9.0.17")
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	err = yaml_utils.UpdateOrAppendToYaml(path.Join(global.Global.ProgramRunPath, global.Global.YunzaiName, "config/config/bot.yaml"), "ver", nil)
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	err = yaml_utils.UpdateOrAppendToYaml(path.Join(global.Global.ProgramRunPath, global.Global.YunzaiName, "config/config/qq.yaml"), "platform", "2")
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	print_utils.PrintWithEmptyLine("设置签名API成功！")
}
