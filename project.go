package main

import (
	"fmt"
	"sync"

	"github.com/xanzy/go-gitlab"
)

type Project struct {
	Name       string      `json:"name,omitempty" ymal:"name,omitempty"`
	TplName    string      `json:"tpl_name,omitempty" ymal:"tpl_name,omitempty"`
	Visibility string      `json:"visibility,omitempty" ymal:"visibility,omitempty"`
	API        *[]API      `json:"api,omitempty" ymal:"api,omitempty"`
	Invites    []Invite    `json:"invites,omitempty" ymal:"invites,omitempty"`
	Issues     []IssueType `json:"issues,omitempty" ymal:"issues,omitempty"`
	Releases   []Release   `json:"releases,omitempty" ymal:"releases,omitempty"`
}

type API struct {
	Name      string  `json:"name,omitempty" ymal:"name,omitempty"`
	ProjectID *string `json:"project_id,omitempty" ymal:"project_id,omitempty"`
	Filename  *string `json:"filename,omitempty" yaml:"filename,omitempty"`
}

type ProjectCtx struct {
	Client    *gitlab.Client
	Group     *gitlab.Group
	Project   Project
	ProjectID string // projectID == NAMESPACE/PROJECT_NAME
}

type replaceFunc map[string]func(*ProjectCtx, API)

func (pc ProjectCtx) create() {
	//	var visibility gitlab.VisibilityValue = "public"
	_, _, err := pc.Client.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:         &pc.Project.Name,
		NamespaceID:  &pc.Group.ID,
		Path:         &pc.Project.Name,
		TemplateName: &pc.Project.TplName,
		//			Visibility:   &projects[i].Visibility,
	})

	if err != nil {
		fmt.Printf("[ERROR] Project `%s` could not be created -- %s\n", pc.ProjectID, err)
	} else {
		fmt.Printf("[SUCCESS] Created project `%s`\n", pc.ProjectID)
		fmt.Printf("git clone git@gitlab.com:%s.git\n", pc.ProjectID)

		if len(pc.Project.Invites) > 0 {
			pc.createInvites()
		}

		if len(pc.Project.Issues) > 0 {
			pc.createIssues()
		}

		if len(pc.Project.Releases) > 0 {
			pc.createReleases()
		}
	}
}

func (pc ProjectCtx) delete() {
	_, err := pc.Client.Projects.DeleteProject(pc.ProjectID)
	if err != nil {
		fmt.Printf("[ERROR] Project `%s` could not be deleted -- %s\n", pc.ProjectID, err)
	} else {
		fmt.Printf("[SUCCESS] Deleted project `%s`\n", pc.ProjectID)
	}
}

// This doesn't take a receiver because it needs to pass a pointer to the funcmap.
func replace(pc *ProjectCtx) {
	funcmap := replaceFunc{
		"invites":  replaceInvites,
		"issues":   replaceIssues,
		"releases": replaceReleases,
	}
	if pc.Project.API != nil {
		for _, field := range *pc.Project.API {
			funcmap[field.Name](pc, field)
		}
	}
}

func process(g Group, projects []Project, destroy bool) {
	var wg sync.WaitGroup
	wg.Add(len(projects))

	group, err := getGroup(g)
	if err != nil {
		panic(err)
	}

	pc := ProjectCtx{
		Client: getClient(),
		Group:  group,
	}

	for _, project := range projects {
		pc.Project = project
		pc.ProjectID = fmt.Sprintf("%s/%s", group.Path, project.Name)
		go func(pc ProjectCtx) {
			if !destroy {
				replace(&pc)
				pc.create()
			} else {
				pc.delete()
			}
			wg.Done()
		}(pc)
	}

	wg.Wait()
}

func processProjects(g []Group, destroy bool) {
	for i := 0; i < len(g); i++ {
		process(g[i], g[i].Projects, destroy)
	}
}
