package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

type Release struct {
	Name        string                      `json:"name,omitempty" ymal:"name,omitempty"`
	TagName     string                      `json:"tag_name,omitempty" ymal:"tag_name,omitempty"`
	Description string                      `json:"description,omitempty" ymal:"description,omitempty"`
	Ref         *string                     `json:"ref,omitempty" ymal:"ref,omitempty"`
	Milestones  []string                    `json:"milestones,omitempty" ymal:"milestones,omitempty"`
	Assets      gitlab.ReleaseAssetsOptions `json:"assets,omitempty" ymal:"assets,omitempty"`
	ReleasedAt  time.Time                   `json:"released_at,omitempty" ymal:"released_at,omitempty"`
}

//	Assets struct {
//		Count   int `json:"count"`
//		Sources []struct {
//			Format string `json:"format"`
//			URL    string `json:"url"`
//		} `json:"sources"`
//		Links []*ReleaseLink `json:"links"`
//	} `json:"assets"`

func (pc ProjectCtx) createReleases() {
	for _, release := range pc.Project.Releases {
		if release.Ref == nil {
			branch := "master"
			release.Ref = &branch
			fmt.Printf("[INFO] `ref` not defined for release `%s`, defaulting to `%s`.\n", release.Name, branch)
		}
		_, _, err := pc.Client.Releases.CreateRelease(pc.ProjectID, &gitlab.CreateReleaseOptions{
			Name:        &release.Name,
			Ref:         release.Ref,
			TagName:     &release.TagName,
			Description: &release.Description,
			Milestones:  &release.Milestones,
			Assets:      &release.Assets,
			ReleasedAt:  &release.ReleasedAt,
		})
		if err != nil {
			fmt.Printf("[ERROR] Release `%s` could not be created for project `%s` -- %s\n", release.Name, pc.ProjectID, err)
		} else {
			fmt.Printf("[INFO] Release `%s` created for project `%s`.\n", release.Name, pc.ProjectID)
		}
	}
}

// TODO: Support pagination.
func getReleases(filename string) ([]Release, error) {
	content, err := getFileContents(filename)
	if err != nil {
		panic(err)
	}

	var releases []Release
	extension := filepath.Ext(filename)
	if extension == ".json" {
		err = json.Unmarshal(content, &releases)
	} else if extension == ".yaml" {
		err = yaml.Unmarshal(content, &releases)
	} else {
		err = errors.New("[ERROR] File extension not recognized, must be either `json` or `yaml`.")
	}

	return releases, err
}

func replaceReleases(pc *ProjectCtx, api API) {
	// Note that the `releases` var in each block is a different type:
	// - []*gitlab.Release
	// - []Release

	if api.ProjectID != nil {
		releases, _, err := pc.Client.Releases.ListReleases(*api.ProjectID, &gitlab.ListReleasesOptions{})
		if err != nil {
			panic(err)
		}

		var r []Release
		for _, release := range releases {
			r = append(r, Release{
				Name:        release.Name,
				TagName:     release.TagName,
				Description: release.Description,
				//					Ref:         release.Ref,
				//					Assets:     release.Assets,
				ReleasedAt: *release.ReleasedAt,
			})
		}
		pc.Project.Releases = r
	} else {
		releases, err := getReleases(*api.Filename)
		if err != nil {
			panic(err)
		}
		pc.Project.Releases = releases
	}
}
