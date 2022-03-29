package main

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type Group struct {
	Name     string    `json:"name,omitempty" yaml:"name,omitempty"`
	Projects []Project `json:"projects,omitempty" yaml:"projects,omitempty"`
}

func createGroup() {
	git := getClient()
	groupName := "derp"
	var groupVisibility gitlab.VisibilityValue
	groupVisibility = "public"
	groupPath := "https://gitlab.com/"
	group, _, err := git.Groups.CreateGroup(&gitlab.CreateGroupOptions{
		Name:       &groupName,
		Path:       &groupPath,
		Visibility: &groupVisibility,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("group", group)
}

func getGroup(g Group) (*gitlab.Group, error) {
	git := getClient()
	groups, _, err := git.Groups.ListGroups(&gitlab.ListGroupsOptions{
		Search: &g.Name,
	})
	if err != nil {
		panic(err)
	}
	if len(groups) > 0 {
		return groups[0], nil
	}
	return nil, err
}
