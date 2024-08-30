package app

import (
	"os"
	"path/filepath"
)

type app struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Path            string `json:"path"`
	UpdateUrl       string `json:"update_url"`
	AnnouncementUrl string `json:"announcement_url"`
}

var appInstance app

func init() {
	//获取当前程序的名称
	appInstance.Name = filepath.Base(os.Args[0])
	// 设置版本
	appInstance.Version = "v0.2.71"
	// 获取当前程序的运行路径
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	appInstance.Path = path
	appInstance.UpdateUrl = "https://hyzl.r2.yshs.fun/updater/updater.json"
	appInstance.AnnouncementUrl = "https://hyzl.r2.yshs.fun/announcement/announcement.json"
}

func GetApp() app {
	return appInstance
}
