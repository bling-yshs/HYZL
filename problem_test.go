package main

import (
	"github.com/bling-yshs/YzLauncher-windows/tools"
	"testing"
)

func TestIcqqProblemFix(t *testing.T) {
	err := tools.UpdateValueInJSONFile("D:/AllCodeWorkspace/goland/YzLauncher-windows/Yunzai-Bot/package.json", "dependencies", "puppeteer", "19.8.3")
	if err != nil {
		t.Error(err)

	}
}
