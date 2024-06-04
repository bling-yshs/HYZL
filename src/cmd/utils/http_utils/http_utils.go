package http_utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func DownloadFile(url string, filePath string, showProgress bool) error {
	// 创建HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var size int
	if showProgress {
		// 获取文件总大小
		size, err = strconv.Atoi(resp.Header.Get("Content-Length"))
		if err != nil {
			return err
		}
	}

	// 如果文件已经存在，先删除
	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
	}
	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 根据showProgress参数决定是否使用带进度条的读者
	var reader io.Reader = resp.Body
	if showProgress {
		reader = &ProgressReader{
			Reader:       resp.Body,
			Total:        int64(size),
			Downloaded:   0,
			LastProgress: 0,
		}
	}

	// 将HTTP响应的内容写入文件
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}

	if showProgress {
		fmt.Println("\n下载完成!")
	}
	return nil
}

// ProgressReader 结构体用于跟踪下载进度
type ProgressReader struct {
	Reader       io.Reader
	Total        int64
	Downloaded   int64
	LastProgress int
}

// Read 方法重写，增加进度条功能
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if n > 0 {
		pr.Downloaded += int64(n)
		progress := int(float64(pr.Downloaded) / float64(pr.Total) * 100)
		if progress != pr.LastProgress {
			pr.LastProgress = progress
			fmt.Printf("\r下载进度: %d%%", progress)
		}
	}
	return n, err
}
