package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GiteeAPI struct {
	latestTag string
	body      string
}

var canConnectToGitee = true

func NewGiteeAPI() *GiteeAPI {
	api := &GiteeAPI{}
	err := api.updateReleaseDataFromAPI()
	if err != nil {
		return nil
	}
	return api
}

func (api *GiteeAPI) getLatestTag() (string, error) {
	if !canConnectToGitee {
		return "", fmt.Errorf("无法连接到 Gitee API，当前无法使用启动器的在线更新和在线公告功能！")
	}
	return api.latestTag, nil
}

func (api *GiteeAPI) setLatestTag(latestTag string) {
	api.latestTag = latestTag
}

func (api *GiteeAPI) getBody() (string, error) {
	if !canConnectToGitee {
		return "", fmt.Errorf("无法连接到 Gitee API，当前无法使用启动器的在线更新和在线公告功能！")
	}
	return api.body, nil
}

func (api *GiteeAPI) setBody(body string) {
	api.body = body
}

func (api *GiteeAPI) updateReleaseDataFromAPI() error {
	if !canConnectToGitee {
		//返回error
		return fmt.Errorf("无法连接到 Gitee API，当前无法使用启动器的在线更新和在线公告功能！")
	}
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/releases/latest", ownerAndRepo)

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var release struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
	}

	json.Unmarshal(body, &release)
	if release.TagName == "" {
		printWithRedColor(`无法连接到 Gitee API，当前无法使用启动器的在线更新和在线公告功能！`)
		canConnectToGitee = false
		return fmt.Errorf("无法连接到 Gitee API，当前无法使用启动器的在线更新和在线公告功能！")
	}
	api.setLatestTag(release.TagName)
	api.setBody(release.Body)
	return nil
}
