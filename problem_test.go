package main

import (
	"github.com/bling-yshs/YzLauncher-windows/tools"
	"path/filepath"
	"testing"
)

func TestIcqqProblemFix(t *testing.T) {
	_ = tools.AppendToYaml(filepath.Join(yunzaiName, "config/config/bot.yaml"), "sign_api_addr", "http://127.0.0.1:8080/sign")
}
