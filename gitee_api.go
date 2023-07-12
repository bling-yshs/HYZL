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

func NewGiteeAPI() *GiteeAPI {
	api := &GiteeAPI{}
	api.updateReleaseDataFromAPI()
	return api
}

func (api *GiteeAPI) getLatestTag() string {
	return api.latestTag
}

func (api *GiteeAPI) setLatestTag(latestTag string) {
	api.latestTag = latestTag
}

func (api *GiteeAPI) getBody() string {
	return api.body
}

func (api *GiteeAPI) setBody(body string) {
	api.body = body
}

func (api *GiteeAPI) updateReleaseDataFromAPI() {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/releases/latest", ownerAndRepo)

	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.Body == nil {
		printRedInfo(`无法连接到 Gitee API，请将此界面截图并反馈给作者，err: "resp.Body == nil" `)
	}

	body, _ := io.ReadAll(resp.Body)

	var release struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
	}

	json.Unmarshal(body, &release)
	if release.TagName == "" {
		printRedInfo(`无法连接到 Gitee API，请将此界面截图并反馈给作者，err: "release.TagName == " " `)

	}
	api.setLatestTag(release.TagName)
	api.setBody(release.Body)
}
