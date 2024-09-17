package schedule

import (
	"github.com/bling-yshs/HYZL/src/cmd/structs/app"
	"github.com/bling-yshs/HYZL/src/cmd/structs/yunzai"
	"github.com/bling-yshs/HYZL/src/cmd/utils/http_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/print_utils"
	"github.com/bling-yshs/HYZL/src/cmd/utils/redis_utils"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path"
	"time"
)

func DownloadAnnouncement() {
	err := http_utils.DownloadFile(app.GetApp().AnnouncementUrl, path.Join(app.GetApp().ConfigDir, "announcement.json"), false)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "下载公告失败"))
		return
	}
}

func UpdateYunzaiAndPlugins() {
	// 先升级云崽本体
	var cmds []*exec.Cmd
	cmd := exec.Command("git", "pull")
	cmd.Dir = yunzai.GetYunzai().Path
	cmds = append(cmds, cmd)
	pluginDir := path.Join(yunzai.GetYunzai().Path, "plugins")
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		print_utils.PrintError(errors.Wrap(err, "读取插件目录失败"))
		return
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		subDir := path.Join(pluginDir, file.Name())
		gitDir := path.Join(subDir, ".git")
		packageJsonPath := path.Join(subDir, "package.json")
		if _, err := os.Stat(gitDir); err != nil {
			continue
		}
		if _, err := os.Stat(packageJsonPath); err != nil {
			continue
		}
		exec := exec.Command("git", "pull")
		exec.Dir = subDir
		cmds = append(cmds, exec)
	}
	for _, cmd := range cmds {
		cmd.Start()
	}
	for _, cmd := range cmds {
		_ = cmd.Wait()
	}
	// 最后关闭云崽并重启
	exec.Command("taskkill", "/f", "/im", "node.exe").Run()
	exec.Command("taskkill", "/f", "/im", "chrome.exe").Run()
	if !redis_utils.IsRedisRunning() {
		redis_utils.StartRedis()
		// 等待1秒
		time.Sleep(1 * time.Second)
	}
	// 启动云崽
	cmd = exec.Command("cmd", "/C", "start", "cmd", "/k", "node app")
	cmd.Dir = path.Join(yunzai.GetYunzai().Path)
	cmd.Start()
}

func CheckUpdate() {

}
