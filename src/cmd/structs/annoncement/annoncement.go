package annoncement

import (
	"encoding/json"
	"fmt"
	"github.com/bling-yshs/HYZL/src/cmd/structs/app"
	"github.com/bling-yshs/HYZL/src/cmd/structs/config"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/pkg/errors"
	"os"
	"path"
	"time"
)

type announcement struct {
	Version    int32  `json:"version"`
	Content    string `json:"content"`
	Timestamp  int64  `json:"timestamp"`
	Deprecated bool   `json:"deprecated"`
}

var Announcements = []announcement{}

var url = app.GetApp().AnnouncementUrl

func ShowAnnouncement() {
	// 从 config 文件夹读取公告文件
	file, err := os.ReadFile(path.Join(app.GetApp().ConfigDir, "announcement.json"))
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "原因：读取本地公告文件失败"))
		return
	}
	// 解析json
	err = json.Unmarshal(file, &Announcements)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "原因：解析公告失败"))
		return
	}
	// 找到Announcements中最新的公告，并且deprecated为false
	var latestAnnouncement announcement
	for _, item := range Announcements {
		if !item.Deprecated {
			latestAnnouncement = item
			break
		}
	}
	// 判断当前config中的最新公告版本号是否小于等于最新公告的版本号
	if config.GetConfig().LastAnnouncementVersion >= latestAnnouncement.Version {
		return
	}
	// 展示公告
	printAnnouncement(latestAnnouncement)
	config.GetConfig().LastAnnouncementVersion = latestAnnouncement.Version
	config.SaveConfig()
}

func printAnnouncement(item announcement) {
	format := time.Unix(item.Timestamp, 0).Format("2006-01-02")
	text := fmt.Sprintf("%s: %s", format, item.Content)
	print_utils.PrintWithColor(ct.Yellow, true, text)
}
