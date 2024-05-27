package global

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type config struct {
	GitInstalled            bool  `json:"git_installed"`
	NodeInstalled           bool  `json:"node_installed"`
	NpmInstalled            bool  `json:"npm_installed"`
	LastAnnouncementVersion int32 `json:"last_announcement_version"`
	JustFinishedUpdating    bool  `json:"just_finished_updating"`
}

var Config = config{
	GitInstalled:            false,
	NodeInstalled:           false,
	NpmInstalled:            false,
	LastAnnouncementVersion: 0,
	JustFinishedUpdating:    false,
}

type global struct {
	YunzaiName     string
	ProgramName    string
	ProgramRunPath string
	ProgramVersion string
}

var Global = global{
	YunzaiName:     yunzaiName(),
	ProgramName:    programName(),
	ProgramRunPath: programRunPath(),
	ProgramVersion: "v0.2.34",
}

func yunzaiName() string {
	_, err := os.Stat("./Miao-Yunzai")
	// 如果文件夹存在，说明是Miao-Yunzai
	if err == nil {
		return "Miao-Yunzai"
	} else {
		return "Yunzai-Bot"
	}
}

func programRunPath() string {
	// 获取程序运行的所在文件夹
	exePath, _ := os.Executable()
	return filepath.Dir(exePath)
}

func programName() string {
	// 获取程序名
	exePath, _ := os.Executable()
	return filepath.Base(exePath)
}

func WriteConfig() {
	// 写入./config/config.json
	bytes, _ := json.MarshalIndent(Config, "", "    ")
	_ = os.WriteFile("./config/config.json", bytes, os.ModePerm)
}
