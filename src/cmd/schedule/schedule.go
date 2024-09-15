package schedule

import (
	"github.com/bling-yshs/HYZL/src/cmd/structs/app_cron"
	"github.com/bling-yshs/HYZL/src/cmd/structs/config"
)

func InitSchedule() {
	// 遍历 config 的 timing_task，如果有 enabled 的，那么就添加到 cron 里面
	for _, task := range config.GetConfig().TimingTasks {
		if task.Enabled {
			if task.RunOnStart {
				// 如果 run_on_start 为 true，那么就立即执行
				execute(task.Id)
			}
			entryId, err := app_cron.AppCronInstance.AddFunc(task.Spec, func() {
				// 这里是执行任务
				// 如果 run_now 为 true，那么就立即执行
				execute(task.Id)
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
