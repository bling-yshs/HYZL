package schedule

import (
	"github.com/bling-yshs/HYZL/src/cmd/structs/app"
	"github.com/bling-yshs/HYZL/src/cmd/structs/config"
	"github.com/bling-yshs/HYZL/src/cmd/structs/updater"
	"github.com/hashicorp/go-version"
	"time"
)

func InitSchedule() {
	checkUpdate := func() {
		latestUpdater, err := updater.GetLatestUpdater()
		if err != nil {
			return
		}
		latest, err := version.NewVersion(latestUpdater.Version)
		if err != nil {
			return
		}
		current, err := version.NewVersion(app.GetApp().Version)
		if err != nil {
			return
		}
		if latest.GreaterThan(current) {
			config.GetConfig().HaveUpdate = true
			config.WriteConfig()
		} else {
			config.GetConfig().HaveUpdate = false
			config.WriteConfig()
		}
	}
	// 立刻执行一次
	go checkUpdate()
	// 每三小时执行一次
	startTicker(time.Hour*3, checkUpdate)
}

// 定时任务函数，传入时间间隔和需要定时执行的函数
func startTicker(duration time.Duration, task func()) {
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			task()
		}
	}()
}
