package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

type GroupMap map[string]int

var GroupIDs GroupMap

type Group struct {
	Group    string    `json:"group,omitempty" yaml:"group,omitempty"`
	Parent   *string   `json:"parent,omitempty" yaml:"parent,omitempty"`
	Projects []Project `json:"projects,omitempty" yaml:"projects,omitempty"`
}

func addGroup(name string, gid int) {
	GroupIDs[name] = gid
}

func createSubgroup(g Group) *gitlab.Group {
	parentID, exists := GroupIDs[*g.Parent]
	if !exists {
		panic("Parent doesn't exist in groups map.")
	}

	git := getClient()
	var groupVisibility gitlab.VisibilityValue
	groupVisibility = "public"
	groupPath := url.PathEscape(fmt.Sprintf("%s%s", g.Group, "1"))
	groupName := url.PathEscape(g.Group)
	group, _, err := git.Groups.CreateGroup(&gitlab.CreateGroupOptions{
		Name:       &groupName,
		ParentID:   &parentID,
		Path:       &groupPath,
		Visibility: &groupVisibility,
	})
	if err != nil {
		panic(err)
	}
	return group
}

//func LookupGroup(group *gitlab.Group) int {
//	if group.Name != nil {
//		fmt.Println("got here")
//		groupID, exists := GroupIDs[*group.Name]
//		if !exists {
//			GroupIDs[*group.Name] = 4
//		}
//		return groupID
//	}
//	return 0
//}

func deleteGroup(gid int) (*gitlab.Response, error) {
	git := getClient()
	return git.Groups.DeleteGroup(gid)
}

func getGroup(g Group) (*gitlab.Group, error) {
	git := getClient()
	groups, _, err := git.Groups.ListGroups(&gitlab.ListGroupsOptions{
		Search: &g.Group,
	})
	if err != nil {
		panic(err)
	}
	if len(groups) > 0 {
		return groups[0], nil
	}
	return nil, err
}

func getGroupConfigs(filename string) ([]Group, error) {
	content, err := getFileContents(filename)
	if err != nil {
		panic(err)
	}

	var groups []Group
	extension := filepath.Ext(filename)
	if extension == ".json" {
		err = json.Unmarshal(content, &groups)
	} else if extension == ".yaml" {
		err = yaml.Unmarshal(content, &groups)
	} else {
		err = errors.New("[ERROR] File extension not recognized, must be either `json` or `yaml`.")
	}

	return groups, err
}

func getGroupIDs() GroupMap {
	return GroupIDs
}
