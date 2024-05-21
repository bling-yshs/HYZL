//go:generate goversioninfo
package main

import (
	"encoding/json"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/pages"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/schedule"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/structs/annoncement"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/structs/global"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/structs/updater"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/cmd_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/git_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/global_utils"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/print_utils"
	ct "github.com/daviddengcn/go-colortext"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 前置检查
	checkBeforeRun()
	// 显示公告
	annoncement.ShowAnnouncement()
	// 检查更新
	if updater.CheckUpdate() {
		// 询问是否更新
		if updater.AskUpdate() {
			// 更新启动器
			updater.UpdateRightNow()
		}
	}
	// 显示更新日志
	updater.ShowChangelog()
	// 初始化定时任务
	schedule.InitSchedule()
	// 显示主菜单
	pages.IndexMenu()
}

func checkBeforeRun() {
	// 检查程序权限
	checkProgramPermission()
	// 检查程序环境,git,node,npm
	checkProgramEnv()
	// 检查是否存在云崽文件
	if checkYunzaiFileExist() {
		print_utils.PrintWithColor(ct.Red, true, "检测到当前目录下可能存在云崽文件，请注意云崽启动器需要在云崽根目录的上一级目录下运行!")
	}
	// 读取配置文件
	readConfig()
	// 检查redis是否存在
	checkRedisExist()
	// 清理更新脚本
	updater.ClearUpdater()
}

func checkRedisExist() {
	_, err := os.Stat("redis-windows-7.0.4")
	if os.IsNotExist(err) {
		print_utils.PrintWithColor(ct.Yellow, true, "正在自动下载 redis ...")
		git_utils.Clone("https://gitee.com/bling_yshs/redis-windows-7.0.4", "", "", "")
	}
}

func readConfig() {
	// 检查是否存在./config/config.json
	_, err := os.Stat("./config/config.json")
	if os.IsNotExist(err) {
		bytes, _ := json.MarshalIndent(global.Config, "", "    ")
		// 创建./config/config.json
		_ = os.Mkdir("./config", os.ModePerm)
		_ = os.WriteFile("./config/config.json", bytes, os.ModePerm)
	}
	// 读取./config/config.json
	file, _ := os.ReadFile("./config/config.json")
	_ = json.Unmarshal(file, &global.Config)
	var needWrite = false
	if global.Config.GitInstalled == false {
		// 检查是否安装了git
		if !cmd_utils.CheckCommandExist("git -v") {
			print_utils.PrintWithEmptyLine("检测到未安装 Git ，请安装后继续")
			global_utils.ShutDownProgram()
		} else {
			global.Config.GitInstalled = true
			needWrite = true
		}
	}
	if global.Config.NodeInstalled == false {
		// 检查是否安装了node
		if !cmd_utils.CheckCommandExist("node -v") {
			print_utils.PrintWithEmptyLine("检测到未安装 Node.js ，请安装后继续")
			global_utils.ShutDownProgram()
		} else {
			global.Config.NodeInstalled = true
			needWrite = true
		}
	}
	if global.Config.NpmInstalled == false {
		// 检查是否安装了npm
		if !cmd_utils.CheckCommandExist("npm -v") {
			print_utils.PrintWithEmptyLine("检测到未安装 npm ，请安装后继续")
			global_utils.ShutDownProgram()
		} else {
			global.Config.NpmInstalled = true
			needWrite = true
		}
	}
	if needWrite {
		// 写入./config/config.json
		bytes, _ := json.MarshalIndent(global.Config, "", "    ")
		_ = os.WriteFile("./config/config.json", bytes, os.ModePerm)
	}
}

func checkProgramPermission() {

	if !cmd_utils.CheckCommandExist("dir") {
		print_utils.PrintWithEmptyLine("当前软件权限不足，请用管理员权限运行，若使用管理员权限依然无效，那么我也没有办法")
		global_utils.ShutDownProgram()
	}
	file, err := os.Create("test.txt")
	if err != nil {
		print_utils.PrintWithEmptyLine("当前软件权限不足，请用管理员权限运行，若使用管理员权限依然无效，那么我也没有办法")
		global_utils.ShutDownProgram()
	}
	file.Close()
	_ = os.Remove("test.txt")
}

func checkProgramEnv() {
	//获取当前程序路径，判断程序路径是否有空格，有则提示并shutdown
	path, err := os.Executable()
	if err != nil {
		return
	}
	path = filepath.Dir(path)
	if strings.Contains(path, " ") {
		print_utils.PrintWithColor(ct.Red, true, "当前程序路径下存在空格，请勿将本程序放在有空格的路径下")
		global_utils.ShutDownProgram()
	}
}

//func startRedis() error {
//	err := wd.changeToRedis()
//	if err != nil {
//		//err则说明没有redis文件夹
//		return err
//	}
//	print_utils.PrintWithEmptyLine("正在启动 Redis ...")
//	dir, _ := os.Getwd()
//	redisPath := filepath.Join(dir, "redis-server.exe")
//	redisConfigPath := filepath.Join(dir, "redis.conf")
//	cmd := exec.Command("cmd.exe", "/c", "start", redisPath, redisConfigPath)
//	err = cmd.Start()
//	if err != nil {
//		printErr(err)
//	}
//	println("Redis 启动成功！")
//	return nil
//}

//func isRedisRunning() bool {
//	processList, err := ps.Processes()
//	if err != nil {
//		print_utils.PrintWithColor(ct.Red, true, "无权限获取进程列表!")
//		return true
//	}
//
//	isRedisRunning := false
//	for _, process := range processList {
//		if strings.Contains(process.Executable(), "redis") {
//			isRedisRunning = true
//			break
//		}
//	}
//	if isRedisRunning {
//		return true
//	} else {
//		return false
//	}
//}

func checkYunzaiFileExist() bool {

	if _, err := os.Stat("./package.json"); err == nil {
		return true
	}
	if _, err := os.Stat("./plugins"); err == nil {
		return true
	}
	return false
}
