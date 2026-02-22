package dl

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const indexURL = "https://go.dev/dl"
const downloadUrl = indexURL + "/?mode=json&include=all"

type File struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Kind     string `json:"kind"`
	SHA256   string `json:"sha256"`
}

type Release struct {
	Version string `json:"version"`
	Files   []File `json:"files"`
}

func fetchReleases() ([]Release, error) {
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rels []Release
	if err := json.NewDecoder(resp.Body).Decode(&rels); err != nil {
		return nil, err
	}
	return rels, nil

}

func Resolve(version, os, arch string) (Release, File, string, error) {
	rels, err := fetchReleases()
	if err != nil {
		return Release{}, File{}, "", err
	}
	want := version
	// 1.20.0 -> go1.20.0
	if version[0] != 'g' {
		want = "go" + version
	}
	for _, r := range rels {
		if r.Version != want {
			continue
		}
		for _, f := range r.Files {
			if f.OS == os && f.Arch == arch && f.Kind == "archive" {
				dlUrl := indexURL + "/" + f.Filename
				return r, f, dlUrl, nil
			}
		}
		return r, File{}, "", fmt.Errorf("no archive for %s/%s", os, arch)
	}

	return Release{}, File{}, "", fmt.Errorf("version %s not found", version)
}
