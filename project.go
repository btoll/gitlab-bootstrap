// TODO:
// - add proper error handling
// - add support accepting `json` or `yaml` from `stdin`?

package main

import (
	"fmt"
	"sync"

	"github.com/xanzy/go-gitlab"
)

func addProjectInvites(group *gitlab.Group, project Project) {
	git := getClient()

	projectID := fmt.Sprintf("%s/%s", group.Path, project.Name)
	for _, invite := range project.Invites {
		_, _, err := git.Invites.ProjectInvites(projectID, &gitlab.InvitesOptions{
			Email:       &invite.Email,
			AccessLevel: getAccessLevel(invite.AccessLevel),
		})
		if err != nil {
			fmt.Printf("[ERROR] Invite for `%s` could not be sent for project `%s` -- %s\n", invite.Email, projectID, err)
		} else {
			fmt.Printf("[INFO] Invite for `%s` sent for project `%s`.\n", invite.Email, projectID)
		}
	}
}

func create(group *gitlab.Group, project Project) {
	git := getClient()

	//	var visibility gitlab.VisibilityValue = "public"
	_, _, err := git.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:         &project.Name,
		NamespaceID:  &group.ID,
		Path:         &project.Name,
		TemplateName: &project.TplName,
		//			Visibility:   &projects[i].Visibility,
	})

	// projectID == NAMESPACE/PROJECT_NAME
	projectID := fmt.Sprintf("%s/%s", group.Path, project.Name)
	if err != nil {
		fmt.Printf("[ERROR] Project `%s` could not be created -- %s\n", projectID, err)
	} else {
		fmt.Printf("[SUCCESS] Created project `%s`\n", projectID)
		fmt.Printf("git clone git@gitlab.com:%s.git\n", projectID)

		if len(project.Invites) > 0 {
			addProjectInvites(group, project)
		}

		if len(project.Issues) > 0 {
			createIssues(group, project)
		}

		if len(project.Releases) > 0 {
			createReleases(group, project)
		}
	}
}

func createIssues(group *gitlab.Group, project Project) {
	git := getClient()

	projectID := fmt.Sprintf("%s/%s", group.Path, project.Name)
	for _, issue := range project.Issues {
		_, _, err := git.Issues.CreateIssue(projectID, &gitlab.CreateIssueOptions{
			Title:     &issue.Title,
			IssueType: getIssueType(issue.Type),
		})
		if err != nil {
			fmt.Printf("[ERROR] Issue `%s` could not be created for project `%s` -- %s\n", issue.Title, projectID, err)
		} else {
			fmt.Printf("[INFO] Issue `%s` created for project `%s`.\n", issue.Title, projectID)
		}
	}
}

func createReleases(group *gitlab.Group, project Project) {
	git := getClient()

	projectID := fmt.Sprintf("%s/%s", group.Path, project.Name)
	for _, release := range project.Releases {
		_, _, err := git.Releases.CreateRelease(projectID, &gitlab.CreateReleaseOptions{
			Name:    &release.Name,
			Ref:     &release.Ref,
			TagName: &release.TagName,
		})
		if err != nil {
			fmt.Printf("[ERROR] Release `%s` could not be created for project `%s` -- %s\n", release.Name, projectID, err)
		} else {
			fmt.Printf("[INFO] Release `%s` created for project `%s`.\n", release.Name, projectID)
		}
	}
}

func delete(group *gitlab.Group, project Project) {
	git := getClient()

	// projectID == NAMESPACE/PROJECT_NAME
	projectID := fmt.Sprintf("%s/%s", group.Path, project.Name)
	_, err := git.Projects.DeleteProject(projectID)
	if err != nil {
		fmt.Printf("[ERROR] Project `%s` could not be deleted -- %s\n", projectID, err)
	} else {
		fmt.Printf("[SUCCESS] Deleted project `%s`\n", projectID)
	}
}

func process(group *gitlab.Group, projects []Project, destroy bool) {
	var wg sync.WaitGroup
	wg.Add(len(projects))

	var fn func(*gitlab.Group, Project)

	if !destroy {
		fn = create
	} else {
		fn = delete
	}

	for _, project := range projects {
		go func(project Project) {
			fn(group, project)
			wg.Done()
		}(project)
	}

	wg.Wait()
}

func processProjects(g []Group, destroy bool) {
	for i := 0; i < len(g); i++ {
		group, err := getGroup(g[i])
		if err != nil {
			panic(err)
		}
		process(group, g[i].Projects, destroy)
	}
}
