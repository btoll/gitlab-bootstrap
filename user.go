package main

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

func getUser(username string) (*gitlab.User, error) {
	git := getClient()
	users, _, err := git.Users.ListUsers(&gitlab.ListUsersOptions{
		Username: &username,
	})
	if err != nil {
		panic(err)
	}
	if len(users) > 0 {
		return users[0], nil
	}
	return nil, fmt.Errorf("[ERROR] User `%s` does not exist.", username)
}

func getUserProjects(userID int) {
	git := getClient()
	projects, _, err := git.Projects.ListUserProjects(userID, &gitlab.ListProjectsOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("len(projects)", len(projects))
	fmt.Println("projects", projects)
}
