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

type ReleaseService struct {
	*BaseService
}

type Release struct {
	Name        string                      `json:"name,omitempty" yaml:"name,omitempty"`
	TagName     string                      `json:"tag_name,omitempty" yaml:"tag_name,omitempty"`
	Description string                      `json:"description,omitempty" yaml:"description,omitempty"`
	Ref         *string                     `json:"ref,omitempty" yaml:"ref,omitempty"`
	Milestones  []string                    `json:"milestones,omitempty" yaml:"milestones,omitempty"`
	Assets      gitlab.ReleaseAssetsOptions `json:"assets,omitempty" yaml:"assets,omitempty"`
	ReleasedAt  time.Time                   `json:"released_at,omitempty" yaml:"released_at,omitempty"`
}

func NewReleaseService(p *Provisioner) *ReleaseService {
	return &ReleaseService{
		BaseService: &BaseService{
			provisioner: p,
		},
	}
}

//	Assets struct {
//		Count   int `json:"count"`
//		Sources []struct {
//			Format string `json:"format"`
//			URL    string `json:"url"`
//		} `json:"sources"`
//		Links []*ReleaseLink `json:"links"`
//	} `json:"assets"`

func (r *ReleaseService) Create(pc *ProjectCtx) {
	for _, release := range pc.Project.Releases {
		if release.Ref == nil {
			branch := "master"
			release.Ref = &branch
			fmt.Printf("[INFO] `ref` not defined for release `%s`, defaulting to `%s`.\n", release.Name, branch)
		}
		_, _, err := r.provisioner.Client.Releases.CreateRelease(pc.ProjectID, &gitlab.CreateReleaseOptions{
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
			fmt.Printf("[INFO] Created release `%s` for project `%s`.\n", release.Name, pc.ProjectID)
		}
	}
}

// TODO: Support pagination.
func (r *ReleaseService) Get(filename string) ([]Release, error) {
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

//func (r *ReleaseService) Replace(pc *ProjectCtx, api API) {
//	// Note that the `releases` var in each block is a different type:
//	// - []*gitlab.Release
//	// - []Release
//
//	if api.ProjectID != nil {
//		releases, _, err := r.provisioner.Client.Releases.ListReleases(*api.ProjectID, &gitlab.ListReleasesOptions{})
//		if err != nil {
//			panic(err)
//		}
//
//		var r []Release
//		for _, release := range releases {
//			r = append(r, Release{
//				Name:        release.Name,
//				TagName:     release.TagName,
//				Description: release.Description,
//				//					Ref:         release.Ref,
//				//					Assets:     release.Assets,
//				ReleasedAt: *release.ReleasedAt,
//			})
//		}
//		pc.Project.Releases = r
//	} else {
//		releases, err := r.Get(*api.Filename)
//		if err != nil {
//			panic(err)
//		}
//		pc.Project.Releases = releases
//	}
//}

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
