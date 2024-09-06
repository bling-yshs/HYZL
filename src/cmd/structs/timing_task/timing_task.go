package timing_task

type TimingTask struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Spec       string `json:"spec"`
	RunOnStart bool   `json:"run_on_start"`
	Enabled    bool   `json:"enabled"`
}
