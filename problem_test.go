package main

import (
	"github.com/bling-yshs/YzLauncher-windows/tools"
	"testing"
)

func TestIcqqProblemFix(t *testing.T) {
	tools.UpdateYAMLFile("D:/AllCodeWorkspace/goland/YzLauncher-windows/Miao-Yunzai/config/config/bot.yaml", "sign_api_addr", "http://127.0.0.1:8080/sign")
}
