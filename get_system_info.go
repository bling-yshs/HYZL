package main

import "os"

func getSystemTempPath() string {
	// 获取系统临时目录
	tempDir := os.Getenv("TEMP")
	return tempDir
}

func writeSystemTempPath(config *Config) {
	if config.SystemTempPath == "" {
		config.SystemTempPath = getSystemTempPath()
	}
}
