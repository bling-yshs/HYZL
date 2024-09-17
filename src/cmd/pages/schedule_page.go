package pages

import (
	"fmt"
	"github.com/bling-yshs/HYZL/src/cmd/schedule"
	"github.com/bling-yshs/HYZL/src/cmd/structs/app_cron"
	"github.com/bling-yshs/HYZL/src/cmd/structs/config"
	"github.com/bling-yshs/HYZL/src/cmd/structs/menu_option"
	"github.com/bling-yshs/HYZL/src/cmd/structs/timing_task"
	"github.com/bling-yshs/HYZL/src/cmd/structs/yunzai"
	"github.com/bling-yshs/HYZL/src/cmd/utils/cmd_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/input_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/menu_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
	"os"
)

func ScheduleMenu() {
	//检查是否存在Global.YunzaiName文件夹
	_, err := os.Stat(yunzai.GetYunzai().Name)
	if os.IsNotExist(err) {
		print_utils.PrintWithEmptyLine("未检测到云崽文件夹，请先下载云崽！")
		return
	}
	options := []menu_option.MenuOption{
		{"定时更新云崽和所有插件", updateYunzaiAndAllPlugins},
		{"定时获取启动器公告", downloadAnnouncement},
	}

	for {
		menu_utils.PrintMenu("定时任务", options, false)
		choice := input_utils.ReadUint32()
		if choice == 0 {
			cmd_utils.ClearLog()
			return
		}
		menu_utils.DealChoice(choice, options, false)
	}

}

func downloadAnnouncement() {
	taskTemplate("download_announcement")
}

func taskTemplate(taskName string) {
	// 初始化空的定时任务
	var newTask timing_task.TimingTask
	for _, item := range timing_task.BuiltInTasks {
		// 找到定时任务
		if item.Name == taskName {
			newTask = item
			break
		}
	}
	// 是否启用
	fmt.Print("是否启用(y/n):")
	enable := input_utils.ReadChoice([]string{"y", "n"})
	if enable == "y" {
		newTask.Enabled = true
	} else {
		newTask.Enabled = false
	}
	// 是否在启动器打开时立刻运行
	fmt.Print("是否在启动器打开时立刻运行(y/n):")
	runOnStart := input_utils.ReadChoice([]string{"y", "n"})
	if runOnStart == "y" {
		newTask.RunOnStart = true
	} else {
		newTask.RunOnStart = false
	}
	// 设置定时任务的 go cron 表达式
	fmt.Printf("请输入定时任务的 go cron 表达式(参考 https://pkg.go.dev/github.com/robfig/cron#hdr-Usage )(默认%s):", newTask.Spec)
	cron := input_utils.ReadString()
	if cron == "" {
		cron = newTask.Spec
	}
	newTask.Spec = cron
	// 保存定时任务
	tasks := &config.GetConfig().TimingTasks
	var found bool
	for i, task := range *tasks {
		if task.Id == newTask.Id {
			found = true
			(*tasks)[i] = newTask
		}
	}
	if !found {
		*tasks = append(*tasks, newTask)
	}
	entryId := app_cron.TaskIdEntryIdMap[newTask.Id]
	app_cron.AppCronInstance.Remove(entryId)
	entryId, err := app_cron.AppCronInstance.AddFunc(newTask.Spec, func() {
		schedule.Execute(newTask)
	})
	app_cron.TaskIdEntryIdMap[newTask.Id] = entryId
	if err != nil {
		return
	}
	config.SaveConfig()
}

func updateYunzaiAndAllPlugins() {
	taskTemplate("update_yunzai_and_plugins")
}
