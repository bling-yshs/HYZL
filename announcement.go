package main

import (
	"net/url"
	"os"
	"path/filepath"
)

var (
	newAnnouncementPath = ""
	oldAnnouncementPath = ""
)

func downloadAnnouncement(latestVersion string) error {
	announcementDownloadLink, _ := url.JoinPath(globalRepositoryLink, "releases", "download", latestVersion, "announcement.txt")
	downloadPath := filepath.Join(config.SystemTempPath, "YzLauncher", "download")
	err := downloadFile(announcementDownloadLink, downloadPath)
	return err
}

func compareAnnouncement() bool {
	_, err := os.Stat(oldAnnouncementPath)
	//如果文件不存在
	if err != nil {
		//复制文件
		copyFile(newAnnouncementPath, oldAnnouncementPath)
		return true
	}
	//如果文件存在，比较文件的md5值
	newAnnouncementMD5 := getFileMD5(newAnnouncementPath)
	oldAnnouncementMD5 := getFileMD5(oldAnnouncementPath)
	if newAnnouncementMD5 == oldAnnouncementMD5 {
		return false
	} else {
		copyFile(newAnnouncementPath, oldAnnouncementPath)
		return true
	}
}

func printAnnouncement() {
	announcementPath := filepath.Join(config.ConfigPath, "announcement.txt")
	_, err := os.Stat(announcementPath)
	//如果文件不存在
	if err != nil {
		return
	}
	content, err := getFileContent(announcementPath)
	if err != nil {
		return
	}
	printWithEmptyLine("公告：\n" + content)
	//删除文件
	_ = os.Remove(announcementPath)
}

func getAnnouncement() {
	latestVersion, _ := giteeAPI.getLatestTag()
	err := downloadAnnouncement(latestVersion)
	if err != nil {
		return
	}
	if compareAnnouncement() {
		//复制到Config文件夹
		copyFile(oldAnnouncementPath, filepath.Join(config.ConfigPath, "announcement.txt"))
	}
}

func getAndPrintAnnouncement() {
	newAnnouncementPath = filepath.Join(config.SystemTempPath, "YzLauncher", "download", "announcement.txt")
	oldAnnouncementPath = filepath.Join(config.SystemTempPath, "YzLauncher", "announcement.txt")
	printAnnouncement()
	go getAnnouncement()
}
