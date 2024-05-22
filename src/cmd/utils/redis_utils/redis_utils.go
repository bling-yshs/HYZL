package redis_utils

import (
	"github.com/bling-yshs/HYZL/src/cmd/structs/global"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/windows_utils"
	ct "github.com/daviddengcn/go-colortext"
	"os/exec"
	"path"
)

func IsRedisRunning() bool {
	// 检查进程列表，看是否有redis进程, 如果有则返回true, 否则返回false
	return windows_utils.IsProcessRunning("redis-server.exe")
}

func StartRedis() {
	// 启动redis
	cmd := exec.Command("cmd", "/c", "start", "redis-server.exe", "redis.conf")
	cmd.Dir = path.Join(global.Global.ProgramRunPath, "redis-windows-7.0.4")
	err := cmd.Start()
	if err != nil {
		print_utils.PrintWithColor(ct.Red, true, "启动Redis失败！", err)
	}
}
