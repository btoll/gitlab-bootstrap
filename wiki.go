package main

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type WikiService struct {
	*BaseService
}

type Wiki struct {
	Content string `json:"content,omitempty" yaml:"content,omitempty"`
	Title   string `json:"title,omitempty" yaml:"title,omitempty"`
}

func NewWikiService(p *Provisioner) *WikiService {
	return &WikiService{
		BaseService: &BaseService{
			provisioner: p,
		},
	}
}

func (w *WikiService) Create(pc *ProjectCtx) {
	// The format will default to markdown.
	_, _, err := w.provisioner.Client.Wikis.CreateWikiPage(pc.ProjectID, &gitlab.CreateWikiPageOptions{
		Content: &pc.Project.Wiki.Content,
		Title:   &pc.Project.Wiki.Title,
	})
	if err != nil {
		fmt.Printf("[ERROR] Wiki `%s` could not be created for project `%s` -- %s\n", pc.Project.Wiki.Title, pc.ProjectID, err)
	} else {
		fmt.Printf("[INFO] Created wiki `%s` for project `%s`.\n", pc.Project.Wiki.Title, pc.ProjectID)
	}
}
