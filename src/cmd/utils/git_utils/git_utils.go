package git_utils

import (
	"fmt"
	"github.com/bling-yshs/YzLauncher-windows/src/cmd/utils/cmd_utils"
)

func Clone(url string, dir string, beforeMsg string, afterMsg string) {
	cmd_utils.ExecuteCmd(fmt.Sprintf("git clone --depth=1 %s", url), dir, beforeMsg, afterMsg)
}
