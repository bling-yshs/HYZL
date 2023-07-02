// 编译： go build
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func createNormalConfig(config Config) {
	//检查当前目录下是否存在config文件夹
	_, err := os.Stat("./config")

	//如果存在就不创建
	if err == nil {
		return
	}

	//创建文件
	err = os.Mkdir("./config", os.ModePerm)
	if err != nil {
		printErr(err)
		return
	}

	//再创建config.json
	file, err := os.Create("./config/config.json")
	if err != nil {
		printErr(err)
		return
	}
	defer file.Close()

	//写入文件
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		printErr(err)
		return
	}
	_, err = file.Write(data)
	if err != nil {
		printErr(err)
		return
	}
}

func checkEnv(config *Config) bool {
	var willWrite = false
	if !config.GitInstalled {
		if !checkCommand("git -v") {
			printWithEmptyLine("检测到未安装 Git ，请安装后继续")
			return false
		} else {
			config.GitInstalled = true
			willWrite = true
		}
	}
	if !config.NodeJSInstalled {
		if !checkCommand("node -v") {
			printWithEmptyLine("检测到未安装 Node.js ，请安装后继续")
			return false
		} else {
			config.NodeJSInstalled = true
			willWrite = true

		}
	}
	if !config.NpmInstalled {
		if !checkCommand("npm -v") {
			fmt.Print("检测到未安装 npm ，请手动安装Node.js，具体请看：https://note.youdao.com/s/ImCA210l")
		} else {
			config.NpmInstalled = true
			willWrite = true
		}
	}
	if willWrite {
		//写入到文件
		data, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			printErr(err)
			return false
		}
		err = os.WriteFile("./config/config.json", data, os.ModePerm)
		if err != nil {
			printErr(err)
			return false
		}
	}
	return true
}

func startRedis() *exec.Cmd {
	wd.changeToRedis()
	printWithEmptyLine("正在启动 Redis ...")
	dir, _ := os.Getwd()
	redisPath := filepath.Join(dir, "redis-server.exe")
	redisConfigPath := filepath.Join(dir, "redis.conf")
	cmd := exec.Command("cmd.exe", "/c", "start", redisPath, redisConfigPath)
	err := cmd.Start()
	if err != nil {
		printErr(err)
	}
	println("Redis 启动成功！")
	return cmd
}

func isRedisRunning() bool {
	// 执行 tasklist 命令并获取输出结果
	cmd := exec.Command("tasklist")
	output, err := cmd.Output()
	if err != nil {
		printErr(err)
	}

	// 检查输出结果中是否包含 redis-server.exe 进程
	if strings.Contains(string(output), "redis-server.exe") {
		return true
	} else {
		return false
	}
}

// 主菜单函数
func mainMenu() {
	options := []string{
		"安装云崽",
		"云崽管理",
		"BUG修复",
		"立即更新启动器",
		"获取自建签名API下载地址",
	}

	for {
		wd.changeToRoot()
		choice := showMenu("主菜单", options, true)

		switch choice {
		case 0:
			os.Exit(0)
			return
		case 1:
			clearLog()
			downloadYunzaiFromGitee()
		case 2:
			clearLog()
			manageYunzaiMenu()
		case 3:
			clearLog()
			bugsFixMenu()
		case 4:
			clearLog()
			updateLauncherRightNow()
		case 5:
			clearLog()
			getSelfSignAPI()
		}
	}
}

type Config struct {
	GitInstalled    bool   `json:"git_installed"`
	NodeJSInstalled bool   `json:"nodejs_installed"`
	NpmInstalled    bool   `json:"npm_installed"`
	SystemTempPath  string `json:"system_temp_path"`
	ConfigPath      string `json:"config_path"`
}

var (
	yunzaiName           = "Yunzai-bot"
	programName          = "YzLauncher-windows"
	globalRepositoryLink = `https://gitee.com/bling_yshs/YzLauncher-windows`
	programRunPath       = ""
	ownerAndRepo         = "bling_yshs/YzLauncher-windows"
	giteeAPI             = NewGiteeAPI()
	config               Config
	wd                         = &WorkingDirectory{}
	updating                   = false
	windowsVersion       int64 = 10
	configPath                 = ""
	updatedVersion             = version
)

const (
	version = "v0.1.28"
)

func main() {
	getAppInfoInt(&windowsVersion)
	getAppInfo(&programRunPath, &programName, &configPath, &yunzaiName)
	if checkYunzaiFileExist() {
		printRedInfo("检测到当前目录下可能存在云崽文件，请注意云崽启动器需要在云崽根目录的上一级目录下运行!")
	}
	createNormalConfig(config)
	readAndWriteSomeConfig(&config)
	updateThisProgram()
	if !checkEnv(&config) {
		shutdownApp()
	}
	checkRedis()
	println("当前版本:", version)
	getAndPrintAnnouncement()
	scheduleList()
	mainMenu()
}

func checkYunzaiFileExist() bool {
	wd.changeToRoot()
	if _, err := os.Stat("./package.json"); err == nil {
		return true
	}
	if _, err := os.Stat("./plugins"); err == nil {
		return true
	}
	return false
}

func readAndWriteSomeConfig(config *Config) {
	//读取配置文件
	file, err := os.Open("./config/config.json")
	if err != nil {
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return
	}
	writeSystemTempPath(config)
}

func getSelfSignAPI() {
	printWithEmptyLine("下载地址 https://www.123pan.com/s/tsd9-boNJv.html ，解压后放到与启动器同级目录下，然后进入解压出来的文件夹，查阅里面的 一小段说明.txt ，然后运行云崽管理->启动签名API，等待弹出的窗口显示 [FEKit_]info: task_handle.h:74 TaskSystem not allow 即为成功")
}
