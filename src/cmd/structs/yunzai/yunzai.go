package yunzai

import (
	"os"
	"path"
)

type yunzai struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

var yunzaiInstance yunzai

func init() {
	// 获取程序运行路径
	if _, err := os.Stat("./Miao-Yunzai"); err == nil {
		yunzaiInstance.Name = "Miao-Yunzai"
	} else {
		yunzaiInstance.Name = "Yunzai-Bot"
	}
	dir, _ := os.Getwd()
	yunzaiInstance.Path = path.Join(dir, yunzaiInstance.Name)
}

func GetYunzai() yunzai {
	return yunzaiInstance
}
