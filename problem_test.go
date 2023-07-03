package main

import (
	"github.com/bling-yshs/YzLauncher-windows/tools"
	"testing"
)

func TestIcqqProblemFix(t *testing.T) {
	b, _ := tools.CheckKeyInJSONFile("./Yunzai-bot/package.json", "puppeteer")
	println(b)
}
