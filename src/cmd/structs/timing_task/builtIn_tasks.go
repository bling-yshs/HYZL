package timing_task

var BuiltInTasks = []TimingTask{
	{
		Id:         1,
		Name:       "CheckUpdate",
		Spec:       "0 4 * * *",
		RunOnStart: false,
		Enabled:    false,
	},
	{
		Id:         2,
		Name:       "update_yunzai_and_plugins",
		Spec:       "0 4 * * *",
		RunOnStart: false,
		Enabled:    false,
	},
	{
		Id:         3,
		Name:       "download_announcement",
		Spec:       "@yearly",
		RunOnStart: true,
		Enabled:    true,
	},
}
