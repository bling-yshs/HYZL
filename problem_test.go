package main

import (
	"os"
	"testing"
)

func TestIcqqProblemFix(t *testing.T) {
	// 备份标准输出
	os.Chdir("./Yunzai-Bot")
	err := executeCmd(`git pull`)
	if err != nil {
		t.Errorf("executeCmd() error = %v", err)
	}
}
