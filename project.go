package main

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type ProjectService struct {
	*BaseService
}

type Project struct {
	Name          string         `json:"name,omitempty" yaml:"name,omitempty"`
	TplName       string         `json:"tpl_name,omitempty" yaml:"tpl_name,omitempty"`
	Visibility    string         `json:"visibility,omitempty" yaml:"visibility,omitempty"`
	API           *[]API         `json:"api,omitempty" yaml:"api,omitempty"`
	Branches      []Branch       `json:"branches,omitempty" yaml:"branches,omitempty"`
	Invites       []Invite       `json:"invites,omitempty" yaml:"invites,omitempty"`
	Issues        []IssueType    `json:"issues,omitempty" yaml:"issues,omitempty"`
	Labels        []Label        `json:"labels,omitempty" yaml:"labels,omitempty"`
	MergeRequests []MergeRequest `json:"merge_requests,omitempty" yaml:"merge_requests,omitempty"`
	Releases      []Release      `json:"releases,omitempty" yaml:"releases,omitempty"`
	Wiki          *Wiki          `json:"wiki,omitempty" yaml:"wiki,omitempty"`
}

type API struct {
	Name      string  `json:"name,omitempty" yaml:"name,omitempty"`
	ProjectID *string `json:"project_id,omitempty" yaml:"project_id,omitempty"`
	Filename  *string `json:"filename,omitempty" yaml:"filename,omitempty"`
}

func NewProjectService(p *Provisioner) *ProjectService {
	return &ProjectService{
		BaseService: &BaseService{
			provisioner: p,
		},
	}
}

type ProjectCtx struct {
	Client    *gitlab.Client
	Group     *gitlab.Group
	Project   Project
	ProjectID string // projectID == NAMESPACE/PROJECT_NAME
}

//type replaceFunc map[string]func(*ProjectCtx, API)

func (p *ProjectService) Create(pc *ProjectCtx) {
	var b bool
	if pc.Project.Wiki != nil {
		b = true
	}

	var visibility gitlab.VisibilityValue = "public"
	_, _, err := p.provisioner.Client.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:         &pc.Project.Name,
		NamespaceID:  &pc.Group.ID,
		Path:         &pc.Project.Name,
		TemplateName: &pc.Project.TplName,
		Visibility:   &visibility,
		WikiEnabled:  &b,
	})

	if err != nil {
		fmt.Printf("[ERROR] Project `%s` could not be created -- %s\n", pc.ProjectID, err)
	} else {
		fmt.Printf("[SUCCESS] Created project `%s`\n", pc.ProjectID)
		fmt.Printf("git clone git@gitlab.com:%s.git\n", pc.ProjectID)

		if len(pc.Project.Branches) > 0 {
			p.provisioner.Branches.Create(pc)
		}

		if len(pc.Project.Invites) > 0 {
			p.provisioner.Invites.Create(pc)
		}

		if len(pc.Project.Issues) > 0 {
			p.provisioner.Issues.Create(pc)
		}

		if len(pc.Project.Labels) > 0 {
			p.provisioner.Labels.Create(pc)
		}

		if len(pc.Project.MergeRequests) > 0 {
			p.provisioner.MergeRequests.Create(pc)
		}

		if len(pc.Project.Releases) > 0 {
			p.provisioner.Releases.Create(pc)
		}

		if pc.Project.Wiki != nil {
			p.provisioner.Wiki.Create(pc)
		}
	}
}

func (p *ProjectService) Delete(pc *ProjectCtx) {
	_, err := p.provisioner.Client.Projects.DeleteProject(pc.ProjectID)
	if err != nil {
		fmt.Printf("[ERROR] Project `%s` could not be deleted -- %s\n", pc.ProjectID, err)
	} else {
		fmt.Printf("[SUCCESS] Deleted project `%s`\n", pc.ProjectID)
	}

	// TODO: Currently, only the parent group (i.e, "gl-group1") is added to the Mapper, so
	// this isn't really doing anything (after all, we don't want to delete the parent group,
	// and this is (poorly) making sure it's not doing that).
	// Perhaps this should be activated when another sibling group to "gl-group1" would be needed
	// to be deleted, but do we even have the ability to create a group at that same level through
	// the API (see the note on why we haven't been able to create a parent group programmatically
	// in the project README)?

	//	for k, v := range p.provisioner.Groups.Mapper {
	//		if k != "gl-group" {
	//			_, err = p.provisioner.Groups.Delete(v)
	//			if err != nil {
	//				fmt.Printf("[ERROR] Group `%s` could not be deleted -- %s\n", k, err)
	//			} else {
	//				fmt.Printf("[SUCCESS] Deleted group `%s`\n", k)
	//			}
	//		}
	//	}
}

// This doesn't take a receiver because it needs to pass a pointer to the funcmap.
//func (p *ProjectService) Replace() {
//	funcmap := replaceFunc{
//		"invites":  replaceInvites,
//		"issues":   replaceIssues,
//		"releases": replaceReleases,
//	}
//	if p.Project.API != nil {
//		for _, field := range *p.Project.API {
//			funcmap[field.Name](p, field)
//		}
//	}
//}
