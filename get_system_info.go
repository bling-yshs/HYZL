package main

import (
	"encoding/json"
	"os"
)

func getSystemTempPath() string {
	// 获取系统临时目录
	tempDir := os.Getenv("TEMP")
	_, err := os.Stat(tempDir)
	//如果不存在就创建
	if err != nil {
		err = os.MkdirAll(tempDir, os.ModePerm)
		if err != nil {
			printErr(err)
			return ""
		}
	}
	return tempDir
}

func writeConfig(config *Config) {
	//将config写入文件
	file, err := os.Create("./config/config.json")
	if err != nil {
		printErr(err)
		return
	}
	defer file.Close()
	//写入文件
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		printErr(err)
		return
	}
	_, err = file.Write(data)
	if err != nil {
		printErr(err)
		return
	}
}

func writeSystemTempPath(config *Config) {
	if config.SystemTempPath == "" {
		config.SystemTempPath = getSystemTempPath()
		writeConfig(config)
	}
}
