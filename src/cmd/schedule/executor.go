package schedule

import (
	"github.com/bling-yshs/HYZL/src/cmd/structs/timing_task"
)

func execute(taskId int) {
	// 这里是执行任务
	// 1. 从内置任务中获取任务信息
	task := timing_task.BuiltInTasks[taskId]
	// 2. 执行任务
	switch task.Name {
	case "CheckUpdate":
		// 检查更新
		CheckUpdate()
	case "update_yunzai_and_plugins":
		// 更新插件
		UpdateYunzaiAndPlugins()
	}
}
