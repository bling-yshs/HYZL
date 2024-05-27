package schedule

import (
	"github.com/bling-yshs/HYZL/src/cmd/structs/updater"
	"time"
)

func InitSchedule() {
	update := func() {
		// 如果有更新
		if updater.CheckForUpdate() {
			// 先检查是否存在已经下载的更新文件
			// 如果有的话
			if updater.UpdateTempExist() {
				// 比较版本号
				// 如果版本号不是最新
				if !updater.IsUpdateTempNew() {
					updater.WriteUpdaterJson()
					updater.DownloadUpdate(false)
				}
			} else {
				updater.WriteUpdaterJson()
				updater.DownloadUpdate(false)
			}
		}
	}
	// 每三小时执行一次
	startTicker(time.Hour*3, update)
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
