package main

import (
	"fmt"
	"net/url"

	"github.com/xanzy/go-gitlab"
)

type GroupMap map[string]int

type GroupService struct {
	*BaseService
	Mapper GroupMap
}

type Group struct {
	Group    string    `json:"group,omitempty" yaml:"group,omitempty"`
	Parent   *string   `json:"parent,omitempty" yaml:"parent,omitempty"`
	Projects []Project `json:"projects,omitempty" yaml:"projects,omitempty"`
}

func NewGroupService(p *Provisioner) *GroupService {
	return &GroupService{
		BaseService: &BaseService{
			provisioner: p,
		},
		Mapper: make(GroupMap),
	}
}

func (g *GroupService) AddGroupToMap(name string, gid int) {
	g.Mapper[name] = gid
}

func (g *GroupService) CreateSubgroup(gr Group) *gitlab.Group {
	parentID, exists := g.Mapper[*gr.Parent]
	if !exists {
		panic("Parent doesn't exist in groups map.")
	}

	var groupVisibility gitlab.VisibilityValue = "public"
	groupPath := url.PathEscape(fmt.Sprintf("%s%s", gr.Group, "1"))
	groupName := url.PathEscape(gr.Group)
	group, _, err := g.provisioner.Client.Groups.CreateGroup(&gitlab.CreateGroupOptions{
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
//		groupID, exists := g.Mapper[*group.Name]
//		if !exists {
//			g.Mapper[*group.Name] = 4
//		}
//		return groupID
//	}
//	return 0
//}

func (g *GroupService) Delete(gid int) (*gitlab.Response, error) {
	return g.provisioner.Client.Groups.DeleteGroup(gid, &gitlab.DeleteGroupOptions{})
}

func (g *GroupService) Get(groupName string) (*gitlab.Group, error) {
	groups, _, err := g.provisioner.Client.Groups.ListGroups(&gitlab.ListGroupsOptions{
		Search: &groupName,
	})
	if err != nil {
		panic(err)
	}
	if len(groups) > 0 {
		return groups[0], nil
	}
	return nil, err
}
