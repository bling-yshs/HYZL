package config

import (
	"encoding/json"
	"github.com/bling-yshs/HYZL/src/cmd/structs/timing_task"
	"github.com/bling-yshs/HYZL/src/cmd/utils/cmd_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
	ct "github.com/daviddengcn/go-colortext"
	"os"
)

type config struct {
	GitInstalled            bool                     `json:"git_installed"`
	NodeInstalled           bool                     `json:"nodejs_installed"`
	NpmInstalled            bool                     `json:"npm_installed"`
	LastAnnouncementVersion int32                    `json:"last_announcement_version"`
	JustFinishedUpdating    bool                     `json:"just_finished_updating"`
	HaveUpdate              bool                     `json:"have_update"`
	TimingTasks             []timing_task.TimingTask `json:"timing_tasks"`
}

var configInstance *config

func init() {
	configInstance = &config{}
	// 从配置文件中读取配置
	file, err := os.ReadFile("./config/config.json")
	if err == nil && len(file) > 0 {
		// 如果文件存在，那么直接读取
		err = json.Unmarshal(file, &configInstance)
		if err != nil {
			panic(err)
		}
	} else {
		// 先创建那个文件
		os.MkdirAll("./config", os.ModePerm)
		create, err := os.Create("./config/config.json")
		if err != nil {
			panic(err)
		}
		_ = create.Close()
	}
	if !configInstance.GitInstalled {
		_, err := cmd_utils.CheckCommand("git -v")
		if err != nil {
			print_utils.PrintWithColor(ct.Red, true, "未检测到 git，请先安装 git")
			os.Exit(1)
			return
		}
		configInstance.GitInstalled = true
	}
	if !configInstance.NodeInstalled {
		_, err := cmd_utils.CheckCommand("node -v")
		if err != nil {
			print_utils.PrintWithColor(ct.Red, true, "未检测到 node，请先安装 node ")
			os.Exit(1)
			return
		}
		configInstance.NodeInstalled = true
	}
	if !configInstance.NpmInstalled {
		_, err := cmd_utils.CheckCommand("npm -v")
		if err != nil {
			print_utils.PrintWithColor(ct.Red, true, "未检测到 npm，请先安装 npm")
			os.Exit(1)
			return
		}
		configInstance.NpmInstalled = true
	}
	SaveConfig()
	return
}

func GetConfig() *config {
	return configInstance
}

func SaveConfig() {
	marshal, err := json.MarshalIndent(configInstance, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("./config/config.json", marshal, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
