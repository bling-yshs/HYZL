package schedule

import (
	"github.com/bling-yshs/HYZL/src/cmd/structs/global"
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
		current, err := version.NewVersion(global.Global.ProgramVersion)
		if err != nil {
			return
		}
		if latest.GreaterThan(current) {
			global.Config.HaveUpdate = true
			global.WriteConfig()
		} else {
			global.Config.HaveUpdate = false
			global.WriteConfig()
		}
	}
	//update := func() {
	//	// 如果有更新
	//	b, instance := updater.CheckForUpdate()
	//	if b {
	//		// 先检查是否存在已经下载的更新文件
	//		// 如果有的话
	//		if updater.UpdateTempExist() {
	//			// 比较版本号
	//			// 如果版本号不是最新
	//			if !updater.IsUpdateTempNew() {
	//				updater.DownloadUpdate(instance.Url, false)
	//			}
	//		} else {
	//			updater.DownloadUpdate(instance.Url, false)
	//		}
	//	}
	//}
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
