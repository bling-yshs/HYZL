package schedule

import (
	"github.com/bling-yshs/HYZL/src/cmd/structs/timing_task"
)

func Execute(task timing_task.TimingTask) {
	// 执行任务
	switch task.Name {
	case "CheckUpdate":
		// 检查更新
		CheckUpdate()
	case "update_yunzai_and_plugins":
		// 更新插件
		UpdateYunzaiAndPlugins()
	case "download_announcement":
		// 下载公告
		DownloadAnnouncement()
	}
}
