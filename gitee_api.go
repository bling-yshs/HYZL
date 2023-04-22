package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ReleaseData struct {
	Body            string `json:"body"`
	CreatedAt       string `json:"created_at"`
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
}
type GiteeAPI struct{}

func (api *GiteeAPI) getReleaseData() (*ReleaseData, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/releases/latest", ownerAndRepo)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release ReleaseData
	err = json.Unmarshal(body, &release)
	if err != nil {
		return nil, err
	}

	return &release, nil
}

func (api *GiteeAPI) getLatestTag() (string, error) {
	releaseData, err := api.getReleaseData()
	if err != nil {
		return "", err
	}

	return releaseData.TagName, nil
}

func (api *GiteeAPI) getBody() (string, error) {
	releaseData, err := api.getReleaseData()
	if err != nil {
		return "", err
	}

	return releaseData.Body, nil
}
