package annoncement

import (
	"encoding/json"
	"fmt"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/structs/global"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/print_utils"
	ct "github.com/daviddengcn/go-colortext"
	"net/http"
	"os"
	"time"
)

type announcement struct {
	Version    int32  `json:"version"`
	Content    string `json:"content"`
	Timestamp  int64  `json:"timestamp"`
	Deprecated bool   `json:"deprecated"`
}

var Announcements = []announcement{}

const url = "https://mirror.ghproxy.com/https://raw.githubusercontent.com/bling-yshs/YzLauncher-windows-announcement/main/announcement.json"

func ShowAnnouncement() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	// 展示公告
	response, err := client.Get(url)
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	defer response.Body.Close()
	// 解析json
	err = json.NewDecoder(response.Body).Decode(&Announcements)
	if err != nil {
		print_utils.PrintError(err)
		return
	}
	if global.Config.LastAnnouncementVersion >= Announcements[0].Version {
		return
	}
	print_utils.PrintWithColor(ct.Yellow, true, "公告:")
	// 展示公告
	if global.Config.LastAnnouncementVersion == 0 {
		// 如果是第一次运行，展示最新的公告
		printAnnouncement(Announcements[0])
		saveLastAnnouncementVersion()
		return
	}

	for _, item := range Announcements {
		// 展示未启用，并且所有版本号大于上次公告版本号的公告
		if !item.Deprecated && item.Version > global.Config.LastAnnouncementVersion {
			printAnnouncement(item)
		}
	}
	saveLastAnnouncementVersion()
}

func saveLastAnnouncementVersion() {
	global.Config.LastAnnouncementVersion = Announcements[0].Version
	bytes, _ := json.MarshalIndent(global.Config, "", "    ")
	_ = os.WriteFile("./config/config.json", bytes, os.ModePerm)
}

func printAnnouncement(item announcement) {
	format := time.Unix(item.Timestamp, 0).Format("2006-01-02")
	text := fmt.Sprintf("%s: %s", format, item.Content)
	print_utils.PrintWithColor(ct.Yellow, true, text)
}
