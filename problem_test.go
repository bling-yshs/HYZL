package main

import (
	"testing"
)

func TestIcqqProblemFix(t *testing.T) {
	config.SystemTempPath = `C:\Users\yshs\AppData\Local\Temp`
	downloadFile("https://gitee.com/bling_yshs/yunzaiv3-ys-plugin/raw/master/%E6%8A%95%E7%A5%A8%E8%B8%A2%E4%BA%BA.js", "C:\\Users\\yshs\\AppData\\Local\\Temp\\YzLauncher\\plugins")
}
