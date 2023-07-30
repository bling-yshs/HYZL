// 编译： go build
package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-ps"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	version = "v0.1.42"
)

func main() {
	checkAppPermissions()
	getAppInfoInt(&windowsVersion)
	getAppInfo(&programRunPath, &programName, &configPath, &yunzaiName)
	if checkYunzaiFileExist() {
		printWithRedColor("检测到当前目录下可能存在云崽文件，请注意云崽启动器需要在云崽根目录的上一级目录下运行!")
	}
	createNormalConfig(config)
	readAndWriteSomeConfig(&config)
	updateThisProgram()
	if !checkEnv(&config) {
		shutdownApp()
	}
	println("当前版本:", version)
	getAndPrintAnnouncement()
	scheduleList()
	mainMenu()
}

// 主菜单函数
func mainMenu() {
	options := []MenuOption{
		{"安装云崽", downloadYunzaiFromGitee},
		{"云崽管理", manageYunzaiMenu},
		{"BUG修复", bugsFixMenu},
		{"立即更新启动器", updateLauncherRightNow},
		{"获取自建签名API下载地址", getSelfSignAPI},
	}

	for {
		choice := showMenu("主菜单", options, true)
		if choice == 0 {
			os.Exit(0)
			return
		}
	}
}

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
		if !checkCommandExist("git -v") {
			printWithEmptyLine("检测到未安装 Git ，请安装后继续")
			return false
		} else {
			config.GitInstalled = true
			willWrite = true
		}
	}
	if !config.NodeJSInstalled {
		if !checkCommandExist("node -v") {
			printWithEmptyLine("检测到未安装 Node.js ，请安装后继续")
			return false
		} else {
			config.NodeJSInstalled = true
			willWrite = true

		}
	}
	if !config.NpmInstalled {
		if !checkCommandExist("npm -v") {
			fmt.Print("检测到未安装 npm ，请手动安装Node.js，具体请看：https://note.youdao.com/s/ImCA210l")
		} else {
			config.NpmInstalled = true
			willWrite = true
		}
	}
	b := checkRedisExist()
	if b == true {
		config.RedisInstalled = true
		willWrite = true
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

func startRedis() error {
	err := wd.changeToRedis()
	if err != nil {
		//err则说明没有redis文件夹
		return err
	}
	printWithEmptyLine("正在启动 Redis ...")
	dir, _ := os.Getwd()
	redisPath := filepath.Join(dir, "redis-server.exe")
	redisConfigPath := filepath.Join(dir, "redis.conf")
	cmd := exec.Command("cmd.exe", "/c", "start", redisPath, redisConfigPath)
	err = cmd.Start()
	if err != nil {
		printErr(err)
	}
	println("Redis 启动成功！")
	return nil
}

func isRedisRunning() bool {
	processList, err := ps.Processes()
	if err != nil {
		printWithRedColor("无权限获取进程列表!")
		return true
	}

	isRedisRunning := false
	for _, process := range processList {
		if strings.Contains(process.Executable(), "redis") {
			isRedisRunning = true
			break
		}
	}
	if isRedisRunning {
		return true
	} else {
		return false
	}
}

type Config struct {
	GitInstalled    bool   `json:"git_installed"`
	NodeJSInstalled bool   `json:"nodejs_installed"`
	NpmInstalled    bool   `json:"npm_installed"`
	SystemTempPath  string `json:"system_temp_path"`
	RedisInstalled  bool   `json:"redis_installed"`
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

func checkAppPermissions() {
	if !checkCommandExist("dir") {
		printWithEmptyLine("当前软件权限不足，请用管理员权限运行，若使用管理员权限依然无效，那么我也没有办法")
		shutdownApp()
	}
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
	printWithEmptyLine("下载地址 https://yshs.lanzouk.com/b0a0q6sbc 密码:0000 ，解压后放到与启动器同级目录下(注意是启动器同级)，然后运行云崽管理->1. 启动签名 API 并启动云崽 即可")
}
