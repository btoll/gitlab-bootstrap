package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/xanzy/go-gitlab"
)

type Project struct {
	Name       string      `json:"name,omitempty" ymal:"name,omitempty"`
	TplName    string      `json:"tpl_name,omitempty" ymal:"tpl_name,omitempty"`
	Visibility string      `json:"visibility,omitempty" ymal:"visibility,omitempty"`
	Invites    []Invite    `json:"invites,omitempty" ymal:"invites,omitempty"`
	Issues     []IssueType `json:"issues,omitempty" ymal:"issues,omitempty"`
	Releases   []Release   `json:"releases,omitempty" ymal:"releases,omitempty"`
}

type Release struct {
	Name        string                      `json:"name,omitempty" ymal:"name,omitempty"`
	TagName     string                      `json:"tag_name,omitempty" ymal:"tag_name,omitempty"`
	Description string                      `json:"description,omitempty" ymal:"description,omitempty"`
	Ref         *string                     `json:"ref,omitempty" ymal:"ref,omitempty"`
	Milestones  []string                    `json:"milestones,omitempty" ymal:"milestones,omitempty"`
	Assets      gitlab.ReleaseAssetsOptions `json:"assets,omitempty" ymal:"assets,omitempty"`
	ReleasedAt  time.Time                   `json:"released_at,omitempty" ymal:"released_at,omitempty"`
}

type Invite struct {
	AccessLevel string `json:"access_level,omitempty" yaml:"access_level,omitempty"`
	Email       string `json:"email,omitempty" yaml:"email,omitempty"`
}

type IssueType struct {
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	Type  string `json:"type,omitempty" yaml:"type,omitempty"`
}

const (
	Incident string = "incident"
	Issue    string = "issue"
	TestCase string = "test_case"
)

func getIssueType(issueType string) *string {
	var s string
	switch issueType {
	case Incident:
		s = "incident"
	case Issue:
		s = "issue"
	case TestCase:
		s = "test_case"
	default:
		s = "issue"
	}
	return gitlab.String(s)
}

// https://docs.gitlab.com/ee/development/permissions.html#members
const (
	None       gitlab.AccessLevelValue = 0
	Minimal    gitlab.AccessLevelValue = 5
	Guest      gitlab.AccessLevelValue = 10
	Reporter   gitlab.AccessLevelValue = 20
	Developer  gitlab.AccessLevelValue = 30
	Maintainer gitlab.AccessLevelValue = 40
	Owner      gitlab.AccessLevelValue = 50
)

func getAccessLevel(accessLevel string) *gitlab.AccessLevelValue {
	var v gitlab.AccessLevelValue
	switch accessLevel {
	case "None":
		v = None
	case "Minimal":
		v = Minimal
	case "Guest":
		v = Guest
	case "Reporter":
		v = Reporter
	case "Maintainer":
		v = Maintainer
	case "Owner":
		v = Owner
	default:
		v = Developer
	}
	return gitlab.AccessLevel(v)
}

type ProjectCtx struct {
	Client    *gitlab.Client
	Group     *gitlab.Group
	Project   Project
	ProjectID string // projectID == NAMESPACE/PROJECT_NAME
}

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
			pc.invites()
		}

		if len(pc.Project.Issues) > 0 {
			pc.issues()
		}

		if len(pc.Project.Releases) > 0 {
			pc.releases()
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

func (pc ProjectCtx) invites() {
	for _, invite := range pc.Project.Invites {
		_, _, err := pc.Client.Invites.ProjectInvites(pc.ProjectID, &gitlab.InvitesOptions{
			Email:       &invite.Email,
			AccessLevel: getAccessLevel(invite.AccessLevel),
		})
		if err != nil {
			fmt.Printf("[ERROR] Invite for `%s` could not be sent for project `%s` -- %s\n", invite.Email, pc.ProjectID, err)
		} else {
			fmt.Printf("[INFO] Invite for `%s` sent for project `%s`.\n", invite.Email, pc.ProjectID)
		}
	}
}

func (pc ProjectCtx) issues() {
	for _, issue := range pc.Project.Issues {
		_, _, err := pc.Client.Issues.CreateIssue(pc.ProjectID, &gitlab.CreateIssueOptions{
			Title:     &issue.Title,
			IssueType: getIssueType(issue.Type),
		})
		if err != nil {
			fmt.Printf("[ERROR] Issue `%s` could not be created for project `%s` -- %s\n", issue.Title, pc.ProjectID, err)
		} else {
			fmt.Printf("[INFO] Issue `%s` created for project `%s`.\n", issue.Title, pc.ProjectID)
		}
	}
}

func (pc ProjectCtx) releases() {
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
