package schedule

import (
	"time"
)

func InitSchedule() {
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
