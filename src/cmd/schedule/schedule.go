package schedule

import (
	"github.com/bling-yshs/HYZL/src/cmd/structs/app_cron"
	"github.com/bling-yshs/HYZL/src/cmd/structs/config"
	"github.com/bling-yshs/HYZL/src/cmd/structs/timing_task"
)

func InitSchedule() {
	InitBuiltInTasks()
	// 遍历 config 的 timing_task，如果有 enabled 的，那么就添加到 cron 里面
	for _, task := range config.GetConfig().TimingTasks {
		if task.Enabled {
			if task.RunOnStart {
				// 如果 run_on_start 为 true，那么就立即执行
				Execute(task)
			}
			entryId, err := app_cron.AppCronInstance.AddFunc(task.Spec, func() {
				// 这里是执行任务
				// 如果 run_now 为 true，那么就立即执行
				Execute(task)
			})
			if err != nil {
				panic(err)
				return
			}
			app_cron.TaskIdEntryIdMap[task.Id] = entryId
		}
	}
	app_cron.AppCronInstance.Start()
}

func InitBuiltInTasks() {
	// 初始化内置定时任务
	// 检查当前config中是否有download_announcement
	var hasDownloadAnnouncement = false
	for _, task := range config.GetConfig().TimingTasks {
		if task.Name == "download_announcement" {
			hasDownloadAnnouncement = true
			break
		}
	}
	if !hasDownloadAnnouncement {
		var downloadAnnouncementTask timing_task.TimingTask
		for i, task := range timing_task.BuiltInTasks {
			if task.Name == "download_announcement" {
				downloadAnnouncementTask = timing_task.BuiltInTasks[i]
				break
			}
		}
		config.GetConfig().TimingTasks = append(config.GetConfig().TimingTasks, downloadAnnouncementTask)
		config.SaveConfig()
	}
}
