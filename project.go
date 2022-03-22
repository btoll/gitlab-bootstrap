// TODO:
//   - add proper error handling
//   - add support accepting `json` or `yaml` from `stdin`?

package main

import (
	"fmt"
	"sync"

	"github.com/xanzy/go-gitlab"
)

func create(group *gitlab.Group, project Project) {
	git := getClient()

	//		var visibility gitlab.VisibilityValue = "public"
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
